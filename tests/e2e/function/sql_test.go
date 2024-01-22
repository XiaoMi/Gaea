package function

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

// This script, titled "Simple SQL Queries" is designed to verify basic SQL operations in a database system.
// It sets up a test environment and then performs a series of SQL commands such as SELECT, DELETE, UPDATE, and INSERT.
// The test checks if these operations yield the expected results, using a custom function checkFunc to validate the outcomes.
// This approach ensures the database system accurately processes and reflects changes made by standard SQL queries, highlighting its capability to handle essential database operations reliably.
var _ = ginkgo.Describe("Simple SQL Queries", func() {
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
		time.Sleep(500 * time.Millisecond)
	})

	ginkgo.Context("When executing basic SQL operations", func() {
		ginkgo.It("should handle SELECT operations correctly", func() {
			gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)
			mysqlConn, err := slice.GetMasterAdminConn(0)
			util.ExpectNoError(err)

			// 定义 SQL 测试用例
			sqlCases := []struct {
				GaeaConn   *sql.DB
				GaeaSQL    string
				MasterConn *sql.DB
				CheckSQL   string
				ExpectRes  [][]string
			}{
				{
					GaeaConn:   gaeaConn,
					GaeaSQL:    fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 9", db, table),
					MasterConn: mysqlConn,
					CheckSQL:   fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 9", db, table),
					ExpectRes: [][]string{
						{"9", "nameValue"},
					},
				},
				{
					GaeaConn:   gaeaConn,
					GaeaSQL:    fmt.Sprintf("DELETE FROM %s.%s WHERE id=1", db, table),
					MasterConn: mysqlConn,
					CheckSQL:   fmt.Sprintf("SELECT * FROM %s.%s WHERE id=1", db, table),
					ExpectRes:  [][]string{},
				},
				{
					GaeaConn:   gaeaConn,
					GaeaSQL:    fmt.Sprintf("UPDATE %s.%s SET name='Tom' WHERE id=2", db, table),
					MasterConn: mysqlConn,
					CheckSQL:   fmt.Sprintf("SELECT * FROM %s.%s WHERE id=2 AND name='Tom'", db, table),
					ExpectRes: [][]string{{
						"2", "Tom",
					}},
				},
				{
					GaeaConn:   gaeaConn,
					GaeaSQL:    fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('Alex')", db, table),
					MasterConn: mysqlConn,
					CheckSQL:   fmt.Sprintf("SELECT * FROM %s.%s WHERE name='Alex'", db, table),
					ExpectRes: [][]string{{
						"11", "Alex",
					}},
				},
			}

			// 执行 SQL 测试用例
			for _, sqlCase := range sqlCases {
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				err = checkFunc(sqlCase.MasterConn, sqlCase.CheckSQL, sqlCase.ExpectRes)
				util.ExpectNoError(err)

			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})

func checkFunc(db *sql.DB, sqlStr string, values [][]string) error {
	rows, err := db.Query(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows && len(values) == 0 {
			return nil
		}
		return fmt.Errorf("db Exec Error %v", err)
	}
	defer rows.Close()
	res, err := util.GetDataFromRows(rows)
	if err != nil {
		return fmt.Errorf("get data from rows error:%v", err)
	}
	// res为空代表没有查到数据
	if (len(res) == 0 || res == nil) && len(values) == 0 {
		return nil
	}
	if !reflect.DeepEqual(values, res) {
		return fmt.Errorf("mismatch. Actual: %v, Expect: %v", res, values)
	}
	return nil
}