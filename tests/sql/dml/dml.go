package dml

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sql_support_test", func() {

	var planManagers []*config.PlanManager
	var err error
	var allowedDBS map[string]bool
	var slice []*models.Slice
	var users []*models.User
	var sqlPlanNamespace *models.Namespace
	ginkgo.BeforeEach(func() {
		// 创建数据库实例
		ginkgo.By("create slice")
		slice = []*models.Slice{
			{
				Name:     "slice-0",
				UserName: "superroot",
				Password: "superroot",
				Master:   "127.0.0.1:3319",
				Slaves: []string{
					"127.0.0.1:3329",
					"127.0.0.1:3339",
				},
				StatisticSlaves: []string{},
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
				InitConnect:     "",
			},
		}
		allowedDBS = map[string]bool{
			"db_test":         true,
			"db_test_delete":  true,
			"db_test_insert":  true,
			"db_test_replace": true,
			"db_test_select":  true,
			"db_test_update":  true,
		}
		//  创建数据库用户
		ginkgo.By("create user")
		users = []*models.User{
			{
				UserName:      "superroot",
				Password:      "superroot",
				Namespace:     "test_namespace_sql_plan",
				RWFlag:        2,
				RWSplit:       1,
				OtherProperty: 0,
			},
		}
		ginkgo.By("create namespace")
		sqlPlanNamespace = config.NewNamespace(
			config.WithIsEncrypt(true),
			config.WithName("test_namespace_sql_plan"),
			config.WithAllowedDBS(allowedDBS),
			config.WithSlices(slice),
			config.WithUsers(users),
			config.WithDefaultSlice("slice-0"),
		)

		// 注册命令空间
		ginkgo.By("add namespace")
		err = config.GetNamespaceRegisterManager().AddNamespace(sqlPlanNamespace)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("register namespace")
		err = config.GetNamespaceRegisterManager().RegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Cluster conn ")
		mysqlCluster, err := config.InitMysqlClusterConn(slice)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Gaea conn ")
		gaeaDB, err := config.InitGaeaConn(users[0], "127.0.0.1:13306")
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("init case path")
		casesPath := []string{
			util.GetCasesFileAbsPath("dml/case/sql.json"),
			util.GetCasesFileAbsPath("dml/case/insert.json"),
			util.GetCasesFileAbsPath("dml/case/delete.json"),
			util.GetCasesFileAbsPath("dml/case/update.json"),
			util.GetCasesFileAbsPath("dml/case/replace.json"),
			util.GetCasesFileAbsPath("dml/case/select.json"),
		}
		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := config.NewPlanManager(
				config.WithPlanPath(v),
				config.WithMysqlClusterConn(mysqlCluster),
				config.WithGaeaConn(gaeaDB),
			)
			ginkgo.By("connect all db")
			planManagers = append(planManagers, p)
		}

	})
	ginkgo.Context("sql support test", func() {
		ginkgo.It("sql support", func() {
			for _, p := range planManagers {
				err = p.Init()
				gomega.Expect(err).Should(gomega.BeNil())
				err := p.Run()
				gomega.Expect(err).Should(gomega.BeNil())
			}
		})
	})
	ginkgo.AfterEach(func() {
		for _, p := range planManagers {
			p.MysqlClusterConnClose()
		}
		// 删除注册
		config.GetNamespaceRegisterManager().UnRegisterNamespaces()
	})
})
