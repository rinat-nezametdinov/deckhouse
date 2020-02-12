package hooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Modules :: cloud-instance-manager :: hooks :: remove_csi_taints ::", func() {
	const (
		stateNodes = `
---
apiVersion: v1
kind: Node
metadata:
  name: node-1
spec:
  taints:
  - key: somekey-1
  - effect: NoSchedule
    key: node.flant.com/csi-not-bootstrapped
    value: ""
---
apiVersion: v1
kind: Node
metadata:
  name: node-2
spec:
  taints:
  - key: somekey-2
  - effect: NoSchedule
    key: node.flant.com/csi-not-bootstrapped
    value: ""
---
apiVersion: v1
kind: Node
metadata:
  name: node-3
spec:
  taints:
  - key: somekey-3
---
apiVersion: v1
kind: Node
metadata:
  name: node-4
spec: {}
`
		stateCSINode1 = `
---
apiVersion: storage.k8s.io/v1beta1
kind: CSINode
metadata:
  name: node-1
`
		stateCSINode2 = `
---
apiVersion: storage.k8s.io/v1beta1
kind: CSINode
metadata:
  name: node-2
`
	)

	f := HookExecutionConfigInit(`{"cloudInstanceManager":{"internal": {}}}`, `{}`)

	Context("Cluster is empty", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("Hook must execute successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
		})
	})

	Context("Cluster has four nodes and single CSINode", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateNodes + stateCSINode1))
			f.RunHook()
		})

		It("node-1 must lose taint 'node.flant.com/csi-not-bootstrapped'", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesResource("Node", "", "node-1").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-1"}]`))
			Expect(f.KubernetesResource("Node", "", "node-2").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-2"},{"effect":"NoSchedule","key":"node.flant.com/csi-not-bootstrapped","value":""}]`))
			Expect(f.KubernetesResource("Node", "", "node-3").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-3"}]`))
		})

		Context("CSINode for node-2 added", func() {
			BeforeEach(func() {
				f.BindingContexts.Set(f.KubeStateSet(stateNodes + stateCSINode1 + stateCSINode2))
				f.RunHook()
			})

			It("node-2 must lose taint 'node.flant.com/csi-not-bootstrapped'", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.KubernetesResource("Node", "", "node-1").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-1"}]`))
				Expect(f.KubernetesResource("Node", "", "node-2").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-2"}]`))
				Expect(f.KubernetesResource("Node", "", "node-3").Field("spec.taints").String()).To(MatchJSON(`[{"key":"somekey-3"}]`))
			})
		})

	})
})