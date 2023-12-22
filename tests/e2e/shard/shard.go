package shard

import (
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("shard_dml_support_test", func() {
	nsTemplateFile := "shard/ns/simple/simple.template"
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceMultiMaster]
	planManagers := []*config.PlanManager{}
	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)

	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), sliceMulti)
		util.ExpectNoError(err)

		ginkgo.By("init case path")
		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "shard/case/dml"))
		util.ExpectNoError(err)

		ginkgo.By("get sql plan")
		clusterConn, err := sliceMulti.GetLocalSliceConn()
		util.ExpectNoError(err)

		for _, v := range casesPath {
			p := &config.PlanManager{
				PlanPath:         v,
				MysqlClusterConn: clusterConn,
				GaeaDB:           gaeaConn,
			}
			ginkgo.By(fmt.Sprintf("plan %s init", v))
			planManagers = append(planManagers, p)
		}

	})
	ginkgo.Context("shard support test", func() {
		ginkgo.It("shard support", func() {
			for _, p := range planManagers {
				err = p.Init()
				util.ExpectNoError(err)
				err := p.Run()
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
