package function

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Show Sql Type", func() {
	nsTemplateFile := "e2e/function/ns/default.template"
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMSName]
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())

		master, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())

		_, err = master.Exec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, db))
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.Context("show variables like", func() {
		ginkgo.It("show variables like to master or slave", func() {
			gaeaConns, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())
			sqlCase := []struct {
				sql         string
				gaeaConns   []*sql.DB
				expectAddrs []string
			}{
				{
					sql:         `show variables like "read_only"`,
					gaeaConns:   []*sql.DB{gaeaConns},
					expectAddrs: []string{slice.Slices[0].Master},
				},
			}

			for _, sqlCase := range sqlCase {
				for i := 0; i < len(sqlCase.expectAddrs); i++ {
					sql := fmt.Sprintf("/*i:%d*/ %s", i, sqlCase.sql)
					_, err := sqlCase.gaeaConns[i].Exec(sql)
					gomega.Expect(err).Should(gomega.BeNil())
					res, err := e2eMgr.SearchLog(sql, currentTime)
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(res).Should(gomega.HaveLen(1))
					gomega.Expect(res[0].BackendAddr).Should(gomega.Equal(sqlCase.expectAddrs[i]))
				}
			}

		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
