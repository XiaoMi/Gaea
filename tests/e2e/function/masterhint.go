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

// This script, titled "Force Read from Master" aims to validate that all read operations are correctly directed to the master database, regardless of whether the user connection is configured for read, write, or both.
// This is tested using various SQL queries with a "/master/" hint.
// The script checks if the queries are being executed on the master or slave nodes as expected, based on the type of connection (read, write, or read-write) and the presence of the master hint in the query.
// The results from the SQL log are then verified to ensure that the read queries are directed to the correct backend address as per the configuration, demonstrating the effectiveness of the forced master read functionality.
var _ = ginkgo.Describe("Force Read from Master", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
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

	ginkgo.Context("When all read operations are forced to master", func() {
		ginkgo.It("should evenly distribute read queries", func() {

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			util.ExpectNoError(err)
			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
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
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/ * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[1],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/  * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: slice.Slices[0].Slaves[1],
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 5", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 6", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * /*master*/  FROM `%s`.`%s` WHERE `id`= 5", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT /*master*/ * FROM `%s`.`%s` WHERE `id`= 6", db, table),
					ExpectBackendAddr: slice.Slices[0].Master,
				},
			}
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
