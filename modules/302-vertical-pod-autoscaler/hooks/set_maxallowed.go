package hooks

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	autoscaler "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	"k8s.io/utils/pointer"
)

/*
Overview:
   1. All system components require resource requests, managed by vpa.
   2. Sum of all resource requests should not exceed manually configured resources limits.
   3. We expect that resources limits to be allocated fairly between vpa requests.
We have 3 groups of resources:
   1. ControlPlane - resources for pods in control-plane (kube-controller-manager, kube-scheduler, kube-apiserver, etcd).
   2. Master - vpa resources, working on master nodes (label "workload-resource-policy.deckhouse.io: master").
   3. EveryNode - vpa resources, working on every node (label "workload-resource-policy.deckhouse.io: every-node").
Calculate steps:
   1. We calculate sum of uncappedTargets requests for all vpa resources in Master group, and proportionally sets MaxAllowed values for this resources,
      based on resources requests from global config for Master group.
   2. We calculate sum of uncappedTargets requests for all vpa resources in EveryNode group, and proportionally sets MaxAllowed values for this resources,
      based on resources requests from global config for EveryNode group.
Hook start conditions:
   1. If uncappedTarget value changed in vpa with labels "workload-resource-policy.deckhouse.io: master" or "workload-resource-policy.deckhouse.io: every-node".
   2. If user changed global.modules.resourcesRequests values.
   3. By crontab to process situation, if nodes resources changed.
*/

const (
	groupLabelKey  = "workload-resource-policy.deckhouse.io"
	everyNodeLabel = "every-node"
	masterLabel    = "master"
)

type VPA struct {
	Name                     string
	Namespace                string
	Label                    string
	ContainerRecommendations []autoscaler.RecommendedContainerResources
}

func applyVpaResourcesFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	v := &autoscaler.VerticalPodAutoscaler{}
	err := sdk.FromUnstructured(obj, v)
	if err != nil {
		return nil, fmt.Errorf("cannot parse vpa object from unstructured: %v", err)
	}

	if v.Status.Recommendation == nil {
		return nil, nil
	}
	recommendations := v.Status.Recommendation.ContainerRecommendations

	c := &VPA{}
	c.Name = v.Name
	c.Namespace = v.Namespace
	c.Label = v.Labels[groupLabelKey]
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].ContainerName < recommendations[j].ContainerName
	})
	c.ContainerRecommendations = recommendations
	return c, nil
}

var (
	_ = sdk.RegisterFunc(&go_hook.HookConfig{
		OnBeforeHelm: &go_hook.OrderedConfig{Order: 20},
		Queue:        "/modules/vertical-pod-autoscaler",
		Schedule: []go_hook.ScheduleConfig{
			{Name: "vpaCron", Crontab: "*/30 * * * *"},
		},
		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:                   "Vpa",
				WaitForSynchronization: pointer.BoolPtr(false),
				ExecuteHookOnEvents:    pointer.BoolPtr(false),
				ApiVersion:             "autoscaling.k8s.io/v1",
				Kind:                   "VerticalPodAutoscaler",
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      "heritage",
							Operator: metav1.LabelSelectorOpIn,
							Values:   []string{"deckhouse"},
						},
						{
							Key:      groupLabelKey,
							Operator: metav1.LabelSelectorOpIn,
							Values:   []string{everyNodeLabel, masterLabel},
						},
					},
				},
				FilterFunc: applyVpaResourcesFilter,
			},
		},
	}, updateVpaResources)
)

