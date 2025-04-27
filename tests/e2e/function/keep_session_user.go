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

// Test type constants for KeepSession scenario
const (
	ReadUser      int = 0 // operation by read-only user in KeepSession mode
	WriteUser     int = 1 // operation by write-only user in KeepSession mode
	ReadWriteUser int = 3 // operation by read-write user in KeepSession mode
)

// KeepSessionRoutingTestSuite tests the routing behavior in KeepSession mode.
// It verifies that queries are properly routed based on user permissions and session persistence.
var _ = ginkgo.Describe("KeepSession Routing Test Suite", func() {
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
		initNs.SetForKeepSession = true
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(3 * time.Millisecond)
	})

	ginkgo.Context("When executing queries with KeepSession enabled", func() {
		ginkgo.It("should correctly route queries based on user permissions and session persistence", func() {
			// 定义 SQL 测试用例
			sqlCases := []struct {
				Name              string // 测试用例描述
				User              int    // 测试类型
				GaeaSQL           string
				ExpectBackendAddr []string
			}{
				// Read-only user test cases
				{
					Name:              "Read-only user - Read operation (SELECT)",
					User:              ReadUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case1*/ SELECT * FROM `%s`.`%s` WHERE id=1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Slaves[0], slice.Slices[0].Slaves[1]},
				},
				// Write-only user test cases (CRUD)
				{
					Name:              "Write-only user - Create operation (INSERT)",
					User:              WriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case2*/ INSERT INTO `%s`.`%s` (name) VALUES ('testValue')", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Write-only user - Read operation (SELECT)",
					User:              WriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case3*/ SELECT * FROM `%s`.`%s` WHERE id=1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Write-only user - Update operation (UPDATE)",
					User:              WriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case4*/ UPDATE `%s`.`%s` SET name='updatedValue' WHERE id=1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Write-only user - Delete operation (DELETE)",
					User:              WriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case5*/ DELETE FROM `%s`.`%s` WHERE id=1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},

				// Read-write user test cases (CRUD)
				{
					Name:              "Read-write user - Create operation (INSERT)",
					User:              ReadWriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case6*/ INSERT INTO `%s`.`%s` (name) VALUES ('testValue')", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Read-write user - Read operation (SELECT)",
					User:              ReadWriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case7*/ SELECT * FROM `%s`.`%s` WHERE id=2", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Read-write user - Update operation (UPDATE)",
					User:              ReadWriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case8*/ UPDATE `%s`.`%s` SET name='updatedValue' WHERE id=2", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
				{
					Name:              "Read-write user - Delete operation (DELETE)",
					User:              ReadWriteUser,
					GaeaSQL:           fmt.Sprintf("/*KeepSession Case9*/ DELETE FROM `%s`.`%s` WHERE id=2", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
				},
			}

			e2eMgr.ClearSqlLog()
			currentTime := time.Now()
			for _, sqlCase := range sqlCases {
				// 使用 Ginkgo 的 By 来标记每个测试步骤
				ginkgo.By(fmt.Sprintf("Running test case: %s", sqlCase.Name), func() {
					var gaeaConn *sql.DB
					var err error
					switch sqlCase.User {
					case ReadUser:
						gaeaConn, err = e2eMgr.GetReadGaeaUserDBConn(db)
						util.ExpectNoError(err, "Failed to get read-only user connection")
					case WriteUser:
						gaeaConn, err = e2eMgr.GetWriteGaeaUserDBConn(db)
						util.ExpectNoError(err, "Failed to get write-only user connection")
					case ReadWriteUser:
						gaeaConn, err = e2eMgr.GetReadWriteGaeaUserDBConn(db)
						util.ExpectNoError(err, "Failed to get read-write user connection")
					}
					// 执行 SQL
					_, err = gaeaConn.Exec(sqlCase.GaeaSQL)
					util.ExpectNoError(err, fmt.Sprintf("Test:%s, SQL Exec Faid: %v, SQL: %s", sqlCase.Name, err, sqlCase.GaeaSQL))

					// 查询日志
					res, err := e2eMgr.SearchSqlLog(sqlCase.GaeaSQL, currentTime)
					util.ExpectNoError(err, fmt.Sprintf("Test:%s, Search SQL Faid: %v, SQL: %s", sqlCase.Name, err, sqlCase.GaeaSQL))

					// 检查日志数量
					gomega.Expect(res).Should(gomega.HaveLen(1), "Test:%s, Log length is not 1, SQL: %s", sqlCase.Name, sqlCase.GaeaSQL)

					// 检查后端地址
					if len(sqlCase.ExpectBackendAddr) == 1 {
						gomega.Expect(sqlCase.ExpectBackendAddr[0]).Should(
							gomega.Equal(res[0].BackendAddr),
							"Backend address does not match\nTest: %s\nSQL: %s\nExpect: %s\nActual: %s",
							sqlCase.Name,
							sqlCase.GaeaSQL,
							sqlCase.ExpectBackendAddr[0],
							res[0].BackendAddr,
						)
					} else {
						gomega.Expect(sqlCase.ExpectBackendAddr).Should(
							gomega.ContainElement(res[0].BackendAddr),
							"Backend address is not in the expected range\nTest: %s\nSQL: %s\nExpect: %v\nActual: %s",
							sqlCase.Name,
							sqlCase.GaeaSQL,
							sqlCase.ExpectBackendAddr,
							res[0].BackendAddr,
						)
					}
				})
			}

		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
