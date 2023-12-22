package shard

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"path/filepath"
)

var _ = ginkgo.Describe("shard join support test in mycat mod", func() {
	nsTemplateFile := "shard/ns/mycat/mod.template"
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
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/mod/equal.sql"),
					util.Equal,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/mod/unequal.sql"),
					util.UnEqual,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				util.ExpectNoError(err)
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
