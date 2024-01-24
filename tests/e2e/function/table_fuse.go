package function

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// This script, titled "Test HealthyCheck Sql Fuse Result" is designed to test the slave database failover (fuse) functionality in a database environment.
// The test involves setting up a database and a table, followed by executing various SQL queries to check how operations are handled in a setup with a single slave slice.
// The test includes different SQL cases for scenarios such as non-existent databases, non-existent tables, and table locking.
// It ensures that traffic is correctly switched to the master database when specific health check queries fail against the slave.
// Each test case executes preparation SQL statements (like creating and locking tables) and then runs a health check query,
// followed by validating if the actual backend address corresponds with the expected one (i.e., correctly switching to the master or staying on the slave).
// After executing the test cases, any data generated during the test is cleaned up.
// Upon completion, the script checks the results to ensure they align with the expected behavior of slave database failover.
var _ = ginkgo.Describe("Test HealthyCheck Sql Fuse Result", func() {
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSingleSlave]
	db := config.DefaultE2eDatabase
	table := config.DefaultE2eTable
	var initNs *models.Namespace

	ginkgo.BeforeEach(func() {
		// mysql prepare
		master, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		util.ExpectNoError(util.CleanUpDatabases(master))
		err = util.SetupDatabaseAndInsertData(master, db, table)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(master, "db_fuse", "tbl_fuse")
		util.ExpectNoError(err)
		// namespace prepare
		initNs, err = config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		initNs.DownAfterNoAlive = 1
		initNs.Slices[0].HealthCheckSql = "select * from db_fuse.tbl_fuse"
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait namespace down after no alive effect
		time.Sleep(10 * time.Second)
	})

	ginkgo.Context("When executing basic SQL operations", func() {
		ginkgo.It("should handle SELECT operations correctly", func() {
			master, err := slice.GetMasterAdminConn(0)
			util.ExpectNoError(err)

			slave, err := slice.GetSlaveAdminConn(0, 0)
			util.ExpectNoError(err)

			sqlCases := []struct {
				DB                *sql.DB
				PrePareSql        string
				ExpectBackendAddr string
				CleanSql          string
			}{
				{
					//从库锁表，
					DB: slave,
					// 从库锁表/主库正常/health_check_query语法正确/执行超时/熔断从库，主库正常
					PrePareSql:        "LOCK TABLES db_fuse.tbl_fuse WRITE;",
					ExpectBackendAddr: slice.Slices[0].Master,
					CleanSql:          "UNLOCK TABLES;",
				},
				{
					//主库删除表，从库同步删除表，
					DB: master,
					// 从库表不存在/health_check_query语法正确/配置错误/不熔断从库，主库。
					PrePareSql:        "DROP TABLE  db_fuse.tbl_fuse",
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
				{
					//主库删除库，从库同步删除库
					DB: master,
					// 从库库不存在/health_check_query语法正确/配置错误/不熔断从库，主库。
					PrePareSql:        "DROP DATABASE  db_fuse",
					ExpectBackendAddr: slice.Slices[0].Slaves[0],
				},
			}

			for i, v := range sqlCases {
				if len(v.PrePareSql) > 0 {
					_, err = v.DB.Exec(v.PrePareSql)
					util.ExpectNoError(err, "Prepare Sql Exec Error")
				}
				// 1. wait fuse  effect
				time.Sleep(10 * time.Second)
				// 2. gaea exec sql
				currentTime := time.Now()
				gaea, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				sql := fmt.Sprintf("/* No.%d */ select * from %s.%s", i, db, table)
				_, err = gaea.Exec(sql)
				util.ExpectNoError(err)
				// 3.wait log  write
				time.Sleep(1 * time.Second)
				// 4. search log
				res, err := e2eMgr.SearchSqlLog(sql, currentTime)
				util.ExpectNoError(err)
				// 5. compare res and expect
				gomega.Expect(res).Should(gomega.HaveLen(1))
				gomega.Expect(v.ExpectBackendAddr).Should(gomega.Equal(res[0].BackendAddr))

				if len(v.CleanSql) > 0 {
					_, err = v.DB.Exec(v.CleanSql)
					util.ExpectNoError(err, "Clean Sql Exec Error")
				}

			}

		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
