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
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
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
var collationTestSqlCases = []struct {
	TestSQL      string
	CheckSQL     string
	ExpectMySQL5 string
	ExpectMySQL8 string
}{
	{
		TestSQL:      `SET NAMES 'utf8mb4' COLLATE 'utf8mb4_0900_ai_ci'`,
		CheckSQL:     "show variables like 'collation_connection'",
		ExpectMySQL5: "utf8mb4_general_ci",
		ExpectMySQL8: "utf8mb4_0900_ai_ci",
	},
}

var _ = ginkgo.Describe("only mysql5: SET NAMES utf8mb4 COLLATE utf8mb4_0900_ai_ci", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
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

	ginkgo.It("should set collate 'utf8mb4_0900_ai_ci' get 'utf8mb4_general_ci'", func() {
		gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
		util.ExpectNoError(err)

		for _, sqlCase := range collationTestSqlCases {
			gaeaConn.Exec(sqlCase.TestSQL)
			res, err := util.MysqlQuery(gaeaConn, sqlCase.CheckSQL)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0]).Should(gomega.HaveLen(2))
			gomega.Expect(res[0][1]).Should(gomega.Equal(sqlCase.ExpectMySQL5))
		}
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})

var _ = ginkgo.Describe("only mysql8: SET NAMES utf8mb4 COLLATE utf8mb4_0900_ai_ci", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
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

	ginkgo.It("should set collate 'utf8mb4_0900_ai_ci' get 'utf8mb4_0900_ai_ci", func() {
		gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
		util.ExpectNoError(err)

		for _, sqlCase := range collationTestSqlCases {
			gaeaConn.Exec(sqlCase.TestSQL)
			res, err := util.MysqlQuery(gaeaConn, sqlCase.CheckSQL)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0]).Should(gomega.HaveLen(2))
			gomega.Expect(res[0][1]).Should(gomega.Equal(sqlCase.ExpectMySQL8))
		}
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
