/*
Copyright 2021 Flant CJSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/kube-dns/sts_discovery",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "statefulsets",
			ApiVersion: "apps/v1",
			Kind:       "Statefulset",
			FilterFunc: ObjFilter,
		},
	},
}, foundStsInNamespaces)

func ObjFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return obj.GetNamespace(), nil
}

func foundStsInNamespaces(input *go_hook.HookInput) error {
	uniqNamespaces := make([]string, 0)
	namespaces := make(map[string]struct{})
	for _, o := range input.Snapshots["statefulsets"] {
		namespace := o.(string)
		if _, has := namespaces[namespace]; !has {
			namespaces[namespace] = struct{}{}
			uniqNamespaces = append(uniqNamespaces, namespace)
		}
	}

	input.Values.Set("kubeDns.internal.stsNamespaces", uniqNamespaces)

	return nil
}
