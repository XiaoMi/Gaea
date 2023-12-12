package function

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Client Connection Limit", func() {
	nsTemplateFile := "e2e/function/ns/limit.template"
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMSName]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())
		masterConn, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())
		err = util.SetupDatabaseAndInsertData(masterConn, db, table)
		gomega.Expect(err).Should(gomega.BeNil())
	})
	ginkgo.Context("When handling client limit operations", func() {
		ginkgo.It("should limit exceeded maximum connection ", func() {
			var globalConns = []*sql.DB{}
			for i := 0; i < 11; i++ {
				gaeaReadWriteConn, err := e2eMgr.NewReadWriteGaeaUserConn()
				gomega.Expect(err).Should(gomega.BeNil())
				globalConns = append(globalConns, gaeaReadWriteConn)
				_, err = globalConns[i].Exec(fmt.Sprintf("Use %s", db))
				if i < 10 {
					gomega.Expect(err).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				}
				_, err = globalConns[i].Exec(fmt.Sprintf("select * from %s", table))
				if i < 10 {
					gomega.Expect(err).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				}
			}
			for _, v := range globalConns {
				v.Close()
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
