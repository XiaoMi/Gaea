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
	"log"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("mysql bad connection test", func() {
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

	ginkgo.Context("test mysql bad connection", func() {
		ginkgo.It("should handle test mysql bad connection correctly", func() {
			gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)

			sqlCases := []struct {
				TestSQL   string
				ExpectErr string
			}{
				{
					TestSQL:   `select sleep(100)`,
					ExpectErr: "Error 1105: connection was bad",
				},
			}

			for _, sqlCase := range sqlCases {
				masterAdminConn, err := slice.GetMasterAdminConn(0)
				util.ExpectNoError(err)
				go killMySQLProcess(sqlCase.TestSQL, masterAdminConn)
				_, err = gaeaConn.Exec(sqlCase.TestSQL)
				time.Sleep(500 * time.Millisecond)
				util.ExpectEqual(err.Error(), sqlCase.ExpectErr)
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})

func killMySQLProcess(sql string, mysqlConn *sql.DB) {
	time.Sleep(time.Second)
	rows, err := mysqlConn.Query("select ID from information_schema.processlist where INFO = ?", sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}

		_, err := mysqlConn.Exec(fmt.Sprintf("KILL %s", id))
		if err != nil {
			util.ExpectNoError(err)
		}
	}
}