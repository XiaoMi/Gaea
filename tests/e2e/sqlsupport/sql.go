package sqlsupport

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sql_test", func() {
	var manager *config.NamespaceRegisterManger
	var ns string
	var dbName string
	var tableName string
	var gaeaDB *sql.DB
	var masterDB *sql.DB
	var err error

	var gaeaUsers []*config.DBConfig
	var gaeaHost string
	var gaeaPort int
	ginkgo.BeforeEach(func() {
		ginkgo.By("init common config")
		manager = config.GetNamespaceRegisterManager()
		ns = config.DefaulNamespace.Name
		gaeaHost = "127.0.0.1"
		gaeaPort = 13306
		dbName = "sql_DB"
		tableName = "tbl_sql"
		ginkgo.By("init mysqlDB config")
		m := config.NewMysqlClusterManager(manager.GetNsSlices(ns))
		s := m.GetMysqlSliceByName("slice-0")

		ginkgo.By("connect mysql")
		masterDB, err = util.ConnectMysql(s.Master)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init gaeaDB config")
		gaeaUsers = config.NewGaeaDBConfig(manager.GetNsUsers(ns), fmt.Sprintf("%s:%d", gaeaHost, gaeaPort))
		gomega.Expect(gaeaUsers).Should(gomega.HaveLen(1))

		ginkgo.By("connect gaea")
		gaeaDB, err = util.ConnectMysql(gaeaUsers[0])
		gomega.Expect(err).Should(gomega.BeNil())

		// masterDB 初始数据库和表
		ginkgo.By("create database and insert 10 data")
		err = util.CreateDatabaseAndInsertData(masterDB, dbName, tableName, 10)
		gomega.Expect(err).Should(gomega.BeNil())

	})
	ginkgo.Context("hint", func() {
		ginkgo.It("sql hint", func() {
			// gaea  transaction test
			ginkgo.By("delete one record by gaea")
			gTx, err := gaeaDB.Begin()
			gomega.Expect(err).Should(gomega.BeNil())

			_, err = gTx.Exec(fmt.Sprintf(`delete from %s.%s where id=1;`, dbName, tableName))
			gomega.Expect(err).Should(gomega.BeNil())

			err = gTx.Commit()
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea select test
			ginkgo.By("query remaining data count")
			var count int
			err = gaeaDB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s.%s;`, dbName, tableName)).Scan(&count)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(count).Should(gomega.Equal(9))

		})
	})
	ginkgo.AfterEach(func() {
		ginkgo.By("begin clean database.")
		util.DropDatabase(masterDB, []string{
			dbName,
		})
		if gaeaDB != nil {
			gaeaDB.Close()
		}
		if masterDB != nil {
			masterDB.Close()
		}
	})
})
