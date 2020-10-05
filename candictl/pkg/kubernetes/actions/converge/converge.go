package converge

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"

	"flant/candictl/pkg/config"
	"flant/candictl/pkg/kubernetes/client"
	"flant/candictl/pkg/log"
	"flant/candictl/pkg/terraform"
	"flant/candictl/pkg/util/retry"
)

const masterNodeGroupName = "master"

func BootstrapAdditionalNode(kubeCl *client.KubernetesClient, cfg *config.MetaConfig, index int, step, nodeGroupName, cloudConfig string) error {
	nodeName := fmt.Sprintf("%s-%s-%v", cfg.ClusterPrefix, nodeGroupName, index)
	nodeConfig := cfg.NodeGroupConfig(nodeGroupName, index, cloudConfig)

	runner := terraform.NewRunnerFromConfig(cfg, step).
		WithVariables(nodeConfig).
		WithName(nodeName).
		WithAutoApprove(true)

	defer runner.Close()

	outputs, err := terraform.ApplyPipeline(runner, nodeName, terraform.OnlyState)
	if err != nil {
		_ = runner.Destroy()
		return err
	}

	err = SaveNodeTerraformState(kubeCl, nodeName, nodeGroupName, outputs.TerraformState, cfg.FindStaticNodeGroup(nodeGroupName))
	// If we failed to save state into cluster, node doesn't exist for us. Let's destroy it.
	if err != nil {
		_ = runner.Destroy()
		return err
	}
	return nil
}

func BootstrapAdditionalMasterNode(kubeCl *client.KubernetesClient, cfg *config.MetaConfig, index int, cloudConfig string) error {
	nodeName := fmt.Sprintf("%s-%s-%v", cfg.ClusterPrefix, masterNodeGroupName, index)
	nodeConfig := cfg.NodeGroupConfig(masterNodeGroupName, index, cloudConfig)

	runner := terraform.NewRunnerFromConfig(cfg, "master-node").
		WithVariables(nodeConfig).
		WithName(nodeName).
		WithAutoApprove(true)

	defer runner.Close()

	outputs, err := terraform.ApplyPipeline(runner, nodeName, terraform.GetMasterNodeResult)
	if err != nil {
		_ = runner.Destroy()
		return err
	}

	err = SaveMasterNodeTerraformState(kubeCl, nodeName, outputs.TerraformState, []byte(outputs.KubeDataDevicePath))
	// If we failed to save state into cluster, node doesn't exist for us. Let's destroy it.
	if err != nil {
		_ = runner.Destroy()
		return err
	}
	return nil
}

