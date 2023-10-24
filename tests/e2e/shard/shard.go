package shard

import (
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("shard_dml_support_test", func() {
	var e2eConfig *config.E2eConfig
	var multiMasterCluster *config.MysqlClusterConfig
	var planManagers []*config.PlanManager
	var err error

	ginkgo.BeforeEach(func() {
		e2eConfig = config.GetDefaultE2eConfig()
		multiMasterCluster = e2eConfig.MultiMasterCluster
		// 解析模版
		ns, err := multiMasterCluster.TemplateParse(e2eConfig.FilepathJoin("e2e/shard/ns/simple/simple.template"))
		gomega.Expect(err).Should(gomega.BeNil())
		// 注册namespace
		err = e2eConfig.RegisterNamespaces(ns)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Gaea conn ")
		gaeaDB, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init case path")
		casesPath, err := config.GetJSONFilesFromDir(e2eConfig.FilepathJoin("e2e/shard/case/dml"))
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Cluster conn ")
		err = multiMasterCluster.InitMysqlClusterConn()
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := &config.PlanManager{
				PlanPath:         v,
				MysqlClusterConn: multiMasterCluster.SliceConns,
				GaeaDB:           gaeaDB,
			}
			ginkgo.By(fmt.Sprintf("plan %s init", v))
			planManagers = append(planManagers, p)
		}

	})
	ginkgo.Context("shard support test", func() {
		ginkgo.It("shard support", func() {
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
		err := e2eConfig.UnRegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())

	})
})
