package function

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Force Read from Master", func() {
	nsTemplateFile := "e2e/function/ns/default.template"
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMSName]
	table := config.DefaultE2eTable
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())
		masterConn, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())
		err = util.SetupDatabaseAndInsertData(masterConn, db, table)
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.Context("When all read operations are forced to master", func() {
		ginkgo.It("should evenly distribute read queries", func() {

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())
			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())
			//Select 带 Master Hint：
			//读写分离用户（RWFlag=2，RWSplit=1）：主库
			//只写用户（RWFlag=2，RWSplit=0）：主库
			//只读用户（RWFlag=1，RWSplit=1）：从库（V2.0 以下版本会请求到主库，MiProxy 会打到从库）
			sqlCases := []struct {
				GaeaConn          *sql.DB
				GaeaSQL           string
				ExpectBackendAddr string
			}{
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/ * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 5", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 6", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * /*master*/  FROM `%s`.`%s` WHERE `id`= 5", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/ * FROM `%s`.`%s` WHERE `id`= 6", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
			}
			for _, sqlCase := range sqlCases {
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				gomega.Expect(err).Should(gomega.BeNil())
				res, err := e2eMgr.SearchLog(sqlCase.GaeaSQL, currentTime)
				gomega.Expect(err).Should(gomega.BeNil())
				// 避免扫到以前的数据
				gomega.Expect(res).Should(gomega.HaveLen(1))
				gomega.Expect(sqlCase.ExpectBackendAddr).Should(gomega.Equal(res[0].BackendAddr))
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
