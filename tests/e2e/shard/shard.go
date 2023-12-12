package shard

import (
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("shard_dml_support_test", func() {
	nsTemplateFile := "e2e/shard/ns/simple/simple.template"
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceMMName]
	planManagers := []*config.PlanManager{}
	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	gomega.Expect(err).Should(gomega.BeNil())

	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), sliceMulti)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init case path")
		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "e2e/shard/case/dml"))
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get sql plan")
		clusterConn, err := sliceMulti.GetLocalSliceConn()
		gomega.Expect(err).Should(gomega.BeNil())

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
				gomega.Expect(err).Should(gomega.BeNil())
				err := p.Run()
				gomega.Expect(err).Should(gomega.BeNil())
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
