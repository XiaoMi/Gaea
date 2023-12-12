package dml

import (
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sql_support_test", func() {
	planManagers := []*config.PlanManager{}
	nsTemplateFile := "e2e/dml/ns/simple.template"
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSMName]

	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	gomega.Expect(err).Should(gomega.BeNil())

	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), e2eMgr.NsSlices[config.SliceSMName])
		gomega.Expect(err).Should(gomega.BeNil())

		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "e2e/dml/case"))
		gomega.Expect(err).Should(gomega.BeNil())

		conns, err := slice.GetLocalSliceConn()
		gomega.Expect(err).Should(gomega.BeNil())

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
				gomega.Expect(err).Should(gomega.BeNil())
				err = p.Run()
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
