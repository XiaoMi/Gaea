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

// This script, titled "Read-Write Splitting" focuses on verifying the correct routing of read and write operations in a database system with a dual-slave configuration.
// The test sets up various scenarios using different types of database connections (read-only, write-only, and read-write) to execute SQL queries.
// The goal is to ensure that read operations are directed to the slave nodes and write operations to the master node.
// It uses a mixture of successful and expected-to-fail SQL operations to validate the read-write splitting functionality.
// The results are then verified to confirm that the queries are executed on the intended backend nodes.
var _ = ginkgo.Describe("Read-Write Splitting", func() {
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

	ginkgo.Context("When handling read and write operations", func() {
		ginkgo.It("should direct read operations to replicas", func() {
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaReadWriteConn.Close()

			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaWriteConn.Close()

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaReadConn.Close()

			sqlCases := []struct {
				GaeaConn          *sql.DB
				GaeaSQL           string
				ExpectBackendAddr []string
				IsSuccess         bool
			}{
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Slaves[0], slice.Slices[0].Slaves[1]},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("DELETE FROM %s.%s WHERE id=2", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				// 写用户
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("DELETE FROM %s.%s WHERE id=2", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Master},
					IsSuccess:         true,
				},
				//读用户
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Slaves[1], slice.Slices[0].Slaves[0]},
					IsSuccess:         true,
				},

				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: []string{slice.Slices[0].Slaves[1], slice.Slices[0].Slaves[0]},
					IsSuccess:         true,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("DELETE FROM %s.%s  WHERE `id`= 2", db, table),
					IsSuccess: false,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					IsSuccess: false,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					IsSuccess: false,
				},
			}
			for index, sqlCase := range sqlCases {
				e2eMgr.ClearSqlLog()
				currentTime := time.Now()
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				if !sqlCase.IsSuccess {
					util.ExpectError(err, fmt.Sprintf("sql Case %d", index))
					continue
				} else {
					util.ExpectNoError(err, fmt.Sprintf("sql Case %d", index))
				}

				var res []util.LogEntry
				retryCount := 3 // 设置重试次数
				for i := 0; i < retryCount; i++ {
					time.Sleep(500 * time.Millisecond) // 等待一段时间再重试
					res, err = e2eMgr.SearchSqlLog(sqlCase.GaeaSQL, currentTime)
					if err == nil && len(res) == 1 {
						break
					}
				}
				// 检查结果
				util.ExpectNoError(err)
				gomega.Expect(res).Should(gomega.HaveLen(1))
				if len(sqlCase.ExpectBackendAddr) == 1 {
					gomega.Expect(sqlCase.ExpectBackendAddr[0]).Should(gomega.Equal(res[0].BackendAddr))
				} else {
					gomega.Expect(sqlCase.ExpectBackendAddr).Should(gomega.ContainElement(res[0].BackendAddr))
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
