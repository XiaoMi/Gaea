package dml

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("test dml set variables", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		initNs.Name = "test_dml_variables"
		for i := 0; i < len(initNs.Users); i++ {
			initNs.Users[i].Namespace = initNs.Name
		}
		for i := 0; i < len(initNs.Slices); i++ {
			initNs.Slices[i].Capacity = 1
			initNs.Slices[i].MaxCapacity = 3
		}

		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err, "create namespace")
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err, "get master admin conn")
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err, "setup database and insert data")
	})
	ginkgo.Context("test set session variables take effect", func() {
		gaeaRWConn, err := e2eMgr.GetReadWriteGaeaUserConn()
		gaeaRWConn.SetMaxOpenConns(1)
		gaeaRWConn.SetMaxIdleConns(1)
		util.ExpectNoError(err, "get gaea read write conn")
		ginkgo.It("set sql_mode", func() {
			testExecCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "STRICT_TRANS_TABLES",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "abcdefghijklmnopqrstuvwxyz"),
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "abcdefghijk"),
					},
					expectErr: []bool{true, false},
				},
			}

			for _, tt := range testExecCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Exec(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}

			testQueryCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "ANSI_QUOTES",
					sql: []string{
						`show variables like "sql_mode"`,
					},
					expectErr: []bool{true},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

		ginkgo.It("set session sql_select_limit", func() {
			testQueryCase := []struct {
				setSQL    string
				sql       []string
				expectErr []bool
				resultLen []int
			}{
				{
					setSQL: "set SQL_SELECT_LIMIT=default",
					sql: []string{
						fmt.Sprintf("select * from %s.%s", db, table),
						fmt.Sprintf("select * from %s.%s limit 1", db, table),
					},
					expectErr: []bool{false, false},
					resultLen: []int{10, 1},
				},
				{
					setSQL: "set SQL_SELECT_LIMIT=2",
					sql: []string{
						fmt.Sprintf("select * from %s.%s", db, table),
						fmt.Sprintf("select * from %s.%s limit 1", db, table),
					},
					expectErr: []bool{false, false},
					resultLen: []int{2, 1},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(tt.setSQL)
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.setSQL))
				for i := 0; i < len(tt.sql); i++ {
					res, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
					count := 0
					for res.Next() {
						count++
					}
					util.ExpectEqual(count, tt.resultLen[i], fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

		ginkgo.It("set session read write", func() {
			testExecCase := []struct {
				setSQL    string
				sql       []string
				expectErr []bool
			}{
				{
					setSQL: "set session transaction read only",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "a"),
						fmt.Sprintf("select * from %s.%s", db, table),
					},
					expectErr: []bool{true, false},
				},
				{
					setSQL: "set session transaction read write",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "b"),
						fmt.Sprintf("select * from %s.%s", db, table),
					},
					expectErr: []bool{false, false},
				},
			}

			for _, tt := range testExecCase {
				_, err := gaeaRWConn.Exec(tt.setSQL)
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.setSQL))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Exec(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s.", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}

			testQueryCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "ANSI_QUOTES",
					sql: []string{
						`show variables like "sql_mode"`,
					},
					expectErr: []bool{true},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

		ginkgo.It("set session timezone in different session, maybe use same backend connection", func() {
			sqlCases := []struct {
				GaeaSQL    string
				CheckSQL   string
				CheckTimes int
				ExpectRes  string
			}{
				{
					GaeaSQL:    `/*!40103 SET TIME_ZONE='+00:00' */`,
					CheckSQL:   `show variables like "time_zone"`,
					CheckTimes: 5,
					ExpectRes:  "+00:00",
				},
				{
					GaeaSQL:    `/*!40103 SET TIME_ZONE='+01:00' */`,
					CheckSQL:   `show variables like "time_zone"`,
					CheckTimes: 5,
					ExpectRes:  "+01:00",
				},
				{
					GaeaSQL:    `/*!40103 SET TIME_ZONE='+02:00' */`,
					CheckSQL:   `show variables like "time_zone"`,
					CheckTimes: 5,
					ExpectRes:  "+02:00",
				},
				{
					GaeaSQL:    `/*!40103 SET TIME_ZONE='+03:00' */`,
					CheckSQL:   `show variables like "time_zone"`,
					CheckTimes: 5,
					ExpectRes:  "+03:00",
				},
				{
					GaeaSQL:    `/*!40103 SET TIME_ZONE='+04:00' */`,
					CheckSQL:   `show variables like "time_zone"`,
					CheckTimes: 5,
					ExpectRes:  "+04:00",
				},
			}

			for _, sqlCase := range sqlCases {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				_, err = gaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				err = checkFunc(gaeaConn, sqlCase.CheckSQL, sqlCase.ExpectRes)
				util.ExpectNoError(err)
				for i := 0; i < sqlCase.CheckTimes; i++ {
					newGaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
					util.ExpectNoError(err)
					err = checkFunc(newGaeaConn, sqlCase.CheckSQL, "SYSTEM")
					util.ExpectNoError(err)
				}
			}
		})

		ginkgo.It("set session tx_read_only and transaction_read_only", func() {
			sqlCases := []struct {
				GaeaSQL   string
				CheckSQL  string
				ExpectRes string
			}{
				{
					GaeaSQL:   `set @@tx_read_only=off`,
					CheckSQL:  `show variables like "tx_read_only"`,
					ExpectRes: "OFF",
				},
				{
					GaeaSQL:   `set @@tx_read_only=on`,
					CheckSQL:  `show variables like "tx_read_only"`,
					ExpectRes: "ON",
				},
				{
					GaeaSQL:   `set @@transaction_read_only=off`,
					CheckSQL:  `show variables like "transaction_read_only"`,
					ExpectRes: "OFF",
				},
				{
					GaeaSQL:   `set @@transaction_read_only=on`,
					CheckSQL:  `show variables like "transaction_read_only"`,
					ExpectRes: "ON",
				},
			}

			for _, sqlCase := range sqlCases {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				_, err = gaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				err = checkFunc(gaeaConn, sqlCase.CheckSQL, sqlCase.ExpectRes)
				util.ExpectNoError(err)
			}
		})

	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})

func checkFunc(db *sql.DB, sqlStr string, value string) error {
	rows, err := db.Query(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("db Exec Error %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var variableName string
		var variableValue string
		err = rows.Scan(&variableName, &variableValue)
		if value == variableValue {
			return nil
		}
		return fmt.Errorf("mismatch. Actual: %v, Expect: %v", variableValue, value)
	}

	return nil
}
