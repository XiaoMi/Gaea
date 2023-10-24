package function

import (
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Read-Write Splitting", func() {
	var e2eConfig *config.E2eConfig
	var masterSlaveCluster *config.MysqlClusterConfig
	var db string
	var table string
	var startTime string

	ginkgo.BeforeEach(func() {
		currentTime := time.Now()
		startTime = currentTime.Format("2006-01-02 15:04:05.999")

		e2eConfig = config.GetDefaultE2eConfig()
		masterSlaveCluster = e2eConfig.MasterSlaveCluster

		db = "db_e2e_test"
		// 这里的数据表名最好和其他的表名不同，为了防止日志中拉取出的两条日志相同
		table = "tbl_read_write_test"
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

		statements := []string{
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", db),
			fmt.Sprintf("CREATE DATABASE %s", db),
			fmt.Sprintf("USE %s", db),
			fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table),
		}

		for _, stmt := range statements {
			_, err = mysqlConn.Exec(stmt)
			gomega.Expect(err).Should(gomega.BeNil())
		}
		// 插入10条数据
		for i := 0; i < 10; i++ {
			_, err = mysqlConn.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue")
			gomega.Expect(err).Should(gomega.BeNil())
		}
	})
	ginkgo.Context("When handling read and write operations", func() {
		ginkgo.It("should direct read operations to replicas", func() {
			// 初始化gaea连接

			gaeaConn, err := config.InitConn(e2eConfig.GaeaUser.UserName, e2eConfig.GaeaUser.Password, e2eConfig.GaeaHost, "")
			gomega.Expect(err).Should(gomega.BeNil())
			defer gaeaConn.Close()

			// 初始化数据库连接
			mysqlConn, err := config.InitConn(masterSlaveCluster.Slices[0].UserName, masterSlaveCluster.Slices[0].Password, masterSlaveCluster.Slices[0].Master, "")
			gomega.Expect(err).Should(gomega.BeNil())
			defer mysqlConn.Close()
			// read
			sql := fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table)
			_, err = gaeaConn.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err := util.ReadLog(e2eConfig.FilepathJoin("cmd/logs/gaea_sql.log"), sql, startTime)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0].BackendAddr).Should(gomega.Or(
				gomega.Equal(masterSlaveCluster.Slices[0].Slaves[0]),
				gomega.Equal(masterSlaveCluster.Slices[0].Slaves[1]),
			))

			// write
			sql = fmt.Sprintf("INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "Alex")
			_, err = gaeaConn.Exec(sql)
			gomega.Expect(err).Should(gomega.BeNil())
			res, err = util.ReadLog(e2eConfig.FilepathJoin("cmd/logs/gaea_sql.log"), sql, startTime)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(res).Should(gomega.HaveLen(1))
			gomega.Expect(res[0].BackendAddr).Should(gomega.Equal(masterSlaveCluster.Slices[0].Master))

		})
	})
	ginkgo.AfterEach(func() {
		e2eConfig.NamespaceManager.UnRegisterNamespaces()
	})
})
