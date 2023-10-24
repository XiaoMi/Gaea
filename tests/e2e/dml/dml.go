package dml

import (
	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sql_support_test", func() {

	var planManagers []*config.PlanManager
	var e2eConfig *config.E2eConfig
	var singleMasterCluster *config.MysqlClusterConfig
	var err error

	ginkgo.BeforeEach(func() {
		e2eConfig = config.GetDefaultE2eConfig()
		singleMasterCluster = e2eConfig.SingleMasterCluster
		// 解析模版
		ns, err := singleMasterCluster.TemplateParse(e2eConfig.FilepathJoin("e2e/dml/ns/simple.template"))
		gomega.Expect(err).Should(gomega.BeNil())

		// 注册namespace
		err = e2eConfig.RegisterNamespaces(ns)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Cluster conn ")
		err = singleMasterCluster.InitMysqlClusterConn()
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Gaea conn ")
		gaeaConn, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init case path")

		casesPath, err := config.GetJSONFilesFromDir(e2eConfig.FilepathJoin("e2e/dml/case"))
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := &config.PlanManager{
				PlanPath:         v,
				MysqlClusterConn: singleMasterCluster.SliceConns,
				GaeaDB:           gaeaConn,
			}
			planManagers = append(planManagers, p)
		}
	})

	ginkgo.Context("sql support test", func() {
		ginkgo.It("sql support", func() {
			for _, p := range planManagers {
				err = p.Init()
				gomega.Expect(err).Should(gomega.BeNil())
				err := p.Run()
				gomega.Expect(err).Should(gomega.BeNil())
			}
		})
	})
	ginkgo.AfterEach(func() {
		for _, p := range planManagers {
			p.MysqlClusterConnClose()
		}
		// 删除注册
		e2eConfig.NamespaceManager.UnRegisterNamespaces()
	})
})
