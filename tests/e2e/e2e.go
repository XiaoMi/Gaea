package e2e

import (
	"database/sql"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

var err error
var cfg *config.Config
var mDB *sql.DB

func RunE2ETests(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "gaea e2e suite")
}

var _ = ginkgo.BeforeSuite(func() {
	ginkgo.By("begin prepare database.")
	// init db
	cfg = config.TestConfig
	mDB, err = util.Open(cfg.MasterHost, cfg.MasterPort, cfg.MasterUser, cfg.MasterPassword, 10, 3)
	defer mDB.Close()
	gomega.Expect(err).Should(gomega.BeNil())

	// delete cluster data
	_, err = mDB.Exec(cfg.GenDropDatabaseSqls())
	gomega.Expect(err).Should(gomega.BeNil())
	_, err = mDB.Exec(cfg.GenCreateDatabaseSqls())
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("After prepare database.")

})

var _ = ginkgo.AfterSuite(func() {
	// init db
	ginkgo.By("begin clean database.")
	cfg = config.TestConfig
	mDB, err = util.Open(cfg.MasterHost, cfg.MasterPort, cfg.MasterUser, cfg.MasterPassword, 10, 3)
	defer mDB.Close()
	gomega.Expect(err).Should(gomega.BeNil())

	// delete cluster data
	_, err = mDB.Exec(cfg.GenDropDatabaseSqls())
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("after clean database.")
})

var _ = ginkgo.BeforeEach(func() {
	// load config file
	//TODO : use etcd
	//start gaea
	ginkgo.By("start gaea default.")
	err = util.StartGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	// TODO: assure gaea start
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea default success.")
})

var _ = ginkgo.AfterEach(func() {
	ginkgo.By("stop gaea default.")
	_ = util.StopGaeaDefault()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.By("stop gaea default success.")
})
