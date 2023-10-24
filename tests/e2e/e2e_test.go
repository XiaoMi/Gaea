package e2e

import (
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/tests/config"
	_ "github.com/XiaoMi/Gaea/tests/e2e/dml"
	_ "github.com/XiaoMi/Gaea/tests/e2e/function"
	_ "github.com/XiaoMi/Gaea/tests/e2e/shard"
	_ "github.com/XiaoMi/Gaea/tests/e2e/unshard"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode.")
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "gaea E2E Testing")
}

var _ = ginkgo.BeforeSuite(func() {

	ginkgo.By("start remove sql test result.")
	err := util.RemoveLog(config.GetDefaultE2eConfig().FilepathJoin("e2e/shard/result"))
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("remove old shard result logs success.")

	err = util.RemoveLog(config.GetDefaultE2eConfig().FilepathJoin("e2e/unshard/result"))
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("remove old unshard result logs success.")

	ginkgo.By("start gaea default.")
	err = util.StartGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea default success.")

	ginkgo.By("start gaea-cc default.")
	err = util.StartGaeaCCDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea-cc default success.")
})

var _ = ginkgo.AfterSuite(func() {

	ginkgo.By("stop gaea-cc default.")
	err := util.StopGaeaCCDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("stop gaea-cc default success.")

	ginkgo.By("stop gaea default.")
	err = util.StopGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("stop gaea default success.")

})
