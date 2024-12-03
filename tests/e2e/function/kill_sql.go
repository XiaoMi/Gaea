package function

import (
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"time"
)

var _ = ginkgo.Describe("kill sql", func() {
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

	ginkgo.It("should kill sql success", func() {
		gaeaConn, err := e2eMgr.GetWriteGaeaUserConn()
		util.ExpectNoError(err)
		var sqlCases = []struct {
			TestSQL string
		}{
			{
				TestSQL: `kill -1`,
			},
		}
		for _, sqlCase := range sqlCases {
			gaeaConn.Exec(sqlCase.TestSQL)
			_, err := util.MysqlExec(gaeaConn, sqlCase.TestSQL)
			util.ExpectNoError(err)
		}
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
