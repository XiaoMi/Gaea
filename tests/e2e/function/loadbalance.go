package function

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// This script, titled "Load Balancing," aims to test the load balancing functionality in a database environment.
// The test involves setting up a database and table, then running various SQL queries to check how read and write operations are distributed among the master and slave nodes in a dual-slave slice setup.
// The test includes different SQL cases for read and write operations, ensuring that read queries are evenly distributed between the slave nodes while write queries go to the master node.
// Each query's execution location is verified against the expected backend address.
// After running the test cases, the script checks the results to ensure they align with the expected load balancing behavior..
var _ = ginkgo.Describe("Load Balancing", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
		// namespace prepare
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(3 * time.Millisecond)
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
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case2 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[1],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case3 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case4 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[1],
				},
				// 写用户-读
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case5 */SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case6 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case7 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case8 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				// 写用户-写
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case9 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case10 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case11 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case12 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
			}
			e2eMgr.ClearSqlLog()
			currentTime := time.Now()
			// 执行 SQL 测试用例
			for _, sqlCase := range sqlCases {
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				res, err := e2eMgr.SearchSqlLog(sqlCase.GaeaSQL, currentTime)
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
