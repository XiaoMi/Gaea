package function

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Simple SQL Queries", func() {
	nsTemplateFile := "e2e/function/ns/default.template"
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

	ginkgo.Context("When executing basic SQL operations", func() {
		ginkgo.It("should handle SELECT operations correctly", func() {
			gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())
			mysqlConn, err := slice.GetMasterConn(0)
			gomega.Expect(err).Should(gomega.BeNil())

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
				gomega.Expect(err).Should(gomega.BeNil())
				err = checkFunc(sqlCase.MasterConn, sqlCase.CheckSQL, sqlCase.ExpectRes)
				gomega.Expect(err).Should(gomega.BeNil())

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
