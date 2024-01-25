package dml

import (
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("test dml set variables", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		initNs.Name = "test_dml_variables"
		for i := 0; i < len(initNs.Users); i++ {
			initNs.Users[i].Namespace = initNs.Name
		}

		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err, "create namespace")
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err, "get master admin conn")
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err, "setup database and insert data")
	})
	ginkgo.Context("test set session variables take effect", func() {
		gaeaRWConn, err := e2eMgr.GetReadWriteGaeaUserConn()
		gaeaRWConn.SetMaxOpenConns(1)
		gaeaRWConn.SetMaxIdleConns(1)
		util.ExpectNoError(err, "get gaea read write conn")
		ginkgo.It("set sql_mode", func() {
			testExecCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "STRICT_TRANS_TABLES",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "abcdefghijklmnopqrstuvwxyz"),
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "abcdefghijk"),
					},
					expectErr: []bool{true, false},
				},
			}

			for _, tt := range testExecCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Exec(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}

			testQueryCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "ANSI_QUOTES",
					sql: []string{
						`show variables like "sql_mode"`,
					},
					expectErr: []bool{true},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

		ginkgo.It("set session sql_select_limit", func() {
			testQueryCase := []struct {
				setSQL    string
				sql       []string
				expectErr []bool
				resultLen []int
			}{
				{
					setSQL: "set SQL_SELECT_LIMIT=default",
					sql: []string{
						fmt.Sprintf("select * from %s.%s", db, table),
						fmt.Sprintf("select * from %s.%s limit 1", db, table),
					},
					expectErr: []bool{false, false},
					resultLen: []int{10, 1},
				},
				{
					setSQL: "set SQL_SELECT_LIMIT=2",
					sql: []string{
						fmt.Sprintf("select * from %s.%s", db, table),
						fmt.Sprintf("select * from %s.%s limit 1", db, table),
					},
					expectErr: []bool{false, false},
					resultLen: []int{2, 1},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(tt.setSQL)
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.setSQL))
				for i := 0; i < len(tt.sql); i++ {
					res, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
					count := 0
					for res.Next() {
						count++
					}
					util.ExpectEqual(count, tt.resultLen[i], fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

		ginkgo.It("set session read write", func() {
			testExecCase := []struct {
				setSQL    string
				sql       []string
				expectErr []bool
			}{
				{
					setSQL: "set session transaction read only",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "a"),
						fmt.Sprintf("select * from %s.%s", db, table),
					},
					expectErr: []bool{true, false},
				},
				{
					setSQL: "set session transaction read write",
					sql: []string{
						fmt.Sprintf("insert into %s.%s(name) values('%s')", db, table, "b"),
						fmt.Sprintf("select * from %s.%s", db, table),
					},
					expectErr: []bool{false, false},
				},
			}

			for _, tt := range testExecCase {
				_, err := gaeaRWConn.Exec(tt.setSQL)
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.setSQL))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Exec(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s.", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}

			testQueryCase := []struct {
				sqlMode   string
				sql       []string
				expectErr []bool
			}{
				{
					sqlMode: "ANSI_QUOTES",
					sql: []string{
						`show variables like "sql_mode"`,
					},
					expectErr: []bool{true},
				},
			}

			for _, tt := range testQueryCase {
				_, err := gaeaRWConn.Exec(fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				util.ExpectNoError(err, fmt.Sprintf("set sql_mode='%s'", tt.sqlMode))
				for i := 0; i < len(tt.sql); i++ {
					_, err := gaeaRWConn.Query(tt.sql[i])
					if tt.expectErr[i] {
						util.ExpectError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
						continue
					}
					util.ExpectNoError(err, fmt.Sprintf("exec sql: %s", tt.sql[i]))
				}
			}
		})

	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
