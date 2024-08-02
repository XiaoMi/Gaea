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
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/XiaoMi/Gaea/proxy/server"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

// This script, titled "prepare stmt test" is designed to verify prepare and exec in a database system.
var _ = ginkgo.Describe("prepare stmt test", func() {
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
		time.Sleep(500 * time.Millisecond)
	})

	ginkgo.Context("test prepare and exec stmts", func() {
		ginkgo.It("should handle prepare sqls correctly", func() {
			gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)

			sqlCases := []struct {
				PrepareSQL string
				ExecSQL    []interface{}
				ExpectRes  [][]string
				ExpectErr  string
			}{
				{
					PrepareSQL: fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ? and name = ?", db, table),
					ExecSQL:    []interface{}{9, "nameValue"},
					ExpectRes: [][]string{{
						"9", "nameValue",
					}},
				},
				{
					PrepareSQL: fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ? and name = ? limit 1", db, table),
					ExecSQL:    []interface{}{9, "nameValue"},
					ExpectRes: [][]string{{
						"9", "nameValue",
					}},
				},
				{
					PrepareSQL: fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ? and name = ?", db, table),
					ExecSQL:    []interface{}{9},
					ExpectErr:  "exec params mismatch",
				},
			}

			for _, sqlCase := range sqlCases {
				stmt, err := gaeaConn.Prepare(sqlCase.PrepareSQL)
				if sqlCase.ExpectErr != "" {
					count, _, _, _ := server.CalcParams(sqlCase.PrepareSQL)
					if len(sqlCase.ExecSQL) != count {
						err = errors.New("exec params mismatch")
					}

					util.ExpectEqual(err.Error(), sqlCase.ExpectErr)
					continue
				}
				util.ExpectNoError(err)

				err = checkPrepareExecRes(stmt, sqlCase.ExecSQL, sqlCase.ExpectRes)
				util.ExpectNoError(err)
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})

func checkPrepareExecRes(db *sql.Stmt, params []interface{}, values [][]string) error {
	rows, err := db.Query(params[0].(int), params[1].(string))
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows && len(values) == 0 {
			return nil
		}
		return fmt.Errorf("db Exec Error %v", err)
	}

	res, err := util.GetDataFromRows(rows)
	if err != nil {
		return fmt.Errorf("get data from rows error:%v", err)
	}
	if (len(res) == 0 || res == nil) && len(values) == 0 {
		return nil
	}
	if !reflect.DeepEqual(values, res) {
		return fmt.Errorf("mismatch. Actual: %v, Expect: %v", res, values)
	}

	return nil
}
