package function

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// This script, titled "Test slave fuse" is designed to evaluate the behavior of a database system when slave nodes encounter specific issues.
// The test involves modifying database privileges and configurations to simulate scenarios where slave nodes have restricted access or altered roles.
// It checks whether the system correctly handles these changes without causing any service interruptions or incorrect data routing.
// The test ensures the resilience and stability of the database system in handling slave node issues, maintaining consistent query distribution and data integrity.
var _ = ginkgo.Describe("Test slave fuse", func() {
	e2eMgr := config.NewE2eManager()
	db, table := e2eMgr.Db, e2eMgr.Table
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
	initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
	util.ExpectNoError(err, "parse Namespace Tmpl")
	masterAdminConn, err := slice.GetMasterAdminConn(0)
	util.ExpectNoError(err, "get master admin conn")

	ginkgo.BeforeEach(func() {
		// mysql prepare
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err, "setup database and insert data")
		// namespace prepare
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err, "create namespace")
		// wait mysql data  sync and namespace load
		time.Sleep(1 * time.Second)
	})

	ginkgo.It("slave will not fuse when no privilege to show slave status", func() {
		e2eMgr.ClearSqlLog()
		e2eMgr.StartTime = time.Now()
		// step1: revoke mysql cluster privilege
		_, err = masterAdminConn.Exec(fmt.Sprintf(`REVOKE REPLICATION SLAVE, REPLICATION CLIENT ON *.* from '%s'@'%%'`, slice.Slices[0].UserName))
		util.ExpectNoError(err, "revoke replication slave")

		// step2: change cluster master to Gaea config slave, cluster slave to Gaea config master for test.
		sql := fmt.Sprintf("SELECT /*check slave fuse*/ * FROM %s.%s WHERE `id`= 1", db, table)
		counts := 30
		ns := initNs
		// set slave to one slave for checking log
		ns.Slices[0].Slaves = []string{slice.Slices[0].Slaves[0]}
		ns.SecondsBehindMaster = 5
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err, "modify namespace")

		// step3: continue query and check the query distribution.
		gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
		util.ExpectNoError(err)
		for i := 0; i < counts; i++ {
			_, err := gaeaReadConn.Exec(sql)
			util.ExpectNoError(err)
			time.Sleep(300 * time.Millisecond)
		}

		// step4: check the gaea log for distribution.
		res, err := e2eMgr.SearchSqlLog(sql, e2eMgr.StartTime)
		util.ExpectNoError(err)
		gomega.Expect(res).Should(gomega.HaveLen(counts))
		for _, r := range res {
			gomega.Expect(ns.Slices[0].Slaves[0]).Should(gomega.Equal(r.BackendAddr))
		}

		// step5: reset mysql cluster privilege
		_, err = masterAdminConn.Exec(fmt.Sprintf(`GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* to '%s'@'%%'`, slice.Slices[0].UserName))
		util.ExpectNoError(err, "grant replication slave")
	})

	ginkgo.It("slave will not fuse when show slave status is empty", func() {
		e2eMgr.ClearSqlLog()
		e2eMgr.StartTime = time.Now()
		// step1: change cluster master to Gaea config slave, cluster slave to Gaea config master for test.
		ns := initNs
		ns.Slices[0].Master = slice.Slices[0].Slaves[0]
		ns.Slices[0].Slaves = []string{slice.Slices[0].Master}
		ns.SecondsBehindMaster = 10
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err, "modify namespace")

		// step3: continue query and check the query distribution.
		sql := fmt.Sprintf("SELECT /*check slave fuse*/ * FROM %s.%s WHERE `id`= 1", db, table)
		counts := 30
		gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
		util.ExpectNoError(err)
		for i := 0; i < counts; i++ {
			_, err := gaeaReadConn.Exec(sql)
			util.ExpectNoError(err)
			time.Sleep(300 * time.Millisecond)
		}

		// step4: check the gaea log for distribution.
		res, err := e2eMgr.SearchSqlLog(sql, e2eMgr.StartTime)
		util.ExpectNoError(err)
		gomega.Expect(res).Should(gomega.HaveLen(counts))
		for _, r := range res {
			gomega.Expect(ns.Slices[0].Slaves[0]).Should(gomega.Equal(r.BackendAddr))
		}
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
