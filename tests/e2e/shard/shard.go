package shard

import (
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("shard_dml_support_test", func() {
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceDualMaster]
	planManagers := []*config.PlanManager{}
	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)

	ginkgo.BeforeEach(func() {
		// 注册
		ns, err := config.ParseNamespaceTmpl(config.ShardNamespaceTmpl, sliceMulti)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
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
