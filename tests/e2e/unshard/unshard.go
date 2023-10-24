package noshard

import (
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("unshard_dml_support_test", func() {

	var e2eConfig *config.E2eConfig
	var singleMasterCluster *config.MysqlClusterConfig

	ginkgo.BeforeEach(func() {

		e2eConfig = config.GetDefaultE2eConfig()
		singleMasterCluster = e2eConfig.SingleMasterCluster

		// 解析模版
		ns, err := singleMasterCluster.TemplateParse(e2eConfig.FilepathJoin("e2e/unshard/ns/simple.template"))
		gomega.Expect(err).Should(gomega.BeNil())

		// 注册namespace
		err = e2eConfig.RegisterNamespaces(ns)
		gomega.Expect(err).Should(gomega.BeNil())

		// 初始化数据库实例
		mysqlConn, err := config.InitConn(singleMasterCluster.Slices[0].UserName, singleMasterCluster.Slices[0].Password, singleMasterCluster.Slices[0].Master, "")
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get setup file")
		sqls, err := util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/unshard/case/0-prepare.sql"))
		gomega.Expect(err).Should(gomega.BeNil())

		for _, v := range sqls {
			err := util.VerifySqlParsable(v)
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.MysqlExec(mysqlConn, v)
			gomega.Expect(err).Should(gomega.BeNil())
		}
	})
	ginkgo.Context("unshard support test", func() {
		ginkgo.It("When testing sql support", func() {

			ginkgo.By("get Gaea conn ")
			gaeaConn, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
			gomega.Expect(err).Should(gomega.BeNil())

			mysqlConn, err := config.InitConn(singleMasterCluster.Slices[0].UserName, singleMasterCluster.Slices[0].Password, singleMasterCluster.Slices[0].Master, "")
			gomega.Expect(err).Should(gomega.BeNil())

			_, err = gaeaConn.Exec("use sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())
			defer gaeaConn.Close()

			_, err = mysqlConn.Exec("use sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())
			defer mysqlConn.Close()

			casesPath := []string{
				e2eConfig.FilepathJoin("e2e/unshard/case/test_join_unsharding.sql"),
				e2eConfig.FilepathJoin("e2e/unshard/case/test_join.sql"),
				e2eConfig.FilepathJoin("e2e/unshard/case/test_select_global_old.sql"),
				e2eConfig.FilepathJoin("e2e/unshard/case/test_simple.sql"),
				e2eConfig.FilepathJoin("e2e/unshard/case/test_subquery_global.sql"),
			}
			for _, file := range casesPath {
				sqls, err := util.GetSqlFromFile(file)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, v := range sqls {
					// parser
					err := util.VerifySqlParsable(v)
					gomega.Expect(err).Should(gomega.BeNil())
					// test
					res1, err := util.MysqlQuery(mysqlConn, v)
					gomega.Expect(err).Should(gomega.BeNil())
					// gaea
					res2, err := util.MysqlQuery(gaeaConn, v)
					if err != nil {
						util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaError:%v\n", v, res1, err), e2eConfig.FilepathJoin("e2e/unshard/result/gaea_unsupport.sql"))
						continue
					}
					// compare
					_, err = util.Compare(res1, res2)
					if err != nil {
						util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaRes:%v\n", v, res1, res2), e2eConfig.FilepathJoin("e2e/unshard/result/unequal.sql"))
					} else {
						util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaRes:%v\n", v, res1, res2), e2eConfig.FilepathJoin("e2e/unshard/result/equal.sql"))
					}
				}
			}
			// 清理数据库
			sqls, err := util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/unshard/case/1-clean.sql"))
			gomega.Expect(err).Should(gomega.BeNil())
			for _, v := range sqls {
				err := util.VerifySqlParsable(v)
				gomega.Expect(err).Should(gomega.BeNil())
				_, err = util.MysqlQuery(mysqlConn, v)
				gomega.Expect(err).Should(gomega.BeNil())
			}

		})
	})
	ginkgo.AfterEach(func() {
		// 删除注册
		err := e2eConfig.UnRegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())
	})
})
