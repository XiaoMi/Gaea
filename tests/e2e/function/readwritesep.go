package function

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Read-Write Splitting", func() {
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

		table = "tbl_read_write_test"
		// 创建表
		_, err = masterDB.Exec(fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table))
		gomega.Expect(err).Should(gomega.BeNil())
		// 插入10条数据
		for i := 0; i < 10; i++ {
			_, err = masterDB.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue")
			gomega.Expect(err).Should(gomega.BeNil())
		}
	})
	ginkgo.Context("When handling read and write operations", func() {
		ginkgo.It("should direct read operations to replicas", func() {

			// read
			sql := fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table)
			_, err = gaeaDB.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err := util.ReadLog(util.GetTestLogFileAbsPath(), sql)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0].BackendAddr).Should(gomega.Or(
				gomega.Equal("127.0.0.1:3329"),
				gomega.Equal("127.0.0.1:3339"),
			))

			// write
			sql = fmt.Sprintf("INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "Alex")
			_, err = gaeaDB.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err = util.ReadLog(util.GetTestLogFileAbsPath(), sql)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0].BackendAddr).Should(gomega.Equal("127.0.0.1:3319"))
		})
	})
	ginkgo.AfterEach(func() {
		masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", db))
		config.MysqlClusterConnClose(mysqlCluster)
		// 删除注册
		config.GetNamespaceRegisterManager().UnRegisterNamespaces()
	})
})
