// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package function

import (
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

// This script, titled "Show Sql Type" is designed to evaluate the behavior of SQL queries, particularly SHOW VARIABLES LIKE, in a dual-slave database configuration.
// The script sets up the environment and then executes the SHOW VARIABLES LIKE SQL command to determine whether it's correctly directed to the master node in the setup.
// This test ensures that specific queries are handled by the appropriate database nodes (master or slave) as expected in a read-write split environment.
// The results from the SQL log are checked to confirm that the execution happens on the correct backend address, validating the routing logic in the database setup.
var _ = ginkgo.Describe("SET SESSION transaction_read_only Test", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	table := config.DefaultE2eTable
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)

		// namespace prepare
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")

		// limit slice capacity = 1 and  max_capacity = 1
		initNs.Slices[0].Capacity = 1
		initNs.Slices[0].MaxCapacity = 1

		// modify namespace
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(3 * time.Millisecond)
	})

	ginkgo.Context("test set session readonly", func() {
		ginkgo.It("should correctly handle read-only transactions and write operations", func() {
			// 获取两个连接，模拟两个不同的客户端
			gaeaConn1, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
			gaeaConn2, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)

			// 客户端1开始事务，并设置为只读
			_, err = gaeaConn1.Exec("BEGIN;")
			util.ExpectNoError(err)
			_, err = gaeaConn1.Exec("SET SESSION transaction_read_only = 1;")
			util.ExpectNoError(err)
			_, err = gaeaConn1.Exec(fmt.Sprintf("SELECT * FROM %s.%s;", db, table))
			util.ExpectNoError(err)
			_, err = gaeaConn1.Exec("COMMIT;")
			util.ExpectNoError(err)

			// 客户端2尝试写入操作
			_, err = gaeaConn2.Exec("BEGIN;")
			util.ExpectNoError(err)
			_, err = gaeaConn2.Exec(fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", db, table), "TestSessionReadOnly")
			util.ExpectNoError(err)
			_, err = gaeaConn2.Exec("COMMIT;")
			util.ExpectNoError(err)

		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
