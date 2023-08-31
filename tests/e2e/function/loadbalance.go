package function

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Load Balancing", func() {
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

		// 这里的数据表名最好和其他的表名不同，为了防止日志中拉取出的两条日志相同
		table = "tbl_balance_test"
		// 创建表
		_, err = masterDB.Exec(fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table))
		gomega.Expect(err).Should(gomega.BeNil())

		// 插入10条数据
		for i := 0; i < 10; i++ {
			_, err = masterDB.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue")
			gomega.Expect(err).Should(gomega.BeNil())
		}

	})
	ginkgo.Context("When distributing queries among replicas", func() {
		ginkgo.It("should evenly distribute read queries", func() {
			result := make([]string, 2)
			sql := fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table)
			_, err = gaeaDB.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err := util.ReadLog(util.GetTestLogFileAbsPath(), sql)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			result[0] = res[0].BackendAddr

			sql = fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table)
			_, err = gaeaDB.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err = util.ReadLog(util.GetTestLogFileAbsPath(), sql)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			result[1] = res[0].BackendAddr
			// check
			actualMap := map[string]int{}
			for _, v := range result {
				actualMap[v]++
			}
			exceptMap := map[string]int{
				"127.0.0.1:3329": 1,
				"127.0.0.1:3339": 1,
			}
			gomega.Expect(actualMap).To(gomega.Equal(exceptMap))
		})
	})
	ginkgo.AfterEach(func() {
		ginkgo.By("begin clean database.")
		masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", db))
		config.MysqlClusterConnClose(mysqlCluster)

		// 删除注册
		config.GetNamespaceRegisterManager().UnRegisterNamespaces()

	})
})