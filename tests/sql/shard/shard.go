package shard

import (
	"fmt"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	defaultSlice = []string{
		"slice-0",
		"slice-1",
	}
	defaultLocation = []int{
		2,
		2,
	}

	kingshardHashRule = config.NewShard(
		config.WithDB("db_kingshard_hash"),
		config.WithTable("tbl_shard"),
		config.WithType("hash"),
		config.WithKey("id"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice(defaultSlice),
	)
	kingshardModRule = config.NewShard(
		config.WithDB("db_kingshard_mod"),
		config.WithTable("tbl_shard"),
		config.WithType("mod"),
		config.WithKey("id"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice(defaultSlice),
	)

	kingshardRangeRule = config.NewShard(
		config.WithDB("db_kingshard_range"),
		config.WithTable("tbl_shard"),
		config.WithType("range"),
		config.WithKey("id"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice(defaultSlice),
		config.WithTableRowLimit(3),
	)

	kingshardDateYearRule = config.NewShard(
		config.WithDB("db_kingshard_date_year"),
		config.WithTable("tbl_shard"),
		config.WithType("date_year"),
		config.WithKey("create_time"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice(defaultSlice),
		config.WithDateRange([]string{
			"2016-2017",
			"2018-2019",
		}),
	)

	kingshardDateMonthRule = config.NewShard(
		config.WithDB("db_kingshard_date_month"),
		config.WithTable("tbl_shard"),
		config.WithType("date_month"),
		config.WithKey("create_time"),
		config.WithLocations([]int{
			2,
			2,
		}),
		config.WithShardSlice([]string{
			"slice-0",
			"slice-1",
		}),
		config.WithDateRange([]string{
			"201405-201406",
			"201408-201409",
		}),
	)

	kingshardDateDayRule = config.NewShard(
		config.WithDB("db_kingshard_date_day"),
		config.WithTable("tbl_shard"),
		config.WithType("date_day"),
		config.WithKey("create_time"),
		config.WithLocations([]int{
			2,
			2,
		}),
		config.WithShardSlice(defaultSlice),
		config.WithDateRange([]string{
			"20201201-20201202",
			"20201203-20201204",
		}),
		config.WithTableRowLimit(3),
	)
	mycatshardMod = config.NewShard(
		config.WithDB("db_mycat_mod"),
		config.WithTable("tbl_mycat"),
		config.WithType("mycat_mod"),
		config.WithKey("id"),
		config.WithLocations([]int{
			2,
			2,
		}),
		config.WithShardSlice([]string{
			"slice-0",
			"slice-1",
		}),
		config.WithDatabases([]string{
			"db_mycat_mod_[0-3]",
		}),
	)
	mycatshardLong = config.NewShard(
		config.WithDB("db_mycat_long"),
		config.WithTable("tbl_mycat"),
		config.WithType("mycat_long"),
		config.WithKey("id"),
		config.WithLocations([]int{
			2,
			2,
		}),
		config.WithShardSlice([]string{
			"slice-0",
			"slice-1",
		}),
		config.WithDatabases([]string{
			"db_mycat_long_[0-3]",
		}),
		config.WithPartitionCount("4"),
		config.WithPartitionLength("256"),
	)
	mycatshardMurmur = config.NewShard(
		config.WithDB("db_mycat_murmur"),
		config.WithTable("tbl_mycat"),
		config.WithParentTable(""),
		config.WithType("mycat_murmur"),
		config.WithKey("id"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice([]string{
			"slice-0",
			"slice-1",
		}),
		config.WithVirtualBucketTimes("160"),
		config.WithDatabases([]string{
			"db_mycat_murmur_[0-3]",
		}),
		config.WithSeed("0"),
	)
	mycatshardString = config.NewShard(
		config.WithDB("db_mycat_string"),
		config.WithTable("tbl_mycat"),
		config.WithParentTable(""),
		config.WithType("mycat_string"),
		config.WithKey("col1"),
		config.WithLocations(defaultLocation),
		config.WithShardSlice([]string{
			"slice-0",
			"slice-1",
		}),
		config.WithTableRowLimit(0),
		config.WithDatabases([]string{
			"db_mycat_string_[0-3]",
		}),
		config.WithPartitionCount("4"),
		config.WithPartitionLength("256"),
		config.WithHashSlice(":"),
	)
)

var _ = ginkgo.Describe("shard_dml_support_test", func() {

	var planManagers []*config.PlanManager
	var err error
	var allowedDBS map[string]bool
	var slice []*models.Slice
	var users []*models.User
	var noshardPlanNamespace *models.Namespace
	ginkgo.BeforeEach(func() {
		// 创建数据库实例
		ginkgo.By("create slice")
		slice = []*models.Slice{
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
			{
				Name:            "slice-1",
				UserName:        "superroot",
				Password:        "superroot",
				Master:          "127.0.0.1:3369",
				Slaves:          []string{},
				StatisticSlaves: []string{},
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
				InitConnect:     "",
			},
		}

		allowedDBS = map[string]bool{
			"db_kingshard_hash":       true,
			"db_kingshard_mod":        true,
			"db_kingshard_range":      true,
			"db_kingshard_date_year":  true,
			"db_kingshard_date_month": true,
			"db_kingshard_date_day":   true,
			"db_mycat_mod":            true,
			"db_mycat_long":           true,
			"db_mycat_murmur":         true,
			"db_mycat_string":         true,
		}
		//  创建数据库用户
		ginkgo.By("create user")
		users = []*models.User{
			{
				UserName:      "superroot",
				Password:      "superroot",
				Namespace:     "test_namespace_noshard_plan",
				RWFlag:        2,
				RWSplit:       1,
				OtherProperty: 0,
			},
		}
		ginkgo.By("create namespace")
		noshardPlanNamespace = config.NewNamespace(
			config.WithIsEncrypt(true),
			config.WithName("test_namespace_noshard_plan"),
			config.WithAllowedDBS(allowedDBS),
			config.WithSlices(slice),
			config.WithUsers(users),
			config.WithDefaultSlice("slice-0"),
			config.WithShardRules(
				[]*models.Shard{
					kingshardHashRule,
					kingshardModRule,
					kingshardRangeRule,
					kingshardDateYearRule,
					kingshardDateMonthRule,
					kingshardDateDayRule,
					mycatshardMod,
					mycatshardLong,
					mycatshardMurmur,
					mycatshardString,
				},
			),
		)
		// 注册命令空间
		ginkgo.By("add namespace")
		err = config.GetNamespaceRegisterManager().AddNamespace(noshardPlanNamespace)
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
			util.GetCasesFileAbsPath("shard/case/kingdate_year.json"),
			util.GetCasesFileAbsPath("shard/case/kingdate_month.json"),
			util.GetCasesFileAbsPath("shard/case/kinghash.json"),
			util.GetCasesFileAbsPath("shard/case/kingmod.json"),
			util.GetCasesFileAbsPath("shard/case/kingrange.json"),
			util.GetCasesFileAbsPath("shard/case/mycatmod.json"),
			util.GetCasesFileAbsPath("shard/case/mycatlong.json"),
			util.GetCasesFileAbsPath("shard/case/mycatstring.json"),
			util.GetCasesFileAbsPath("shard/case/mycatmurmur.json"),
		}
		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := config.NewPlanManager(
				config.WithPlanPath(v),
				config.WithMysqlClusterConn(mysqlCluster),
				config.WithGaeaConn(gaeaDB),
			)
			ginkgo.By(fmt.Sprintf("plan %s init", v))
			planManagers = append(planManagers, p)
		}

	})
	ginkgo.Context("shard support test", func() {
		ginkgo.It("shard support", func() {
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
