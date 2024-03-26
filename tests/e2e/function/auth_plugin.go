package function

import (
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

// test three cases:
// 1. mysql_native_password
// 2. sha256_password
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
