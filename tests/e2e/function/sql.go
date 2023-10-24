package function

import (
	"fmt"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Simple SQL Queries", func() {
	var e2eConfig *config.E2eConfig
	var masterSlaveCluster *config.MysqlClusterConfig
	var db string
	var table string
	ginkgo.BeforeEach(func() {
		// 注册命令空间
		e2eConfig = config.GetDefaultE2eConfig()
		masterSlaveCluster = e2eConfig.MasterSlaveCluster
		db = "db_e2e_test"
		// 这里的数据表名最好和其他的表名不同，为了防止日志中拉取出的两条日志相同
		table = "tbl_sql_test"
		// 解析模版
		ns, err := masterSlaveCluster.TemplateParse(e2eConfig.FilepathJoin("e2e/function/ns/default.template"))
		gomega.Expect(err).Should(gomega.BeNil())

		// 注册namespace
		err = e2eConfig.RegisterNamespaces(ns)
		gomega.Expect(err).Should(gomega.BeNil())

		// 初始化数据库实例
		mysqlConn, err := config.InitConn(masterSlaveCluster.Slices[0].UserName, masterSlaveCluster.Slices[0].Password, masterSlaveCluster.Slices[0].Master, "")
		gomega.Expect(err).Should(gomega.BeNil())
		defer mysqlConn.Close()

		commands := []string{
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", db),
			fmt.Sprintf("CREATE DATABASE %s", db),
			fmt.Sprintf("USE %s", db),
			fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table),
		}

		for _, cmd := range commands {
			_, err = mysqlConn.Exec(cmd)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		// 插入10条数据
		for i := 0; i < 10; i++ {
			_, err = mysqlConn.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue")
			gomega.Expect(err).Should(gomega.BeNil())
		}

	})
	ginkgo.Context("When executing basic SQL operations", func() {
		ginkgo.It("should handle SELECT operations correctly", func() {
			// 初始化gaea连接
			gaeaConn, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
			gomega.Expect(err).Should(gomega.BeNil())
			defer gaeaConn.Close()

			sql := fmt.Sprintf("/*master*/ SELECT * FROM `%s`.`%s` WHERE `id`= 9", db, table)
			_, err = gaeaConn.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			// 初始化数据库连接
			mysqlConn, err := config.InitConn(masterSlaveCluster.Slices[0].UserName, masterSlaveCluster.Slices[0].Password, masterSlaveCluster.Slices[0].Master, "")
			gomega.Expect(err).Should(gomega.BeNil())
			defer mysqlConn.Close()

			// gaea  delete  test
			ginkgo.By("delete one record by gaea")
			err = util.DeleteRow(gaeaConn, db, table, 1)
			gomega.Expect(err).Should(gomega.BeNil())

			err = util.DeleteVerify(mysqlConn, db, table, 1)
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea updata test
			ginkgo.By("update one record by gaea")
			err = util.UpdateRow(gaeaConn, db, table, 2, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())
			err = util.UpdateVerify(mysqlConn, db, table, 2, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea insert test
			ginkgo.By("insert one record by gaea")
			id, err := util.InsertRow(gaeaConn, db, table, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())
			err = util.InsertVerify(mysqlConn, db, table, id, "Tom")
			gomega.Expect(err).Should(gomega.BeNil())

			// gaea transaction test
			ginkgo.By("test gaea transaction")
			id, err = util.ExecuteTransactionAndReturnId(gaeaConn, fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", db, table), "Alex")
			gomega.Expect(err).Should(gomega.BeNil())
			util.InsertVerify(mysqlConn, db, table, id, "Alex")

		})
	})
	ginkgo.AfterEach(func() {
		e2eConfig.NamespaceManager.UnRegisterNamespaces()
	})
})
