package noshard

import (
	"database/sql"
	"fmt"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("unshard_dml_support_test", func() {

	var allowedDBS map[string]bool
	var users []*models.User
	var noShardPlanNamespace *models.Namespace
	var gaeaDB *sql.DB
	var gaeaMaster *sql.DB
	var casesPath []string
	var err error
	ginkgo.BeforeEach(func() {
		// 创建数据库实例
		ginkgo.By("create slice")
		gaeaSlice := []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        "superroot",
				Password:        "superroot",
				Master:          "127.0.0.1:3359",
				Slaves:          []string{},
				StatisticSlaves: []string{},
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
				InitConnect:     "",
			},
		}
		allowedDBS = map[string]bool{
			"sbtest1":       true,
			"sbtest1_shard": true,
		}
		//  创建数据库用户
		ginkgo.By("create user")
		users = []*models.User{
			{
				UserName:      "superroot",
				Password:      "superroot",
				Namespace:     "test_namespace_1",
				RWFlag:        2,
				RWSplit:       1,
				OtherProperty: 0,
			},
		}
		ginkgo.By("create namespace")
		noShardPlanNamespace = config.NewNamespace(
			config.WithIsEncrypt(true),
			config.WithName("test_namespace_1"),
			config.WithAllowedDBS(allowedDBS),
			config.WithSlices(gaeaSlice),
			config.WithUsers(users),
			config.WithDefaultSlice("slice-0"),
		)
		// 注册命令空间
		ginkgo.By("add namespace")
		err = config.GetNamespaceRegisterManager().AddNamespace(noShardPlanNamespace)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("register namespace")
		err = config.GetNamespaceRegisterManager().RegisterNamespaces()
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get Cluster conn ")
		gaeaCluster, err := config.InitMysqlClusterConn(gaeaSlice)
		gomega.Expect(err).Should(gomega.BeNil())
		if gaeaSlice, ok := gaeaCluster["slice-0"]; ok {
			gaeaMaster = gaeaSlice.Master
		} else {
			gomega.Expect(fmt.Errorf("slice not exist")).Should(gomega.BeNil())
		}
		ginkgo.By("get Gaea conn ")
		gaeaDB, err = config.InitGaeaConn(users[0], "127.0.0.1:13306")
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("get setup file")
		setupFile := util.GetCasesFileAbsPath("unshard/case/0-prepare.sql")
		sqls, err := util.GetSqlFromFile(setupFile)
		gomega.Expect(err).Should(gomega.BeNil())

		for _, v := range sqls {
			sqlOperation, err := util.GetSqlType(v)
			if err != nil {
				continue
			}
			if sqlOperation == util.UnKnown {
				continue
			}
			_, err = util.ExecuteSQL(gaeaMaster, "Exec", v)
			gomega.Expect(err).Should(gomega.BeNil())
		}
		ginkgo.By("init case path")
		casesPath = []string{
			util.GetCasesFileAbsPath("unshard/case/test_join_nosharding.sql"),
			util.GetCasesFileAbsPath("unshard/case/test_join.sql"),
			util.GetCasesFileAbsPath("unshard/case/test_select_global_old.sql"),
			util.GetCasesFileAbsPath("unshard/case/test_simple.sql"),
			util.GetCasesFileAbsPath("unshard/case/test_subquery_global.sql"),
		}
	})
	ginkgo.Context("unshard support test", func() {
		ginkgo.It("When testing sql support", func() {
			_, err := gaeaMaster.Exec("use sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = gaeaDB.Exec("use sbtest1")
			gomega.Expect(err).Should(gomega.BeNil())
			for _, file := range casesPath {
				sqls, err := util.GetSqlFromFile(file)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, v := range sqls {
					operationType, err := util.GetSqlType(v)
					gomega.Expect(err).Should(gomega.BeNil())
					if operationType == util.UnKnown {
						util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/unknown.sql"))
						continue
					}
					res1, err1 := util.ExecuteSQL(gaeaMaster, operationType, v)
					if err1 == util.ErrSkippedSQLOperation || err1 == util.ErrUnsupportedSQLOperation {
						util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/gaeaSkip.sql"))
					}
					res2, err2 := util.ExecuteSQL(gaeaDB, operationType, v)
					if err2 == util.ErrSkippedSQLOperation || err2 == util.ErrUnsupportedSQLOperation {
						util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/mysqlSkip.sql"))
						continue
					}
					if err1 != nil {
						gomega.Expect(fmt.Errorf("mysql Exec%s", err1.Error())).Should(gomega.BeNil())
					}
					if err2 != nil {
						gomega.Expect(fmt.Errorf("gaea Exec%s", err2.Error())).Should(gomega.BeNil())
					}
					equal, err := util.CompareDBExecResults(res1, res2)
					gomega.Expect(err).Should(gomega.BeNil())
					if equal {
						util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/pass.sql"))
					} else {
						util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/unpass.sql"))
					}
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		setupFile := util.GetCasesFileAbsPath("unshard/case/1-clean.sql")
		sqls, err := util.GetSqlFromFile(setupFile)
		gomega.Expect(err).Should(gomega.BeNil())
		for _, v := range sqls {
			operationType, err := util.GetSqlType(v)
			if operationType == util.UnKnown {
				util.OutPutResult(v, util.GetCasesFileAbsPath("unshard/result/unknown.sql"))
				continue
			}
			gomega.Expect(err).Should(gomega.BeNil())
			_, err = util.ExecuteSQL(gaeaMaster, "Exec", v)
			gomega.Expect(err).Should(gomega.BeNil())
		}
		// 删除注册
		config.GetNamespaceRegisterManager().UnRegisterNamespaces()
	})
})
