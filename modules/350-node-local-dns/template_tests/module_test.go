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
	hec := SetupHelmConfig("{}")

	BeforeEach(func() {
		hec.ValuesSet("global.discovery.clusterVersion", "1.15.6")
		hec.ValuesSet("global.discovery.clusterDNSAddress", "192.168.0.10")
		hec.ValuesSet("global.modulesImages.tags.common.kubeCaAuthProxy", "testtag")
	})

	Context("Test helm render", func() {

		It("Shoud successfull render helm", func() {
			hec.HelmRender()
		})

	})
})