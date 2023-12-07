package shard

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"path/filepath"
)

var _ = ginkgo.Describe("shard join support test in mycat mod", func() {
	nsTemplateFile := "e2e/shard/ns/mycat/mod.template"
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
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice0-mycat.sql"),
			},
			{
				DB:   multiMasterB,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice1-mycat.sql"),
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
			cases := []struct {
				path       string
				resultType util.ResultType
			}{
				{
					filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/mycat/mod/equal.sql"),
					util.Equal,
				},
				{
					filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/mycat/mod/unequal.sql"),
					util.UnEqual,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, sql := range sqls {
					switch c.resultType {
					case util.Equal:
						gaeaRes, err := util.MysqlQuery(gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec equal sql err:%s, err:%s\n", sql, err))
						mysqlRes, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec equal sql err:%s, err:%s\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectNoError(err, fmt.Sprintf("compare equal sql get uneuqal res.sql:%s, err:%s\n", sql, err))
					case util.UnSupport:
						_, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unsupport sql err:%s, err:%s\n", sql, err))
						gaeaRes, err := util.MysqlQuery(gaeaConn, sql)
						util.ExpectError(err, fmt.Sprintf("gaea exec unsupport sql get no error, sql:%s,res:%s\n", sql, gaeaRes))
					case util.UnEqual:
						gaeaRes, err := util.MysqlQuery(gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec unequal sql err:%s, err:%s\n", sql, err))
						mysqlRes, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unequal sql err.sql:%s, err:%s\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectError(err, fmt.Sprintf("compare unequal sql get equal res.sql:%s, err:%s\n", sql, err))
					}
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
