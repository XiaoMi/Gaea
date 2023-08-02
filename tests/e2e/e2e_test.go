package e2e

import (
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	_ "github.com/XiaoMi/Gaea/tests/e2e/sqlsupport"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "gaea e2e test")
}

var _ = ginkgo.BeforeSuite(func() {

	ginkgo.By("start gaea default.")
	err := util.StartGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea default success.")

	ginkgo.By("start gaea-cc default.")
	err = util.StartGaeaCCDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea-cc default success.")
	// 注册
	config.GetNamespaceRegisterManager().AddNamespace(config.DefaulNamespace)
	err = config.GetNamespaceRegisterManager().RegisterNamespaces()
	gomega.Expect(err).Should(gomega.BeNil())
})

var _ = ginkgo.AfterSuite(func() {
	// 删除注册
	config.GetNamespaceRegisterManager().UnRegisterNamespaces()

	ginkgo.By("stop gaea-cc default.")
	err := util.StopGaeaCCDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("stop gaea-cc default success.")
	ginkgo.By("stop gaea default.")
	err = util.StopGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("stop gaea default success.")
})
