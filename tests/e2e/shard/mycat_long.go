package shard

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("shard join support test in mycat long", func() {
	e2eMgr := config.NewE2eManager()
	sliceTest := e2eMgr.NsSlices[config.SliceSingleTestMaster]
	sliceMulti := e2eMgr.NsSlices[config.SliceDualMaster]
	ginkgo.BeforeEach(func() {
		// 注册
		ns, err := config.ParseNamespaceTmpl(config.MycatLongNamespaceTmpl, sliceMulti)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err)

		// AdminConn 可以DROP/CREATE
		multiMasterA, err := sliceMulti.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		util.ExpectNoError(util.CleanUpDatabases(multiMasterA))

		// AdminConn 可以DROP/CREATE
		multiMasterB, err := sliceMulti.GetMasterAdminConn(1)
		util.ExpectNoError(err)
		util.ExpectNoError(util.CleanUpDatabases(multiMasterB))
		// AdminConn 可以DROP/CREATE
		singleMaster, err := sliceTest.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		util.ExpectNoError(util.CleanUpDatabases(singleMaster))
		// 获取gaea连接
		gaeaConn, err := e2eMgr.GetReadWriteGaeaUserDBConn("sbtest")
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
				_, err = util.MysqlExec(v.DB, sql)
				util.ExpectNoError(err)
			}
		}

	})

	ginkgo.Context("shard support test", func() {
		ginkgo.It("When testing sql support", func() {
			// gaea使用的ns中的用户是普通权限的用户
			singleMaster, err := sliceTest.GetMasterCommonDBConn(0, "sbtest")
			util.ExpectNoError(err)
			defaultShardConn, err := sliceMulti.GetMasterCommonDBConn(0, "sbtest_0")
			util.ExpectNoError(err)
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserDBConn("sbtest")
			util.ExpectNoError(err)

			cases := []struct {
				path           string
				resultType     util.ResultType
				gaeaConn       *sql.DB
				comparisonConn *sql.DB
			}{
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/long/equal.sql"),
					util.Equal,
					gaeaConn,
					singleMaster,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/long/show.sql"),
					util.Equal,
					gaeaConn,
					defaultShardConn,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/long/unequal.sql"),
					util.UnEqual,
					gaeaConn,
					singleMaster,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/mycat/long/unsupport.sql"),
					util.UnSupport,
					gaeaConn,
					singleMaster,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				util.ExpectNoError(err)
				for _, sql := range sqls {
					switch c.resultType {
					case util.Equal:
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec equal sql err:%v, err:%v\n", sql, err))
						mysqlRes, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec equal sql err:%v, err:%v\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectNoError(err, fmt.Sprintf("compare equal sql get uneuqal res.sql:%v, err:%v\n", sql, err))
					case util.UnSupport:
						_, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unsupport sql err:%v, err:%v\n", sql, err))
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectError(err, fmt.Sprintf("gaea exec unsupport sql get no error, sql:%v,res:%v\n", sql, gaeaRes))
					case util.UnEqual:
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec unequal sql err:%v, err:%v\n", sql, err))
						mysqlRes, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unequal sql err.sql:%v, err:%v\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectError(err, fmt.Sprintf("compare unequal sql get equal res.sql:%v, err:%v\n", sql, err), gaeaRes, mysqlRes)
					}
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
