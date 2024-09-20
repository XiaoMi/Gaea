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
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

// Test SET Invalid Value to Session Variables
// This test suite verifies the system's handling of invalid session variable values.
// Specifically, it focuses on how the system reacts when an invalid sql_mode is set for a session.

// The core test checks the following:
// 1. Establish a connection to Gaea with read-write permissions.
// 2. Retrieve the default sql_mode from the session.
// 3. Attempt to set an invalid sql_mode value and expect the system to handle it without crashing or applying the invalid value.
// 4. Verify that the system does not apply the invalid value and continues to operate normally.

// The test ensures robust error handling by:
// - An error is expected on the first query after setting an invalid sql_mode.
// - Confirming that no lasting changes are made to the sql_mode after the error.
// - Ensure that subsequent SQL requests are processed normally

// Cleanup after each test:
// - Resetting modifications and cleaning up resources to ensure isolation between tests.
var _ = ginkgo.Describe("Test SET Invalid Value to Session Variables", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	table := config.DefaultE2eTable
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
	// 获取默认的 SQL mode
	getDefaultSQLMode := func(conn *sql.DB) (string, error) {
		var sqlMode string
		rows, err := conn.Query("SELECT @@SESSION.sql_mode;")
		if err != nil {
			return "", err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&sqlMode)
			if err != nil {
				return "", err
			}
		}
		return sqlMode, nil
	}
	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)

		// namespace prepare
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")

		// modify namespace
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data sync and namespace load
		time.Sleep(3 * time.Millisecond)
	})

	ginkgo.Context("test set invalid sql_mode value", func() {
		ginkgo.It("should handle invalid sql_mode setting and fall back to default value", func() {
			// Get connection
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)

			//Try to get the default sql_mode
			defaultMode, err := getDefaultSQLMode(gaeaConn)
			util.ExpectNoError(err, "Expected error when getting default sql_mode")

			// Attempt to set invalid sql_mode
			_, err = gaeaConn.Exec("SET SESSION sql_mode = 'aaa';")
			util.ExpectNoError(err, "Expected error when setting an invalid sql_mode")

			// An error occurred during the first query
			_, err = getDefaultSQLMode(gaeaConn)
			util.ExpectError(err)

			// No error occurred during the second query
			setMode, err := getDefaultSQLMode(gaeaConn)
			util.ExpectNoError(err)
			util.ExpectEqual(defaultMode, setMode)

		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
