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

// This script, titled "Show Sql Type" is designed to evaluate the behavior of SQL queries, particularly SHOW VARIABLES LIKE, in a dual-slave database configuration.
// The script sets up the environment and then executes the SHOW VARIABLES LIKE SQL command to determine whether it's correctly directed to the master node in the setup.
// This test ensures that specific queries are handled by the appropriate database nodes (master or slave) as expected in a read-write split environment.
// The results from the SQL log are checked to confirm that the execution happens on the correct backend address, validating the routing logic in the database setup.
var _ = ginkgo.Describe("Show Sql Type", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		_, err = masterAdminConn.Exec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, db))
		util.ExpectNoError(err)
		// namespace prepare
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(3 * time.Millisecond)
	})

	ginkgo.Context("show variables like", func() {
		ginkgo.It("show variables like to master or slave", func() {
			gaeaConns, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
			sqlCase := []struct {
				sql         string
				gaeaConns   []*sql.DB
				expectAddrs []string
			}{
				{
					sql:         `show variables like "read_only"`,
					gaeaConns:   []*sql.DB{gaeaConns},
					expectAddrs: []string{slice.Slices[0].Master},
				},
			}

			for _, sqlCase := range sqlCase {
				for i := 0; i < len(sqlCase.expectAddrs); i++ {
					sql := fmt.Sprintf("/*i:%d*/ %s", i, sqlCase.sql)
					_, err := sqlCase.gaeaConns[i].Exec(sql)
					util.ExpectNoError(err)
					res, err := e2eMgr.SearchSqlLog(sql, currentTime)
					util.ExpectNoError(err)
					gomega.Expect(res).Should(gomega.HaveLen(1))
					gomega.Expect(res[0].BackendAddr).Should(gomega.Equal(sqlCase.expectAddrs[i]))
				}
			}

		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
