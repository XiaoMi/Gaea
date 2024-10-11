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
)

var _ = ginkgo.Describe("default support information_schema db test", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
	table := config.DefaultE2eTable
	initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
	util.ExpectNoError(err, "parse namespace template")

	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(50 * time.Millisecond)
	})

	ginkgo.Context("support information_schema db by default", func() {
		ginkgo.It("should support information_schema db by default", func() {
			sqlCases := []struct {
				testSql  string
				expected [][]string
			}{
				{
					testSql:  `show databases`,
					expected: [][]string{{"db_e2e_test"}, {"information_schema"}},
				},
				{
					testSql:  `use information_schema`,
					expected: [][]string{},
				},
				{
					testSql:  `use information_SCHEMA`,
					expected: [][]string{},
				},
			}

			gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)
			for _, test := range sqlCases {
				err = checkFunc(gaeaConn, test.testSql, test.expected)
				util.ExpectNoError(err)
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
