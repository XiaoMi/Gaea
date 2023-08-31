package function

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Simple SQL Queries", func() {
	var err error
	var mysqlCluster map[string]*config.SliceConn
	var gaeaDB *sql.DB
	var masterDB *sql.DB
	var table string
	ginkgo.BeforeEach(func() {
		// 注册命令空间
		ginkgo.By("add namespace")
		err = config.GetNamespaceRegisterManager().AddNamespace(defaultNamespace)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("register namespace")
		err = config.GetNamespaceRegisterManager().RegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Cluster conn ")
		mysqlCluster, err = config.InitMysqlClusterConn(slice)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Gaea conn ")
		gaeaDB, err = config.InitGaeaConn(users[0], "127.0.0.1:13306")
		gomega.Expect(err).Should(gomega.BeNil())

		if slice, ok := mysqlCluster["slice-0"]; ok {
			gomega.Expect(len(slice.Slaves)).Should(gomega.BeNumerically(">=", 2), "Expected at least 2 slaves in slice")
			masterDB = slice.Master
		} else {
			gomega.Expect(ok).Should(gomega.BeTrue(), "Expected slice-0 to be present in mysqlCluster")
		}
		_, err = masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", db))
		gomega.Expect(err).Should(gomega.BeNil())

		_, err = masterDB.Exec(fmt.Sprintf("CREATE DATABASE  %s", db))
		gomega.Expect(err).Should(gomega.BeNil())

		_, err = masterDB.Exec(fmt.Sprintf("USE  %s", db))
		gomega.Expect(err).Should(gomega.BeNil())

		table = "tbl_sql_test"
		// 创建表
		tablesql := fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table)
		_, err = masterDB.Exec(tablesql)
		gomega.Expect(err).Should(gomega.BeNil())
		// 插入10条数据
		for i := 0; i < 10; i++ {
			_, err = masterDB.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue")
			gomega.Expect(err).Should(gomega.BeNil())
		}

	})
	ginkgo.Context("When executing basic SQL operations", func() {
		ginkgo.It("should handle SELECT operations correctly", func() {

			// gaea  delete  test
			ginkgo.By("delete one record by gaea")
			err = util.DeleteRow(gaeaDB, db, table, 1)
			gomega.Expect(err).Should(gomega.BeNil())
			err = util.DeleteVerify(masterDB, db, table, 1)
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea updata test
			ginkgo.By("update one record by gaea")
			err = util.UpdateRow(gaeaDB, db, table, 2, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())
			err = util.UpdateVerify(masterDB, db, table, 2, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea insert test
			ginkgo.By("insert one record by gaea")
			id, err := util.InsertRow(gaeaDB, db, table, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())
			err = util.InsertVerify(masterDB, db, table, id, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea transaction test
			ginkgo.By("test gaea transaction")
			id, err = util.ExecuteTransactionAndReturnId(gaeaDB, fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", db, table), "Alex")
			gomega.Expect(err).Should(gomega.BeNil())
			util.InsertVerify(masterDB, db, table, id, "Alex")
		})
	})
	ginkgo.AfterEach(func() {
		masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", db))
		config.MysqlClusterConnClose(mysqlCluster)
		// 删除注册
		config.GetNamespaceRegisterManager().UnRegisterNamespaces()
	})
})
