package unshard

import (
	"path/filepath"

	config "github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("sql_support_test", func() {
	planManagers := []*config.PlanManager{}
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSingleSlave]
	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)

	ginkgo.BeforeEach(func() {
		// 注册
		ns, err := config.ParseNamespaceTmpl(config.UnShardDMLNamespaceTmpl, slice)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err)

		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "unshard/case/dml"))
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
