package shard

import (
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("shard_join_support_test", func() {
	var e2eConfig *config.E2eConfig
	var multiMasterCluster *config.MysqlClusterConfig
	var singleMasterCluster *config.MysqlClusterConfig
	ginkgo.BeforeEach(func() {

		e2eConfig = config.GetDefaultE2eConfig()
		multiMasterCluster = e2eConfig.MultiMasterCluster
		singleMasterCluster = e2eConfig.SingleMasterCluster
		// 解析模版
		ns, err := multiMasterCluster.TemplateParse(e2eConfig.FilepathJoin("e2e/shard/ns/kingshard/hash.template"))
		gomega.Expect(err).Should(gomega.BeNil())

		// 注册namespace
		err = e2eConfig.RegisterNamespaces(ns)
		gomega.Expect(err).Should(gomega.BeNil())

		multiMasterA, err := config.InitConn(multiMasterCluster.Slices[0].UserName, multiMasterCluster.Slices[0].Password, multiMasterCluster.Slices[0].Master, "")
		gomega.Expect(err).Should(gomega.BeNil())
		defer multiMasterA.Close()
		multiMasterB, err := config.InitConn(multiMasterCluster.Slices[1].UserName, multiMasterCluster.Slices[1].Password, multiMasterCluster.Slices[1].Master, "")
		gomega.Expect(err).Should(gomega.BeNil())
		defer multiMasterB.Close()

		singleMasterTest, err := config.InitConn(singleMasterCluster.Slices[0].UserName, singleMasterCluster.Slices[0].Password, singleMasterCluster.Slices[0].Master, "")
		gomega.Expect(err).Should(gomega.BeNil())
		defer singleMasterTest.Close()

		// 测试数据库执行
		sqls, err := util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/shard/case/join/0-test-prepare.sql"))
		gomega.Expect(err).Should(gomega.BeNil())
		for _, v := range sqls {
			err := util.VerifySqlParsable(v)
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.MysqlExec(singleMasterTest, v)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		// Gaea 对应的分片数据库-1
		sqls, err = util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/shard/case/join/0-slice0-kingshard.sql"))
		gomega.Expect(err).Should(gomega.BeNil())
		for _, v := range sqls {
			err := util.VerifySqlParsable(v)
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.MysqlExec(multiMasterA, v)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		// Gaea 对应的分片数据库-2
		sqls, err = util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/shard/case/join/0-slice1-kingshard.sql"))
		gomega.Expect(err).Should(gomega.BeNil())
		for _, v := range sqls {
			err := util.VerifySqlParsable(v)
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.MysqlExec(multiMasterB, v)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		gaeaDB, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = gaeaDB.Exec("USE sbtest1")
		gomega.Expect(err).Should(gomega.BeNil())
		defer gaeaDB.Close()

		sqls, err = util.GetSqlFromFile(e2eConfig.FilepathJoin("e2e/shard/case/join/0-gaea-prepare.sql"))
		gomega.Expect(err).Should(gomega.BeNil())
		for _, v := range sqls {
			err := util.VerifySqlParsable(v)
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.MysqlExec(gaeaDB, v)
			gomega.Expect(err).Should(gomega.BeNil())
		}
	})

	ginkgo.Context("shard support test", func() {
		ginkgo.It("When testing sql support", func() {
			casesPath := []string{
				e2eConfig.FilepathJoin("e2e/shard/case/join/test_join.sql"),
				e2eConfig.FilepathJoin("e2e/shard/case/join/test_select_global_old.sql"),
				e2eConfig.FilepathJoin("e2e/shard/case/join/test_simple.sql"),
				e2eConfig.FilepathJoin("e2e/shard/case/join/test_subquery_global.sql"),
			}

			gaea, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())

			test, err := config.InitConn(singleMasterCluster.Slices[0].UserName, singleMasterCluster.Slices[0].Password, singleMasterCluster.Slices[0].Master, "sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())

			for _, file := range casesPath {
				sqls, err := util.GetSqlFromFile(file)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, v := range sqls {
					// parser
					err := util.VerifySqlParsable(v)
					gomega.Expect(err).Should(gomega.BeNil())
					// test
					res1, err := util.MysqlQuery(test, v)
					gomega.Expect(err).Should(gomega.BeNil())
					// gaea
					res2, err := util.MysqlQuery(gaea, "Query")
					if err != nil {
						util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaError:%v\n", v, res1, err), e2eConfig.FilepathJoin("e2e/shard/result/gaea_unsupport.sql"))
						continue
					}
					// compare
					_, err = util.Compare(res1, res2)
					if err != nil {
						util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaRes:%v\n", v, res1, res2), e2eConfig.FilepathJoin("e2e/shard/result/unequal.sql"))
					} else {
						util.OutPutResult(v, e2eConfig.FilepathJoin("e2e/shard/result/equal.sql"))
					}
				}
			}
			err = gaea.Close()
			gomega.Expect(err).Should(gomega.BeNil())
			err = test.Close()
			gomega.Expect(err).Should(gomega.BeNil())
		})
	})
	ginkgo.AfterEach(func() {
		// 删除注册
		err := e2eConfig.UnRegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())
	})
})
