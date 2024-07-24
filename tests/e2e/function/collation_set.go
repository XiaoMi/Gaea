package function

import (
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"time"
)

// mysql5 和 mysql8 的情况是分开来运行的
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
