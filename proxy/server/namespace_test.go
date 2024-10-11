package server

import (
	"github.com/XiaoMi/Gaea/models"
	"reflect"
	"testing"
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
