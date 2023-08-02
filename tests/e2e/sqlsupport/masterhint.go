package sqlsupport

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("master_hint_test", func() {
	var manager *config.NamespaceRegisterManger
	var ns string
	var dbName string
	var tableName string
	var gaeaDB *sql.DB
	var masterDB *sql.DB
	var slave1DB *sql.DB
	var slave2DB *sql.DB
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
		dbName = "master_hint_DB"
		tableName = "tbl_master_hint"
		ginkgo.By("init mysqlDB config")
		m := config.NewMysqlClusterManager(manager.GetNsSlices(ns))
		s := m.GetMysqlSliceByName("slice-0")

		ginkgo.By("connect masterDB")
		masterDB, err = util.ConnectMysql(s.Master)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("connect slave1DB")
		slave1DB, err = util.ConnectMysql(s.Slaves[0])
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("connect slave2DB")
		slave2DB, err = util.ConnectMysql(s.Slaves[1])
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init gaeaDB config")
		gaeaUsers = config.NewGaeaDBConfig(manager.GetNsUsers(ns), fmt.Sprintf("%s:%d", gaeaHost, gaeaPort))
		gomega.Expect(gaeaUsers).Should(gomega.HaveLen(1))

		ginkgo.By("connect gaeaDB")
		gaeaDB, err = util.ConnectMysql(gaeaUsers[0])
		gomega.Expect(err).Should(gomega.BeNil())
		// 插入10条数据
		err = util.CreateDatabaseAndInsertData(masterDB, dbName, tableName, 10)
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = masterDB.Exec("TRUNCATE TABLE mysql.general_log;")
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = slave1DB.Exec("TRUNCATE TABLE mysql.general_log;")
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = slave2DB.Exec("TRUNCATE TABLE mysql.general_log;")
		gomega.Expect(err).Should(gomega.BeNil())
	})
	ginkgo.Context("hint", func() {
		ginkgo.It("master hint", func() {
			var (
				eventTime   string
				userHost    string
				threadID    int
				serverID    int
				commandType string
				argument    string
			)
			_, err = gaeaDB.Exec(fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 9", dbName, tableName))

			rows, err := masterDB.Query("select * from mysql.general_log order by event_time desc limit 1 offset 1")
			gomega.Expect(err).Should(gomega.BeNil())

			if rows.Next() {
				err = rows.Scan(&eventTime, &userHost, &threadID, &serverID, &commandType, &argument)
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(strings.Contains(argument, "9")).To(gomega.BeTrue())
			}
			rows, err = slave1DB.Query("select * from mysql.general_log order by event_time desc limit 1 offset 1")
			gomega.Expect(err).Should(gomega.BeNil())
			if rows.Next() {
				err = rows.Scan(&eventTime, &userHost, &threadID, &serverID, &commandType, &argument)
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(strings.Contains(argument, "9")).To(gomega.BeFalse())
			}

			rows, err = slave2DB.Query("select * from mysql.general_log order by event_time desc limit 1 offset 1")
			gomega.Expect(err).Should(gomega.BeNil())
			if rows.Next() {
				err = rows.Scan(&eventTime, &userHost, &threadID, &serverID, &commandType, &argument)
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(strings.Contains(argument, "9")).To(gomega.BeFalse())

			}
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
		if slave1DB != nil {
			slave1DB.Close()
		}
		if slave2DB != nil {
			slave2DB.Close()
		}

	})
})
