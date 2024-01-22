package unshard

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Unshard DML Support Test", func() {
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSingleSlave]
	testCases := struct {
		PrepareSqlFiles []string
		TestSqlFiles    []string
		CleanSqlFiles   []string
	}{
		PrepareSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/0-prepare.sql"),
		},
		TestSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/test_join_unsharding.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/test_join.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/test_select_global_old.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/test_simple.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/test_subquery_global.sql"),
		},
		CleanSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/join/1-clean.sql"),
		},
	}
	ginkgo.BeforeEach(func() {
		ns, err := config.ParseNamespaceTmpl(config.UnShardNamespaceTmpl, slice)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err)
		master, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		ginkgo.By("get setup file")
		for _, file := range testCases.PrepareSqlFiles {
			sqls, err := util.GetSqlFromFile(file)
			util.ExpectNoError(err)
			for _, sql := range sqls {
				_, err = util.MysqlExec(master, sql)
				util.ExpectNoError(err)
			}
		}
		// 等待从库追赶主库
		time.Sleep(2 * time.Second)
	})
	ginkgo.Context("unshard support test", func() {
		ginkgo.It("When testing unshard sql support", func() {
			ginkgo.By("get Gaea conn ")
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
			// namespace使用的是普通权限的用户
			// 简单show语句走从库
			slaveConn, err := slice.GetSlaveCommonConn(0, 0)
			util.ExpectNoError(err)
			for _, file := range testCases.TestSqlFiles {
				sqls, err := util.GetSqlFromFile(file)
				util.ExpectNoError(err)
				for _, sql := range sqls {
					_, err = gaeaConn.Exec("use sbtest1")
					util.ExpectNoError(err)
					_, err = slaveConn.Exec("use sbtest1")
					util.ExpectNoError(err)
					gaeaRes, err := util.MysqlQuery(gaeaConn, sql)
					util.ExpectNoError(err, fmt.Sprintf("gaea exec equal sql err:%v, err:%v\n", sql, err))
					mysqlRes, err := util.MysqlQuery(slaveConn, sql)
					util.ExpectNoError(err, fmt.Sprintf("mysql exec equal sql err:%v, err:%v\n", sql, err))
					res, err := util.CompareIgnoreSort(gaeaRes, mysqlRes)
					util.ExpectNoError(err, fmt.Sprintf("compare equal sql get uneuqal res.sql:%v, err:%v\n", sql, err), res)
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		// 删除注册
		e2eMgr.Clean()
	})
})
