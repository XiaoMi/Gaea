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

var _ = ginkgo.Describe("Load Balancing", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMasterSlaves]
	table := config.DefaultE2eTable
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.NsManager.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
	})

	ginkgo.Context("When distributing queries among replicas", func() {
		ginkgo.It("should evenly distribute read queries", func() {

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			util.ExpectNoError(err)

			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)

			// 定义 SQL 测试用例
			sqlCases := []struct {
				GaeaConn          *sql.DB
				GaeaSQL           string
				ExpectBackendAddr string
			}{
				//读用户-读
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case1 */SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case2 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case3 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case4 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				// 写用户-读
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case5 */SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case6 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case7 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case8 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				// 写用户-写
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case9 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case10 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case11 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case12 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
			}
			// 执行 SQL 测试用例
			for _, sqlCase := range sqlCases {
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				res, err := e2eMgr.SearchLog(sqlCase.GaeaSQL, currentTime)
				util.ExpectNoError(err)
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
