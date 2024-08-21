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
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

// This test suite is designed to validate the integration of backend authentication plugins within the Gaea database middleware.
// The test scenarios include:
// 1. Creating a MySQL user with 'mysql_native_password', granting them select permissions, and configuring Gaea to recognize and authenticate using these credentials.
// 2. Similarly, setting up another user with 'sha256_password', which represents a more secure authentication mechanism, and ensuring Gaea can handle this advanced encryption during the authentication process.
// Each test:
// - Configures user credentials specific to the authentication method being tested.
// - Executes SQL commands to create and grant permissions to these users.
// - Modifies the Gaea namespace to utilize these new credentials.
// - Performs database queries to confirm that the connections are appropriately authenticated and operational.
// Overall, this testing strategy is crucial for verifying that Gaea's integration with MySQL authentication plugins is robust
var _ = ginkgo.Describe("test backend auth plugin", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
	table := config.DefaultE2eTable
	masterAdminConn, err := slice.GetMasterAdminConn(0)
	util.ExpectNoError(err, "get master admin conn")
	initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
	util.ExpectNoError(err, "parse namespace template")
	ginkgo.BeforeEach(func() {
		// data prepare
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err, "setup database and insert data")
		// namespace prepare
		err = e2eMgr.DeleteNamespace(initNs)
		util.ExpectNoError(err, "delete namespace")
		initNs.Name = "test_auth_plugin"
		for i := 0; i < len(initNs.Users); i++ {
			initNs.Users[i].Namespace = initNs.Name
		}
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err, "create namespace")
		// wait mysql data  sync and namespace load
		time.Sleep(3 * time.Millisecond)

	})
	ginkgo.Context("backend auth plugin", func() {
		ginkgo.It("auth mysql_native_password", func() {
			mysqlUser, mysqlPass := "native_user", "native_pass"
			grantSql := fmt.Sprintf(`CREATE USER IF NOT EXISTS "%s"@'%%' IDENTIFIED WITH mysql_native_password BY "%s"`, mysqlUser, mysqlPass)
			_, err = masterAdminConn.Exec(grantSql)
			util.ExpectNoError(err, "create user")
			grantSql2 := fmt.Sprintf(`GRANT SELECT ON %s.* TO "%s"@'%%'`, db, mysqlUser)
			_, err = masterAdminConn.Exec(grantSql2)
			util.ExpectNoError(err, "grant user")
			ns := initNs
			ns.Slices[0].UserName = mysqlUser
			ns.Slices[0].Password = mysqlPass
			err = e2eMgr.ModifyNamespace(ns)
			util.ExpectNoError(err, "modify namespace")
			// wait mysql data  sync and namespace load
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err, "get gaea read write conn")
			rs, err := gaeaReadWriteConn.Query(fmt.Sprintf(`select count(*) from %s.%s`, db, table))
			util.ExpectNoError(err, "query")
			for rs.Next() {
				var count int
				err = rs.Scan(&count)
				util.ExpectNoError(err, "scan gaea result")
				util.ExpectEqual(count, 10, "count")
			}
		})

		ginkgo.It("auth sha256_password", func() {
			mysqlUser, mysqlPass := "sha256_user", "sha256_pass"
			backendSql := []string{
				fmt.Sprintf(`DROP USER IF EXISTS "%s"@'%%'`, mysqlUser),
				fmt.Sprintf(`CREATE USER IF NOT EXISTS "%s"@'%%' IDENTIFIED WITH sha256_password BY "%s"`, mysqlUser, mysqlPass),
				fmt.Sprintf(`GRANT SELECT ON %s.* TO "%s"@'%%'`, db, mysqlUser),
			}
			for _, sql := range backendSql {
				_, err = masterAdminConn.Exec(sql)
				util.ExpectNoError(err, "init sql", sql)
			}
			ns := initNs
			ns.Slices[0].UserName = mysqlUser
			ns.Slices[0].Password = mysqlPass
			err = e2eMgr.ModifyNamespace(ns)
			util.ExpectNoError(err, "modify namespace")
			// wait mysql data  sync and namespace load
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err, "get gaea read write conn")
			rs, err := gaeaReadWriteConn.Query(fmt.Sprintf(`select count(*) from %s.%s`, db, table))
			util.ExpectNoError(err, "query")
			for rs.Next() {
				var count int
				err = rs.Scan(&count)
				util.ExpectNoError(err, "scan gaea result")
				util.ExpectEqual(count, 10, "count")
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
