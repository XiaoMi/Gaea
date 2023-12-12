package shard

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// 定义输出结果的方法
func OutputEqual(sql string, mysqlRes, gaeaRes interface{}, err error, basePath string) {
	util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaRes:%v\n", sql, mysqlRes, gaeaRes), filepath.Join(basePath, "e2e/unshard/result/equal.sql"))
}

func OutputUnEqual(sql string, mysqlRes, gaeaRes interface{}, err error, basePath string) {
	util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaRes:%v\n", sql, mysqlRes, gaeaRes), filepath.Join(basePath, "e2e/unshard/result/unequal.sql"))
}

func OutputUnSupport(sql string, mysqlRes, gaeaRes interface{}, err error, basePath string) {
	util.OutPutResult(fmt.Sprintf("sql:%s\nmysqlRes:%v\ngaeaError:%v\n", sql, mysqlRes, err), filepath.Join(basePath, "e2e/unshard/result/gaea_unsupport.sql"))
}

var _ = ginkgo.Describe("shard_join_support_test", func() {
	nsTemplateFile := "e2e/shard/ns/kingshard/hash.template"
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceMMName]
	sliceSingle := e2eMgr.NsSlices[config.SliceSMName]
	multiMasterA, err := sliceMulti.GetMasterConn(0)
	gomega.Expect(err).Should(gomega.BeNil())

	multiMasterB, err := sliceMulti.GetMasterConn(1)
	gomega.Expect(err).Should(gomega.BeNil())

	singleMaster, err := sliceSingle.GetMasterConn(0)
	gomega.Expect(err).Should(gomega.BeNil())

	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.BeforeEach(func() {

		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), sliceMulti)
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = gaeaConn.Exec("USE sbtest1")
		gomega.Expect(err).Should(gomega.BeNil())

		prepareCases := []struct {
			DB   *sql.DB
			file string
		}{
			{
				DB:   multiMasterA,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice0-kingshard.sql"),
			},
			{
				DB:   multiMasterB,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice1-kingshard.sql"),
			},
			{
				DB:   singleMaster,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-test-prepare.sql"),
			},
			{
				DB:   gaeaConn,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-gaea-prepare.sql"),
			},
		}
		for _, v := range prepareCases {
			sqls, err := util.GetSqlFromFile(v.file)
			gomega.Expect(err).Should(gomega.BeNil())
			for _, sql := range sqls {
				err := util.VerifySqlParsable(sql)
				gomega.Expect(err).Should(gomega.BeNil())
				_, err = util.MysqlExec(v.DB, sql)
				gomega.Expect(err).Should(gomega.BeNil())
			}
		}

	})

	ginkgo.Context("shard support test", func() {
		ginkgo.It("When testing sql support", func() {
			files := []string{
				filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/test_join.sql"),
				filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/test_select_global_old.sql"),
				filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/test_simple.sql"),
				filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/test_subquery_global.sql"),
			}
			for _, file := range files {
				sqls, err := util.GetSqlFromFile(file)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, sql := range sqls {
					err := util.VerifySqlParsable(sql)
					gomega.Expect(err).Should(gomega.BeNil())
					res1, err := util.MysqlQuery(gaeaConn, sql)
					if err != nil {
						OutputUnSupport(sql, res1, nil, err, e2eMgr.BasePath)
					}
					res2, err := util.MysqlQuery(singleMaster, sql)
					gomega.Expect(err).Should(gomega.BeNil())
					_, err = util.Compare(res1, res2)
					if err != nil {
						OutputUnEqual(sql, res1, res2, err, e2eMgr.BasePath)
					} else {
						OutputEqual(sql, res1, res2, err, e2eMgr.BasePath)
					}

				}
			}
		})

	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
