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
	"context"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// TestAdminDatabasePassthrough verifies admin user's database operation passthrough capabilities.
// The test suite includes:
// - Database creation and validation
// - Table creation and schema verification
// - Context switching between databases
// - Basic CRUD operation validation
var _ = ginkgo.Describe("Admin User Database Passthrough", func() {
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceAdminDualSlave]

	ginkgo.BeforeEach(func() {
		// Initialize MySQL connections
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.CleanUpDatabases(masterAdminConn)
		util.ExpectNoError(err)
		// Configure namespace template
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// Wait for configuration propagation
		time.Sleep(300 * time.Millisecond)
	})

	ginkgo.Context("When performing database operations", func() {
		ginkgo.It("should handle full database lifecycle operations", func() {
			ctx := context.Background()
			gaeaAdminConn, err := e2eMgr.GetAdminGaeaUserConn()
			util.ExpectNoError(err)

			gaeaCommonConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)

			// Database creation test cases
			sqlCases := []struct {
				GaeaSQL string
			}{
				// Database creation operations
				{GaeaSQL: "Create database metrics_dev"},
				{GaeaSQL: "Create database metrics_prod"},
				{GaeaSQL: "Create database metrics_test"},

				// Database context switching and table creation
				{GaeaSQL: fmt.Sprintf("USE `%s`", "metrics_dev")},
				{GaeaSQL: "CREATE TABLE `t_dev` (`row` INT)"},
				{GaeaSQL: fmt.Sprintf("USE `%s`", "metrics_prod")},
				{GaeaSQL: "CREATE TABLE `t_prod` (`row` INT)"},
				{GaeaSQL: fmt.Sprintf("USE `%s`", "metrics_test")},
				{GaeaSQL: "CREATE TABLE `t_test` (`row` INT)"},
			}
			// 执行 SQL 测试用例
			for _, sqlCase := range sqlCases {
				_, err := gaeaAdminConn.ExecContext(ctx, sqlCase.GaeaSQL)
				util.ExpectNoError(err)
			}

			// 1.验证admin show db是否可以透传
			rows, err := gaeaAdminConn.Query("show databases")
			util.ExpectNoError(err)
			defer rows.Close()
			res, err := util.GetDataFromRows(rows)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.Equal([][]string{{"information_schema"}, {"metrics_dev"}, {"metrics_prod"}, {"metrics_test"}, {"mysql"}, {"performance_schema"}, {"sys"}}))

			// 2.验证非admin show db是否可以透传
			rows, err = gaeaCommonConn.Query("show databases")
			util.ExpectNoError(err)
			defer rows.Close()
			res, err = util.GetDataFromRows(rows)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.ConsistOf([][]string{{"db_e2e_test"}, {"information_schema"}}))

			// 2.验证metrics_dev创建表t_dev是否创建成功
			db_dev, err := e2eMgr.GetAdminGaeaUserDBConn("metrics_dev")
			util.ExpectNoError(err)
			rows, err = db_dev.QueryContext(ctx, "show tables")
			util.ExpectNoError(err)
			defer rows.Close()
			res, err = util.GetDataFromRows(rows)
			util.ExpectNoError(err)

			gomega.Expect(res).Should(gomega.Equal([][]string{{"t_dev"}}))

			// 3.验证metrics_prod创建表t_prod是否创建成功
			db_prod, err := e2eMgr.GetAdminGaeaUserDBConn("metrics_prod")
			util.ExpectNoError(err)
			rows, err = db_prod.QueryContext(ctx, "show tables")
			util.ExpectNoError(err)
			defer rows.Close()
			res, err = util.GetDataFromRows(rows)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.Equal([][]string{{"t_prod"}}))

			// 4.验证metrics_test创建表t_test是否创建成功
			db_test, err := e2eMgr.GetAdminGaeaUserDBConn("metrics_test")
			util.ExpectNoError(err)
			rows, err = db_test.QueryContext(ctx, "show tables")
			util.ExpectNoError(err)
			defer rows.Close()
			res, err = util.GetDataFromRows(rows)
			util.ExpectNoError(err)
			gomega.Expect(res).Should(gomega.Equal([][]string{{"t_test"}}))

		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
