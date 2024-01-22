package e2e

import (
	"testing"
	"time"

	_ "github.com/XiaoMi/Gaea/tests/e2e/function"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	_ "github.com/XiaoMi/Gaea/tests/e2e/shard"
	_ "github.com/XiaoMi/Gaea/tests/e2e/unshard"
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
	ginkgo.By("start gaea default.")
	err := util.StartGaeaDefault()
	util.ExpectNoError(err)
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea default success.")

	ginkgo.By("start gaea-cc default.")
	err = util.StartGaeaCCDefault()
	util.ExpectNoError(err)
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea-cc default success.")
})

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("stop gaea-cc default.")
	err := util.StopGaeaCCDefault()
	util.ExpectNoError(err)
	ginkgo.By("stop gaea-cc default success.")

	ginkgo.By("stop gaea default.")
	err = util.StopGaeaDefault()
	util.ExpectNoError(err)
	ginkgo.By("stop gaea default success.")

})
