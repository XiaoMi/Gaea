package unshard

import (
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Unshard DML Support Test", func() {
	nsTemplateFile := "unshard/ns/simple.template"
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceMaster]
	testCases := struct {
		PrepareSqlFiles []string
		TestSqlFiles    []string
		CleanSqlFiles   []string
	}{
		PrepareSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/0-prepare.sql"),
		},
		TestSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/test_join_unsharding.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/test_join.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/test_select_global_old.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/test_simple.sql"),
			filepath.Join(e2eMgr.BasePath, "unshard/case/test_subquery_global.sql"),
		},
		CleanSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "unshard/case/1-clean.sql"),
		},
	}
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
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

	})
	ginkgo.Context("unshard support test", func() {
		ginkgo.It("When testing unshard sql support", func() {
			ginkgo.By("get Gaea conn ")
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)

			masterConn, err := slice.GetMasterConn(0)
			util.ExpectNoError(err)

			for _, file := range testCases.TestSqlFiles {
				sqls, err := util.GetSqlFromFile(file)
				util.ExpectNoError(err)
				for _, sql := range sqls {
					_, err = gaeaConn.Exec("use sbtest1")
					util.ExpectNoError(err)
					_, err = masterConn.Exec("use sbtest1")
					util.ExpectNoError(err)
					err = util.CompareQueryRows(masterConn, gaeaConn, sql)
					util.ExpectNoError(err, sql)
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		// 删除注册
		e2eMgr.Clean()
	})
})
