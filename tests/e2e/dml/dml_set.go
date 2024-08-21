// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dml

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

// This Ginkgo test suite validates the correct setting and behavior of session variables within a database context,
// focusing on DML operations and session controls.
// It ensures that changes to session-specific variables,
// such as SQL modes and transaction read/write settings, are enforced correctly and affect database operations as expected.
// The tests cover various scenarios including the enforcement of SQL standards, limits on SQL select queries, and settings that govern data modification permissions within a session.
// This suite helps verify that the database handles session settings correctly, providing a robust environment for further application-specific testing.
// BeforeEach and AfterEach are defined directly inside the Describe block,
// but outside the Context block. This means that they will be executed for all It blocks inside Describe blocks,
// including those inside any Context blocks. No matter what level of Context these It blocks are located in,
// BeforeEach and AfterEach will be called before and after each It block runs.
// e2eMgr.Clean() will clean the connection of each open gaea, so make sure that each ginkgo.It uses an independent gaea connection
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
		util.ExpectNoError(err, "get gaea read write conn when test set session variables take effect")
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
			gaeaRWConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gaeaRWConn.SetMaxOpenConns(1)
			gaeaRWConn.SetMaxIdleConns(1)
			util.ExpectNoError(err, "get gaea read write conn when set session sql_select_limit")
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
			gaeaRWConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gaeaRWConn.SetMaxOpenConns(1)
			gaeaRWConn.SetMaxIdleConns(1)
			util.ExpectNoError(err, "get gaea read write conn when set session read write")
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

		ginkgo.It("set character_set_client binary test", func() {
			sqlCases := []struct {
				GaeaSQL   string
				CheckSQL  string
				ExpectRes string
			}{
				{
					GaeaSQL:   `set character_set_client=binary`,
					CheckSQL:  `show variables like "character_set_client"`,
					ExpectRes: "binary",
				},
				{
					GaeaSQL:   `set character_set_connection=utf8, character_set_results=utf8, character_set_client=binary`,
					CheckSQL:  `show variables like "character_set_connection"`,
					ExpectRes: "utf8",
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

		ginkgo.It("set some SQL statements with syntax errors", func() {
			sqlCases := []struct {
				GaeaSQL   string
				CheckSQL  string
				ExpectRes string
			}{
				{
					GaeaSQL:   `/*!40101 SET @@SQL_MODE := @OLD_SQL_MODE, @@SQL_QUOTE_SHOW_CREATE := @OLD_QUOTE */`,
					CheckSQL:  `select @@SQL_MODE`,
					ExpectRes: "",
				},
			}
			for _, sqlCase := range sqlCases {
				gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
				util.ExpectNoError(err)
				_, err = gaeaConn.Exec(sqlCase.GaeaSQL)
				util.ExpectNoError(err)
				// 不检测结果，只验证释放能够继续执行
				_, err = gaeaConn.Exec(sqlCase.CheckSQL)
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
