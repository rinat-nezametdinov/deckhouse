/*
Copyright 2021 Flant CJSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package template_tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/helm"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Module :: node-local-dns :: helm temtplate", func() {
	hec := SetupHelmConfig(`{"nodeLocalDns":{}}`)

	BeforeEach(func() {
		hec.ValuesSet("global.discovery.kubernetesVersion", "1.15.6")
		hec.ValuesSet("global.discovery.clusterDNSAddress", "192.168.0.10")
	})

	Context("Test helm render", func() {

		It("Should successful render helm", func() {
			hec.HelmRender()
			Expect(hec.RenderError).ShouldNot(HaveOccurred())
		})

	})
})
