package sql

import (
	"testing"
	"time"

	_ "github.com/XiaoMi/Gaea/tests/sql/dml"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestIntergration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "gaea intergration test")
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
