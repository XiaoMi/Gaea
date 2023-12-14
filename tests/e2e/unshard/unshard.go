package unshard

import (
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Unshard DML Support Test", func() {
	nsTemplateFile := "e2e/unshard/ns/simple.template"
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSMName]
	testCases := struct {
		PrepareSqlFiles []string
		TestSqlFiles    []string
		CleanSqlFiles   []string
	}{
		PrepareSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/0-prepare.sql"),
		},
		TestSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/test_join_unsharding.sql"),
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/test_join.sql"),
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/test_select_global_old.sql"),
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/test_simple.sql"),
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/test_subquery_global.sql"),
		},
		CleanSqlFiles: []string{
			filepath.Join(e2eMgr.BasePath, "e2e/unshard/case/1-clean.sql"),
		},
	}
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())
		master, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())
		ginkgo.By("get setup file")
		for _, file := range testCases.PrepareSqlFiles {
			sqls, err := util.GetSqlFromFile(file)
			gomega.Expect(err).Should(gomega.BeNil())
			for _, sql := range sqls {
				_, err = util.MysqlExec(master, sql)
				gomega.Expect(err).Should(gomega.BeNil())
			}
		}

	})
	ginkgo.Context("unshard support test", func() {
		ginkgo.It("When testing unshard sql support", func() {
			ginkgo.By("get Gaea conn ")
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

			masterConn, err := slice.GetMasterConn(0)
			gomega.Expect(err).Should(gomega.BeNil())

			for _, file := range testCases.TestSqlFiles {
				sqls, err := util.GetSqlFromFile(file)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, sql := range sqls {
					_, err = gaeaConn.Exec("use sbtest1")
					gomega.Expect(err).Should(gomega.BeNil())
					_, err = masterConn.Exec("use sbtest1")
					gomega.Expect(err).Should(gomega.BeNil())
					err = util.CompareQueryRows(masterConn, gaeaConn, sql)
					gomega.Expect(err).Should(gomega.BeNil())
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		// 删除注册
		e2eMgr.Clean()
	})
})
