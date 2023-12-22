package function

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Show Sql Type", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMasterSlaves]
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.NsManager.ModifyNamespace(initNs)
		util.ExpectNoError(err)

		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)

		_, err = masterAdminConn.Exec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, db))
		util.ExpectNoError(err)
	})

	ginkgo.Context("show variables like", func() {
		ginkgo.It("show variables like to master or slave", func() {
			gaeaConns, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
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
					util.ExpectNoError(err)
					res, err := e2eMgr.SearchLog(sql, currentTime)
					util.ExpectNoError(err)
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