func updateVpaResources(input *go_hook.HookInput) error {
	var (
		configEveryNodeMilliCPU  int64
		configEveryNodeMemory    int64
		configMasterNodeMilliCPU int64
		configMasterNodeMemory   int64

		totalRequestsMasterNodeMilliCPU int64
		totalRequestsMasterNodeMemory   int64
		totalRequestsEveryNodeMilliCPU  int64
		totalRequestsEveryNodeMemory    int64
	)

	configEveryNodeMilliCPU, err := getPathInt(input, "global.modules.resourcesRequests.internal.milliCpuEveryNode")
	if err != nil {
		return err
	}
	configEveryNodeMemory, err = getPathInt(input, "global.modules.resourcesRequests.internal.memoryEveryNode")
	if err != nil {
		return err
	}
	configMasterNodeMilliCPU, err = getPathInt(input, "global.modules.resourcesRequests.internal.milliCpuMaster")
	if err != nil {
		return err
	}
	configMasterNodeMemory, err = getPathInt(input, "global.modules.resourcesRequests.internal.memoryMaster")
	if err != nil {
		return err
	}

	snapshots := input.Snapshots["Vpa"]
	if len(snapshots) == 0 {
		return nil
	}

	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}

		v := snapshot.(*VPA)

		for _, r := range v.ContainerRecommendations {
			ut := r.UncappedTarget
			switch v.Label {
			case masterLabel:
				totalRequestsMasterNodeMilliCPU += ut.Cpu().MilliValue()
				totalRequestsMasterNodeMemory += ut.Memory().Value()
			case everyNodeLabel:
				totalRequestsEveryNodeMilliCPU += ut.Cpu().MilliValue()
				totalRequestsEveryNodeMemory += ut.Memory().Value()
			}
		}
	}

	for _, snapshot := range snapshots {
		var containerPolicies []autoscaler.ContainerResourcePolicy

		if snapshot == nil {
			continue
		}
		v := snapshot.(*VPA)

		for _, r := range v.ContainerRecommendations {
			var (
				recommendationsMilliCPU int64
				recommendationsMemory   int64
			)
			switch v.Label {
			case masterLabel:
				recommendationsMilliCPU = r.UncappedTarget.Cpu().MilliValue() * configMasterNodeMilliCPU / totalRequestsMasterNodeMilliCPU
				recommendationsMemory = r.UncappedTarget.Memory().Value() * configMasterNodeMemory / totalRequestsMasterNodeMemory
			case everyNodeLabel:
				recommendationsMilliCPU = r.UncappedTarget.Cpu().MilliValue() * configEveryNodeMilliCPU / totalRequestsEveryNodeMilliCPU
				recommendationsMemory = r.UncappedTarget.Memory().Value() * configEveryNodeMemory / totalRequestsEveryNodeMemory
			}
			newContainerPolicy := autoscaler.ContainerResourcePolicy{ContainerName: r.ContainerName}
			newContainerPolicy.MaxAllowed = v1.ResourceList{
				v1.ResourceCPU:    *resource.NewMilliQuantity(recommendationsMilliCPU, resource.BinarySI),
				v1.ResourceMemory: *resource.NewQuantity(recommendationsMemory, resource.DecimalExponent),
			}
			containerPolicies = append(containerPolicies, newContainerPolicy)
		}

		newResourcePolicyPatch := autoscaler.VerticalPodAutoscaler{Spec: autoscaler.VerticalPodAutoscalerSpec{ResourcePolicy: &autoscaler.PodResourcePolicy{ContainerPolicies: containerPolicies}}}
		jsonPatch, err := json.Marshal(newResourcePolicyPatch)
		if err != nil {
			return err
		}

		err = input.ObjectPatcher.MergePatchObject(jsonPatch, "autoscaling.k8s.io/v1", "VerticalPodAutoscaler", v.Namespace, v.Name, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func getPathInt(input *go_hook.HookInput, path string) (int64, error) {
	if !input.Values.Exists(path) {
		return 0, fmt.Errorf("%s must be set", path)
	}
	return input.Values.Get(path).Int(), nil
}

/*
err = input.ObjectPatcher.FilterObject(func(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var (
		recommendationsMilliCPU int64
		recommendationsMemory   int64
		containerPolicies       []autoscaler.ContainerResourcePolicy
	)

	v := &autoscaler.VerticalPodAutoscaler{}
	err := sdk.FromUnstructured(obj, v)
	if err != nil {
		return nil, fmt.Errorf("cannot parse vpa object from unstructured: %v", err)
	}

	for _, container := range v.Status.Recommendation.ContainerRecommendations {
		switch v.Labels[groupLabelKey] {
		case masterLabel:
			recommendationsMilliCPU = container.UncappedTarget.Cpu().MilliValue() * configMasterNodeMilliCPU / totalRequestsMasterNodeMilliCPU
			recommendationsMemory = container.UncappedTarget.Memory().Value() * configMasterNodeMemory / totalRequestsMasterNodeMemory
		case everyNodeLabel:
			recommendationsMilliCPU = container.UncappedTarget.Cpu().MilliValue() * configEveryNodeMilliCPU / totalRequestsEveryNodeMilliCPU
			recommendationsMemory = container.UncappedTarget.Memory().Value() * configEveryNodeMemory / totalRequestsEveryNodeMemory
		}
		newContainerPolicy := autoscaler.ContainerResourcePolicy{ContainerName: container.ContainerName}
		for _, cp := range v.Spec.ResourcePolicy.ContainerPolicies {
			if cp.ContainerName == container.ContainerName {
				newContainerPolicy = cp
				break
			}
		}
		newContainerPolicy.MaxAllowed = v1.ResourceList{
			v1.ResourceCPU:    *resource.NewMilliQuantity(recommendationsMilliCPU, resource.BinarySI),
			v1.ResourceMemory: *resource.NewQuantity(recommendationsMemory, resource.DecimalExponent),
		}
		containerPolicies = append(containerPolicies, newContainerPolicy)
	}
	if v.Spec.ResourcePolicy == nil {
		v.Spec.ResourcePolicy = &autoscaler.PodResourcePolicy{}
	}

	v.Spec.ResourcePolicy.ContainerPolicies = containerPolicies

	result, err := sdk.ToUnstructured(v)
	if err != nil {
		return nil, fmt.Errorf("cannot parse unstructured to object: %v", err)
	}
	return result, nil
}, v.APIVersion, v.Kind, v.Namespace, v.Name, "")
*/
