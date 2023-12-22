package dml

import (
	config "github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"path/filepath"
)

var _ = ginkgo.Describe("sql_support_test", func() {
	planManagers := []*config.PlanManager{}
	nsTemplateFile := "dml/ns/simple.template"
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceMaster]

	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)

	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), e2eMgr.NsSlices[config.SliceMaster])
		util.ExpectNoError(err)

		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "dml/case"))
		util.ExpectNoError(err)

		conns, err := slice.GetLocalSliceConn()
		util.ExpectNoError(err)

		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := &config.PlanManager{
				PlanPath:         v,
				MysqlClusterConn: conns,
				GaeaDB:           gaeaConn,
			}
			planManagers = append(planManagers, p)
		}
	})

	ginkgo.Context("sql support test", func() {
		ginkgo.It("sql support", func() {
			for _, p := range planManagers {
				err := p.Init()
				util.ExpectNoError(err)
				err = p.Run()
				util.ExpectNoError(err)
			}
		})
	})

	ginkgo.AfterEach(func() {
		for _, p := range planManagers {
			p.MysqlClusterConnClose()
		}
		e2eMgr.Clean()
	})
})
