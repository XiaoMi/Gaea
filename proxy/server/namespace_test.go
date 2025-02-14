package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type phyDBCase struct {
	defaultPhyDBs map[string]string
	allowedDBs    map[string]bool
	shardRules    []*models.Shard
	realPhyDBs    map[string]string
}

func TestParsePhyDBs(t *testing.T) {
	tests := []phyDBCase{
		{defaultPhyDBs: map[string]string{"db_mycat": "db_mycat"},
			allowedDBs: map[string]bool{"db_mycat": true},
			shardRules: []*models.Shard{{Databases: []string{"db_mycat_[0-1]"}}},
			realPhyDBs: map[string]string{"db_mycat": "db_mycat", "db_mycat_0": "db_mycat_0", "db_mycat_1": "db_mycat_1", "information_schema": "information_schema"}},
		{defaultPhyDBs: map[string]string{},
			allowedDBs: map[string]bool{"db_mycat": true},
			shardRules: []*models.Shard{},
			realPhyDBs: map[string]string{"db_mycat": "db_mycat"}},
	}
	for index, test := range tests {
		t.Run("test", func(t *testing.T) {
			realPhyDBs, _ := parseDefaultPhyDB(test.defaultPhyDBs, test.allowedDBs, test.shardRules)
			if !reflect.DeepEqual(realPhyDBs, test.realPhyDBs) {
				t.Errorf("test %d, parse real phyDBs error, %v", index, realPhyDBs)
			}
		})
	}
}

// TestParseSlice 测试 parseSlice 方法，确保返回的 Slice 结构体 Master、Slave、StatisticSlave 解析正确，不为 nil
func TestParseSlice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// 配置测试数据
	testCases := []struct {
		name                 string
		cfg                  *models.Slice
		expectMasterNodes    int
		expectSlaveNodes     int
		expectStatSlaveNodes int
		expectErr            bool
	}{
		{
			name: "valid master and slaves",
			cfg: &models.Slice{
				Name:                        "slice-0",
				Master:                      "master-db:3306",
				Slaves:                      []string{"slave-db-1:3306", "slave-db-2:3307"},
				StatisticSlaves:             []string{"stat-db-1:3308"},
				FallbackToMasterOnSlaveFail: "on",
				HandshakeTimeout:            3000, // 3s
				HealthCheckSql:              "SELECT 1",
			},
			expectMasterNodes:    1, // Master 存在
			expectSlaveNodes:     2, // 有 2 个 Slave
			expectStatSlaveNodes: 1, // 有 1 个 Statistic Slave
			expectErr:            false,
		},
		{
			name: "valid master and no slaves",
			cfg: &models.Slice{
				Name:            "slice-1",
				Master:          "master-db:3306",
				Slaves:          []string{},
				StatisticSlaves: []string{"stat-db-1:3308"},
			},
			expectMasterNodes:    1, // Master 不为空，但 `s.Master` 应该被初始化
			expectSlaveNodes:     0, // 有 0 个 Slave
			expectStatSlaveNodes: 1, // 有 1 个 Statistic Slave
			expectErr:            false,
		},
		{
			name: "valid master and no slaves",
			cfg: &models.Slice{
				Name:            "slice-1",
				Master:          "master-db:3306",
				Slaves:          []string{"slave-db-1:3306"},
				StatisticSlaves: []string{},
			},
			expectMasterNodes:    1, // Master 不为空，但 `s.Master` 应该被初始化
			expectSlaveNodes:     1, // 有 1 个 Slave
			expectStatSlaveNodes: 0, // 有 0 个 Statistic Slave
			expectErr:            false,
		},
		{
			name: "empty master should still return non-nil Slice",
			cfg: &models.Slice{
				Name:            "slice-1",
				Slaves:          []string{"slave-db-1:3306", "slave-db-2:3307"},
				StatisticSlaves: []string{"stat-db-1:3308"},
			},
			expectMasterNodes:    0, // Master 为空，但 `s.Master` 应该被初始化
			expectSlaveNodes:     2, // 有 2 个 Slave
			expectStatSlaveNodes: 1, // 有 1 个 Statistic Slave
			expectErr:            false,
		},
		{
			name: "empty master and no slaves and statistic slaves ",
			cfg: &models.Slice{
				Name:            "slice-1",
				Slaves:          []string{},
				StatisticSlaves: []string{},
			},
			expectMasterNodes:    0, // Master 为空，但 `s.Master` 应该被初始化
			expectSlaveNodes:     0, // 有 0 个 Slave
			expectStatSlaveNodes: 0, // 有 0 个 Statistic Slave
			expectErr:            false,
		},
	}

	// 运行测试
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := parseSlice(tc.cfg, "utf8mb4", 45, "bj")

			// 如果预期失败，检查错误是否符合预期
			if tc.expectErr {
				assert.NotNil(t, err, "expected an error but got nil")
				return
			}

			assert.Nil(t, err, "unexpected error: %v", err)
			assert.NotNil(t, s, "Slice should not be nil")
			assert.NotNil(t, s.Master, "Master should not be nil")
			assert.NotNil(t, s.Slave, "Slave should not be nil")
			assert.NotNil(t, s.StatisticSlave, "StatisticSlave should not be nil")

			// 验证 Master
			assert.Equal(t, tc.expectMasterNodes, len(s.Master.Nodes), "Master Nodes count mismatch")
			// 验证 Slaves
			assert.Equal(t, tc.expectSlaveNodes, len(s.Slave.Nodes), "Slave Nodes count mismatch")
			// 验证 Statistic Slaves
			assert.Equal(t, tc.expectStatSlaveNodes, len(s.StatisticSlave.Nodes), "StatisticSlave Nodes count mismatch")

			// 确保 ProxyDatacenter、HealthCheckSql 和 HandshakeTimeout 被正确初始化
			assert.Equal(t, "bj", s.ProxyDatacenter)
			assert.Equal(t, tc.cfg.HealthCheckSql, s.HealthCheckSql)
			assert.Equal(t, time.Duration(tc.cfg.HandshakeTimeout)*time.Millisecond, s.HandshakeTimeout)
		})
	}
}
