// Copyright 2021 Flant CJSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converge

import (
	"fmt"

	"github.com/deckhouse/deckhouse/dhctl/pkg/config"
	"github.com/deckhouse/deckhouse/dhctl/pkg/kubernetes/client"
	"github.com/deckhouse/deckhouse/dhctl/pkg/terraform"
	"github.com/deckhouse/deckhouse/dhctl/pkg/util/tomb"
)

func NodeName(cfg *config.MetaConfig, nodeGroupName string, index int) string {
	return fmt.Sprintf("%s-%s-%v", cfg.ClusterPrefix, nodeGroupName, index)
}

func BootstrapAdditionalNode(kubeCl *client.KubernetesClient, cfg *config.MetaConfig, index int, step, nodeGroupName, cloudConfig string, isConverge bool) error {
	nodeName := NodeName(cfg, nodeGroupName, index)

	if isConverge {
		nodeExists, err := IsNodeExistsInCluster(kubeCl, nodeName)
		if err != nil {
			return err
		} else if nodeExists {
			return fmt.Errorf("node with name %s exists in cluster", nodeName)
		}
	}

	nodeConfig := cfg.NodeGroupConfig(nodeGroupName, index, cloudConfig)
	nodeGroupSettings := cfg.FindTerraNodeGroup(nodeGroupName)

	runner := terraform.NewRunnerFromConfig(cfg, step).
		WithVariables(nodeConfig).
		WithName(nodeName).
		WithAutoApprove(true)
	tomb.RegisterOnShutdown(nodeName, runner.Stop)

	if isConverge {
		runner.WithIntermediateStateSaver(NewNodeStateSaver(kubeCl, nodeName, nodeGroupName, nodeGroupSettings))
	}

	outputs, err := terraform.ApplyPipeline(runner, nodeName, terraform.OnlyState)
	if err != nil {
		return err
	}

	if tomb.IsInterrupted() {
		return ErrConvergeInterrupted
	}

	err = SaveNodeTerraformState(kubeCl, nodeName, nodeGroupName, outputs.TerraformState, nodeGroupSettings)
	if err != nil {
		return err
	}

	return nil
}

func BootstrapAdditionalMasterNode(kubeCl *client.KubernetesClient, cfg *config.MetaConfig, index int, cloudConfig string, isConverge bool) (*terraform.PipelineOutputs, error) {
	nodeName := NodeName(cfg, MasterNodeGroupName, index)

	if isConverge {
		nodeExists, existsErr := IsNodeExistsInCluster(kubeCl, nodeName)
		if existsErr != nil {
			return nil, existsErr
		} else if nodeExists {
			return nil, fmt.Errorf("node with name %s exists in cluster", nodeName)
		}
	}

	nodeConfig := cfg.NodeGroupConfig(MasterNodeGroupName, index, cloudConfig)

	runner := terraform.NewRunnerFromConfig(cfg, "master-node").
		WithVariables(nodeConfig).
		WithName(nodeName).
		WithAutoApprove(true)
	tomb.RegisterOnShutdown(nodeName, runner.Stop)

	// Node group settings are not required for master node secret.
	if isConverge {
		runner.WithIntermediateStateSaver(NewNodeStateSaver(kubeCl, nodeName, MasterNodeGroupName, nil))
	}

	outputs, err := terraform.ApplyPipeline(runner, nodeName, terraform.GetMasterNodeResult)
	if err != nil {
		return nil, err
	}

	if tomb.IsInterrupted() {
		return nil, ErrConvergeInterrupted
	}

	err = SaveMasterNodeTerraformState(kubeCl, nodeName, outputs.TerraformState, []byte(outputs.KubeDataDevicePath))
	if err != nil {
		return outputs, err
	}

	return outputs, err
}
