package shard

import (
	"database/sql"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("shard join support test in mycat string", func() {
	nsTemplateFile := "shard/ns/mycat/string.template"
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceMultiMaster]
	sliceSingle := e2eMgr.NsSlices[config.SliceMaster]
	multiMasterA, err := sliceMulti.GetMasterAdminConn(0)
	util.ExpectNoError(err)

	multiMasterB, err := sliceMulti.GetMasterAdminConn(1)
	util.ExpectNoError(err)

	singleMaster, err := sliceSingle.GetMasterAdminConn(0)
	util.ExpectNoError(err)

	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)
	ginkgo.BeforeEach(func() {

		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), sliceMulti)
		util.ExpectNoError(err)
		_, err = gaeaConn.Exec("USE sbtest1")
		util.ExpectNoError(err)

		prepareCases := []struct {
			DB   *sql.DB
			file string
		}{
			{
				DB:   multiMasterA,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-slice0-mycat.sql"),
			},
			{
				DB:   multiMasterB,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-slice1-mycat.sql"),
			},
			{
				DB:   singleMaster,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-test-prepare.sql"),
			},
			{
				DB:   gaeaConn,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-gaea-prepare.sql"),
			},
		}
		for _, v := range prepareCases {
			sqls, err := util.GetSqlFromFile(v.file)
			util.ExpectNoError(err)
			for _, sql := range sqls {
				err := util.VerifySqlParsable(sql)
				util.ExpectNoError(err)
				_, err = util.MysqlExec(v.DB, sql)
				util.ExpectNoError(err)
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
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/string/equal.sql"),
					util.Equal,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/string/unequal.sql"),
					util.UnEqual,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/string/unsupport.sql"),
					util.UnSupport,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				util.ExpectNoError(err)
				for _, sql := range sqls {
					err := util.VerifySqlParsable(sql)
					util.ExpectNoError(err)
					switch c.resultType {
					case util.Equal:
						res1, err := util.MysqlQuery(gaeaConn, sql)
						util.ExpectNoError(err)
						res2, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err)
						_, err = util.CompareIgnoreSort(res1, res2)
						util.ExpectNoError(err)
					case util.UnSupport:
						_, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err)
						_, err = util.MysqlQuery(gaeaConn, sql)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					case util.UnEqual:
						res1, err := util.MysqlQuery(gaeaConn, sql)
						util.ExpectNoError(err)
						res2, err := util.MysqlQuery(singleMaster, sql)
						util.ExpectNoError(err)
						_, err = util.CompareIgnoreSort(res1, res2)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					}
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