func RunConverge(kubeCl *client.KubernetesClient, metaConfig *config.MetaConfig) error {
	if err := updateClusterState(kubeCl, metaConfig); err != nil {
		return err
	}

	var nodesState map[string]NodeGroupTerraformState
	var err error
	err = log.Process("converge", "Gather Nodes Terraform state", func() error {
		nodesState, err = GetNodesStateFromCluster(kubeCl)
		if err != nil {
			return fmt.Errorf("terraform nodes state in Kubernetes cluster not found: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(nodesState) == 0 {
		if !retry.AskForConfirmation("Cluster has no nodes created by Terraform. Do you want to continue and create nodes") {
			log.InfoLn("Aborted")
			return nil
		}
	}

	var nodeGroupsWithStateInCluster []string
	for _, group := range metaConfig.GetStaticNodeGroups() {
		// Skip if node group terraform state exists, we will update node group state below
		if _, ok := nodesState[group.Name]; ok {
			nodeGroupsWithStateInCluster = append(nodeGroupsWithStateInCluster, group.Name)
			continue
		}
		if err := createPreviouslyNotExistentNodeGroup(kubeCl, metaConfig, group); err != nil {
			return err
		}
	}

	for _, nodeGroupName := range sortNodeGroupsStateKeys(nodesState, nodeGroupsWithStateInCluster) {
		controller := NewConvergeController(kubeCl, metaConfig)
		if err := controller.Run(nodeGroupName, nodesState[nodeGroupName]); err != nil {
			return err
		}
	}
	return nil
}

func updateClusterState(kubeCl *client.KubernetesClient, metaConfig *config.MetaConfig) error {
	return log.Process("converge", "Update Cluster Terraform state", func() error {
		clusterState, err := GetClusterStateFromCluster(kubeCl)
		if err != nil {
			return fmt.Errorf("terraform cluster state in Kubernetes cluster not found: %w", err)
		}

		if clusterState == nil {
			return fmt.Errorf("kubernetes cluster has no state")
		}

		baseRunner := terraform.NewRunnerFromConfig(metaConfig, "base-infrastructure").
			WithVariables(metaConfig.MarshalConfig()).
			WithState(clusterState)

		outputs, err := terraform.ApplyPipeline(baseRunner, "Kubernetes cluster", terraform.GetBaseInfraResult)
		if err != nil {
			return err
		}

		if err := SaveClusterTerraformState(kubeCl, outputs); err != nil {
			return err
		}
		return nil
	})
}

func createPreviouslyNotExistentNodeGroup(kubeCl *client.KubernetesClient, metaConfig *config.MetaConfig, group config.StaticNodeGroupSpec) error {
	return log.Process("converge", fmt.Sprintf("Add NodeGroup %s (replicas: %v)️", group.Name, group.Replicas), func() error {
		err := CreateNodeGroup(kubeCl, group.Name, metaConfig.NodeGroupManifest(group))
		if err != nil {
			return err
		}

		nodeCloudConfig, err := GetCloudConfig(kubeCl, group.Name)
		if err != nil {
			return err
		}

		for i := 0; i < group.Replicas; i++ {
			err = BootstrapAdditionalNode(kubeCl, metaConfig, i, "static-node", group.Name, nodeCloudConfig)
			if err != nil {
				return err
			}
		}

		if err := WaitForNodesBecomeReady(kubeCl, group.Name, group.Replicas); err != nil {
			return err
		}
		return nil
	})
}

type Controller struct {
	client *client.KubernetesClient
	config *config.MetaConfig
}

type NodeGroupGroupOptions struct {
	Name        string
	Step        string
	CloudConfig string
	Replicas    int
	State       map[string][]byte
}

func NewConvergeController(kubeCl *client.KubernetesClient, metaConfig *config.MetaConfig) *Controller {
	return &Controller{client: kubeCl, config: metaConfig}
}

func (c *Controller) Run(nodeGroupName string, nodeGroupState NodeGroupTerraformState) error {
	replicas := getReplicasByNodeGroupName(c.config, nodeGroupName)
	step := GetStepByNodeGroupName(nodeGroupName)

	nodeCloudConfig, err := GetCloudConfig(c.client, nodeGroupName)
	if err != nil {
		return err
	}

	nodeGroup := NodeGroupGroupOptions{
		Name:        nodeGroupName,
		Step:        step,
		Replicas:    replicas,
		CloudConfig: nodeCloudConfig,
		State:       nodeGroupState.State,
	}

	if replicas > len(nodeGroupState.State) {
		err := log.Process("converge", fmt.Sprintf("Add Nodes to NodeGroup %s (replicas: %v)", nodeGroupName, replicas), func() error {
			return c.addNewNodeGroup(&nodeGroup)
		})
		if err != nil {
			return err
		}
	}

	var allErrs *multierror.Error
	if replicas != 0 {
		for name := range nodeGroupState.State {
			err := log.Process("converge", fmt.Sprintf("Update Node %s in NodeGroup %s (replicas: %v)", name, nodeGroupName, replicas), func() error {
				return c.updateNode(&nodeGroup, name)
			})
			if err != nil {
				allErrs = multierror.Append(allErrs, fmt.Errorf("%s: %v", name, err))
			}
		}
	}

	if err := allErrs.ErrorOrNil(); err != nil {
		return err
	}

	if replicas < len(nodeGroupState.State) {
		err := log.Process("converge", fmt.Sprintf("Delete Nodes from NodeGroup %s (replicas: %v)", nodeGroupName, replicas), func() error {
			return c.deleteRedundantNodes(&nodeGroup, nodeGroupState.Settings)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) addNewNodeGroup(nodeGroup *NodeGroupGroupOptions) error {
	count := len(nodeGroup.State)
	index := 0

	for nodeGroup.Replicas > count {
		candidateName := fmt.Sprintf("%s-%s-%v", c.config.ClusterPrefix, nodeGroup.Name, index)
		if _, ok := nodeGroup.State[candidateName]; !ok {
			var err error
			if nodeGroup.Name == masterNodeGroupName {
				err = BootstrapAdditionalMasterNode(c.client, c.config, index, nodeGroup.CloudConfig)
			} else {
				err = BootstrapAdditionalNode(c.client, c.config, index, nodeGroup.Step, nodeGroup.Name, nodeGroup.CloudConfig)
			}
			if err != nil {
				return err
			}
			count++
		}
		index++
	}
	return WaitForNodesBecomeReady(c.client, nodeGroup.Name, nodeGroup.Replicas)
}

func (c *Controller) updateNode(nodeGroup *NodeGroupGroupOptions, nodeName string) error {
	state := nodeGroup.State[nodeName]
	index, ok := getIndexFromNodeName(nodeName)
	if !ok {
		log.ErrorF("can't extract index from terraform state secret, skip %s\n", nodeName)
		return nil
	}

	nodeRunner := terraform.NewRunnerFromConfig(c.config, nodeGroup.Step).
		WithVariables(c.config.NodeGroupConfig(nodeGroup.Name, int(index), nodeGroup.CloudConfig)).
		WithState(state).
		WithName(nodeName)

	outputs, err := terraform.ApplyPipeline(nodeRunner, nodeName, terraform.OnlyState)
	if err != nil {
		return err
	}

	nodeGroupSettingsFromConfig := c.config.FindStaticNodeGroup(nodeGroup.Name)
	err = SaveNodeTerraformState(c.client, nodeName, nodeGroup.Name, outputs.TerraformState, nodeGroupSettingsFromConfig)
	if err != nil {
		return err
	}

	return WaitForSingleNodeBecomeReady(c.client, nodeName)
}

func (c *Controller) deleteRedundantNodes(nodeGroup *NodeGroupGroupOptions, settings []byte) error {
	deleteNodesNames := make(map[string][]byte)
	count := len(nodeGroup.State)

	for name, state := range nodeGroup.State {
		deleteNodesNames[name] = state
		delete(nodeGroup.State, name)
		count--

		if count == nodeGroup.Replicas {
			break
		}
	}

	cfg := c.config
	if settings != nil {
		nodeGroupsSettings, err := json.Marshal([]json.RawMessage{settings})
		if err != nil {
			log.ErrorLn(err)
		} else {
			cfg = c.config.DeepCopy().Prepare()
			cfg.ProviderClusterConfig["nodeGroups"] = nodeGroupsSettings
		}
	}

	var allErrs *multierror.Error
	for name, state := range deleteNodesNames {
		index, ok := getIndexFromNodeName(name)
		if !ok {
			log.ErrorF("can't extract index from terraform state secret, skip %s\n", name)
			continue
		}

		nodeRunner := terraform.NewRunnerFromConfig(c.config, nodeGroup.Step).
			WithVariables(cfg.NodeGroupConfig(nodeGroup.Name, int(index), nodeGroup.CloudConfig)).
			WithState(state).
			WithName(name).
			WithAutoApprove(true)

		if err := terraform.DestroyPipeline(nodeRunner, name); err != nil {
			allErrs = multierror.Append(allErrs, fmt.Errorf("%s: %w", name, err))
			continue
		}

		nodeRunner.Close()
		err := DeleteTerraformState(c.client, fmt.Sprintf("d8-node-terraform-state-%s", name))
		if err != nil {
			allErrs = multierror.Append(allErrs, fmt.Errorf("%s: %w", name, err))
			continue
		}
	}
	return allErrs.ErrorOrNil()
}

func getIndexFromNodeName(name string) (int64, bool) {
	index, err := strconv.ParseInt(name[strings.LastIndex(name, "-")+1:], 10, 64)
	if err != nil {
		log.ErrorLn(err)
		return 0, false
	}
	return index, true
}

func getReplicasByNodeGroupName(metaConfig *config.MetaConfig, nodeGroupName string) int {
	replicas := 0
	if nodeGroupName != masterNodeGroupName {
		for _, group := range metaConfig.GetStaticNodeGroups() {
			if group.Name == nodeGroupName {
				replicas = group.Replicas
				break
			}
		}
	} else {
		replicas = metaConfig.MasterNodeGroupSpec.Replicas
	}
	return replicas
}

func GetStepByNodeGroupName(nodeGroupName string) string {
	step := "static-node"
	if nodeGroupName == masterNodeGroupName {
		step = "master-node"
	}
	return step
}

func sortNodeGroupsStateKeys(state map[string]NodeGroupTerraformState, sortedNodeGroupsFromConfig []string) []string {
	nodeGroupsFromConfigSet := make(map[string]struct{}, len(sortedNodeGroupsFromConfig))
	for _, key := range sortedNodeGroupsFromConfig {
		nodeGroupsFromConfigSet[key] = struct{}{}
	}

	sortedKeys := append([]string{masterNodeGroupName}, sortedNodeGroupsFromConfig...)

	for key := range state {
		if key == masterNodeGroupName {
			continue
		}

		if _, ok := nodeGroupsFromConfigSet[key]; !ok {
			sortedKeys = append(sortedKeys, key)
		}
	}

	return sortedKeys
}