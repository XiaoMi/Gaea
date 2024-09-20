// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/XiaoMi/Gaea/models"
)

func TestGetRealDatabases(t *testing.T) {
	tests := []struct {
		databaseList     []string
		realDatabaseList []string
		err              error
	}{
		{[]string{}, []string{}, nil},
		{[]string{"db0"}, []string{"db0"}, nil},
		{[]string{"db0", "db1"}, []string{"db0", "db1"}, nil},
		{[]string{"db0", "db1", "db2"}, []string{"db0", "db1", "db2"}, nil},
		{[]string{"db[0-1]", "db[2-5]"}, []string{"db0", "db1", "db2", "db3", "db4", "db5"}, nil},
		{[]string{"db[0-1]", "db2", "db3"}, []string{"db0", "db1", "db2", "db3"}, nil},
		{[]string{"db0", "db[1-2]", "db3"}, []string{"db0", "db1", "db2", "db3"}, nil},

		{[]string{"db0", "db[1]"}, []string{"db0", "db[1]"}, nil},
		{[]string{"db0", "db[1-1]"}, nil, fmt.Errorf("invalid bound value of database list: db[1-1]")},
		{[]string{"db0", "db[1-0]"}, nil, fmt.Errorf("invalid bound value of database list: db[1-0]")},
		{[]string{"db0", "db[a-0]"}, []string{"db0", "db[a-0]"}, nil},
		{[]string{"db0", "db[0-a]"}, []string{"db0", "db[0-a]"}, nil},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.databaseList), func(t *testing.T) {
			dbList, err := getRealDatabases(test.databaseList)
			if err != nil && test.err != nil {
				if err.Error() != test.err.Error() {
					t.Errorf("err not equal, expect: %v, actual: %v", test.err, err)
					t.FailNow()
				}
			} else if (err != nil && test.err == nil) || (err == nil && test.err != nil) {
				t.Errorf("err not equal, expect: %v, actual: %v", test.err, err)
				t.FailNow()
			} else {
				if len(dbList) != len(test.realDatabaseList) {
					t.Errorf("result not equal, expect: %v, actual: %v", test.realDatabaseList, dbList)
					t.FailNow()
				}
				for i := 0; i < len(dbList); i++ {
					if dbList[i] != test.realDatabaseList[i] {
						t.Errorf("result not equal, expect: %v, actual: %v", test.realDatabaseList, dbList)
						t.FailNow()
					}
				}
			}
		})
	}
}

func TestParseMycatRule(t *testing.T) {
	var s = `
	{
		"name": "gaea_namespace_1",
		"online": true,
		"read_only": true,
		"allowed_dbs": {
			"db1": true,
			"db2": true
		},
		"slices": [
			{
				"Name": "slice-0",
				"UserName": "root",
				"Password": "root",
				"Master": "127.0.0.1:3306",
				"Slaves": [
					"127.0.0.1:3306",
					"127.0.0.1:3306"
				],
				"MaxConnNum": 128,
				"DownAfterNoAlive": 16
			},
			{
				"Name": "slice-1",
				"UserName": "root",
				"Password": "root",
				"Master": "127.0.0.1:3307",
				"Slaves": [
					"127.0.0.1:3307",
					"127.0.0.1:3307"
				],
				"MaxConnNum": 128,
				"DownAfterNoAlive": 16
			}
		],
		"shard_rules": [
			{
				"db": "gaea",
				"table": "test_shard_mycat_mod",
				"type": "mycat_mod",
				"key": "id",
				"locations": [
					1,
					1
				],
				"slices": [
					"slice-0",
					"slice-1"
				],
				"databases": [
					"gaea_0",
					"gaea_1"
				],
				"default_database": "gaea_0"
			},
			{
				"db": "gaea",
				"table": "test_shard_mycat_long",
				"type": "mycat_long",
				"key": "id",
				"locations": [
					1,
					1
				],
				"slices": [
					"slice-0",
					"slice-1"
				],
				"databases": [
					"gaea_0",
					"gaea_1"
				],
				"default_database": "gaea_0",
				"partition_count": "1,1",
				"partition_length": "256,768"
			},
			{
				"db": "gaea",
				"table": "test_shard_mycat_murmur",
				"type": "mycat_murmur",
				"key": "id",
				"locations": [
					1,
					1
				],
				"slices": [
					"slice-0",
					"slice-1"
				],
				"databases": [
					"gaea_0",
					"gaea_1"
				],
				"default_database": "gaea_0",
				"seed": "1",
				"virtual_bucket_times": "160"
			}
		],
		"users": [
			{
				"UserName": "test_shard_hash",
				"Password": "test_shard_hash",
				"Namespace": "gaea_namespace_1",
				"rw_flag": 2,
				"rw_split": 1
			}
		],
		"default_slice": "slice-0"
	}
	`

	var namespace = new(models.Namespace)
	if err := models.JSONDecode(namespace, []byte(s)); err != nil {
		t.Fatal(err)
	}

	rt, err := NewRouter(namespace)
	if err != nil {
		t.Fatal(err)
	}
	if rt.defaultRule.GetSlice(0) != "slice-0" {
		t.Fatal("default rule parse not correct.")
	}

	mycatModRule := rt.GetRule("gaea", "test_shard_mycat_mod")
	if mycatModRule.GetType() != MycatModRuleType {
		t.Fatal(mycatModRule.GetType())
	}

	if len(mycatModRule.GetSlices()) != 2 || mycatModRule.GetSlice(0) != "slice-0" || mycatModRule.GetSlices()[1] != "slice-1" {
		t.Fatal("parse slices not correct.")
	}

	mycatLongRule := rt.GetRule("gaea", "test_shard_mycat_long")
	if mycatLongRule.GetType() != MycatLongRuleType {
		t.Fatal(mycatLongRule.GetType())
	}

	if len(mycatLongRule.GetSlices()) != 2 ||
		mycatLongRule.GetSlice(0) != "slice-0" || mycatLongRule.GetSlices()[1] != "slice-1" {
		t.Fatal("parse slices not correct.")
	}

	mycatMurmurRule := rt.GetRule("gaea", "test_shard_mycat_murmur")
	if mycatMurmurRule.GetType() != MycatMurmurRuleType {
		t.Fatal(mycatLongRule.GetType())
	}

	if len(mycatMurmurRule.GetSlices()) != 2 || mycatMurmurRule.GetSlice(0) != "slice-0" || mycatMurmurRule.GetSlices()[1] != "slice-1" {
		t.Fatal("parse slices not correct.")
	}
}

// TODO YYYY-MM-DD HH:MM:SS,YYYY-MM-DD test
func TestParseDateRule(t *testing.T) {
	var s = `
	{"name": "gaea_namespace_1",
	"online":true,
	"read_only":true,
	"allowed_dbs": {"db1":true,
					"db2":true},
	"slices":[
	   {
		   "Name": "slice-0",
		   "UserName": "root",
		   "Password": "root",
		   "Master": "127.0.0.1:3306",
		   "Slaves": [
			   "127.0.0.1:3306",
			   "127.0.0.1:3306"
		   ],
		   "MaxConnNum": 128,
		   "DownAfterNoAlive": 16
	   },
	   {
		   "Name": "slice-1",
		   "UserName": "root",
		   "Password": "root",
		   "Master": "127.0.0.1:3307",
		   "Slaves": [
			   "127.0.0.1:3307",
			   "127.0.0.1:3307"
		   ],
		   "MaxConnNum": 128,
		   "DownAfterNoAlive": 16
	   }
	],
	 "shard_rules": [
		 {
			 "db": "gaea",
			 "table": "test_shard_year",
			 "type": "date_year",
			 "key": "date",
			 "slices": [
				 "slice-0",
				 "slice-1"
			 ],
			 "date_range": ["2012-2015","2016-2018"]
		 },
		 {
			 "db": "gaea",
			 "table": "test_shard_month",
			 "type": "date_month",
			 "key": "date",
			 "slices": [
				 "slice-0",
				 "slice-1"
			 ],
			 "date_range": ["201512-201603", "201604-201608"]
		 },
		 {
			 "db": "gaea",
			 "table": "test_shard_day",
			 "type": "date_day",
			 "key": "date",
			 "slices": [
				 "slice-0",
				 "slice-1"
			 ],
			 "date_range": ["20151201-20160122", "20160202-20160308"]
		 }
     ],
	 "users": [
		 {
			 "UserName": "test_shard_hash",
			 "Password": "test_shard_hash",
			 "Namespace": "gaea_namespace_1",
			 "rw_flag": 2,
			 "rw_split": 1
		 }
	 ],
	 "default_slice": "slice-0"
	}
	`

	var namespace = new(models.Namespace)
	if err := models.JSONDecode(namespace, []byte(s)); err != nil {
		t.Fatal(err)
	}

	rt, err := NewRouter(namespace)
	if err != nil {
		t.Fatal(err)
	}
	if rt.defaultRule.GetSlice(0) != "slice-0" {
		t.Fatal("default rule parse not correct.")
	}

	yearRule := rt.GetRule("gaea", "test_shard_year")
	if yearRule.GetType() != DateYearRuleType {
		t.Fatal(yearRule.GetType())
	}

	if len(yearRule.GetSlices()) != 2 || yearRule.GetSlice(0) != "slice-0" || yearRule.GetSlices()[1] != "slice-1" {
		t.Fatal("parse slices not correct.")
	}

	monthRule := rt.GetRule("gaea", "test_shard_month")
	if monthRule.GetType() != DateMonthRuleType {
		t.Fatal(monthRule.GetType())
	}

	dayRule := rt.GetRule("gaea", "test_shard_day")
	if dayRule.GetType() != DateDayRuleType {
		t.Fatal(monthRule.GetType())
	}
}

func TestParseRule(t *testing.T) {
	var s = `
	{"name": "gaea_namespace_1",
	"online":true,
	"read_only":true,
	"allowed_dbs": {"db1":true,
					"db2":true},
	"slices":[
	   {
		   "Name": "slice-0",
		   "UserName": "root",
		   "Password": "root",
		   "Master": "127.0.0.1:3306",
		   "Slaves": [
			   "127.0.0.1:3306",
			   "127.0.0.1:3306"
		   ],
		   "MaxConnNum": 128,
		   "DownAfterNoAlive": 16
	   },
	   {
		   "Name": "slice-1",
		   "UserName": "root",
		   "Password": "root",
		   "Master": "127.0.0.1:3307",
		   "Slaves": [
			   "127.0.0.1:3307",
			   "127.0.0.1:3307"
		   ],
		   "MaxConnNum": 128,
		   "DownAfterNoAlive": 16
	   }
	],
	 "shard_rules": [
		 {
			 "db": "gaea",
			 "table": "test_shard_hash",
			 "type": "hash",
			 "key": "id",
			 "locations": [
				 1,
				 1
			 ],
			 "slices": [
				 "slice-0",
				 "slice-1"
			 ],
			 "date_range": null,
			 "table_row_limit": 0
		 },
		 {
			 "db": "gaea",
			 "table": "test_shard_range",
			 "type": "range",
			 "key": "id",
			 "locations": [
				 1,
				 1
			 ],
			 "slices": [
				 "slice-0",
				 "slice-1"
			 ],
			 "date_range": null,
			 "table_row_limit": 10000
		 }
     ],
	 "users": [
		 {
			 "UserName": "test1",
			 "Password": "test1",
			 "Namespace": "gaea_namespace_1",
			 "rw_flag": 2,
			 "rw_split": 1
		 }
	 ],
	 "default_slice": "slice-0"
	}
`
	var namespace = new(models.Namespace)
	if err := models.JSONDecode(namespace, []byte(s)); err != nil {
		t.Fatal(err)
	}

	rt, err := NewRouter(namespace)
	if err != nil {
		t.Fatal(err)
	}
	if rt.defaultRule.GetSlice(0) != "slice-0" {
		t.Fatal("default rule parse not correct.")
	}

	rt.GetRule("", "gaea.test_shard_hash")

	hashRule := rt.GetRule("gaea", "test_shard_hash")
	if hashRule.GetType() != HashRuleType {
		t.Fatal(hashRule.GetType())
	}

	if len(hashRule.GetSlices()) != 2 || hashRule.GetSlice(0) != "slice-0" || hashRule.GetSlice(1) != "slice-1" {
		t.Fatal("parse slices not correct.")
	}

	rangeRule := rt.GetRule("gaea", "test_shard_range")
	if rangeRule.GetType() != RangeRuleType {
		t.Fatal(rangeRule.GetType())
	}

	defaultRule := rt.GetRule("gaea", "defaultRule_table")
	if defaultRule == nil {
		t.Fatal("must not nil")
	}

	if defaultRule.GetType() != DefaultRuleType {
		t.Fatal(defaultRule.GetType())
	}

	if defaultRule.GetShard() == nil {
		t.Fatal("nil error")
	}
}

type MockShard struct {
	FindForKeyFunc func(key interface{}) (int, error)
}

func (m *MockShard) FindForKey(key interface{}) (int, error) {
	return m.FindForKeyFunc(key)
}

type MockRule struct {
	GetDBFunc                       func() string
	GetTableFunc                    func() string
	GetShardingColumnFunc           func() string
	IsLinkedRuleFunc                func() bool
	GetShardFunc                    func() Shard
	FindTableIndexFunc              func(key interface{}) (int, error)
	GetSliceFunc                    func(i int) string
	GetSliceIndexFromTableIndexFunc func(i int) int
	GetSlicesFunc                   func() []string
	GetSubTableIndexesFunc          func() []int
	GetTypeFunc                     func() string
	GetDatabaseNameByTableIndexFunc func(index int) (string, error)
	GetTableIndexByDatabaseNameFunc func(phyDB string) (int, bool)
	GetDatabasesFunc                func() []string
	GetFirstTableIndexFunc          func() int
	GetLastTableIndexFunc           func() int
}

func (m *MockRule) GetDB() string {
	return m.GetDBFunc()
}

func (m *MockRule) GetTable() string {
	return m.GetTableFunc()
}

func (m *MockRule) GetShardingColumn() string {
	return m.GetShardingColumnFunc()
}

func (m *MockRule) IsLinkedRule() bool {
	return m.IsLinkedRuleFunc()
}

func (m *MockRule) GetShard() Shard {
	return m.GetShardFunc()
}

func (m *MockRule) FindTableIndex(key interface{}) (int, error) {
	return m.FindTableIndexFunc(key)
}

func (m *MockRule) GetSlice(i int) string {
	return m.GetSliceFunc(i)
}

func (m *MockRule) GetSliceIndexFromTableIndex(i int) int {
	return m.GetSliceIndexFromTableIndexFunc(i)
}

func (m *MockRule) GetSlices() []string {
	return m.GetSlicesFunc()
}

func (m *MockRule) GetSubTableIndexes() []int {
	return m.GetSubTableIndexesFunc()
}

func (m *MockRule) GetType() string {
	return m.GetTypeFunc()
}

func (m *MockRule) GetDatabaseNameByTableIndex(index int) (string, error) {
	return m.GetDatabaseNameByTableIndexFunc(index)
}

func (m *MockRule) GetTableIndexByDatabaseName(phyDB string) (int, bool) {
	return m.GetTableIndexByDatabaseNameFunc(phyDB)
}

func (m *MockRule) GetDatabases() []string {
	return m.GetDatabasesFunc()
}
func (m *MockRule) GetFirstTableIndex() int {
	return m.GetFirstTableIndexFunc()
}
func (m *MockRule) GetLastTableIndex() int {
	return m.GetLastTableIndexFunc()
}

func TestCreateLinkedRule(t *testing.T) {
	rules := map[string]map[string]Rule{
		"db1": {
			"table1": &MockRule{},
		},
	}
	shard := &models.Shard{
		Type:        LinkedTableRuleType,
		DB:          "db1",
		Table:       "table2",
		ParentTable: "table1",
		Key:         "key1",
	}

	t.Run("valid creation", func(t *testing.T) {
		expectedLinkedRule := &LinkedRule{
			db:             "db1",
			table:          "table2",
			shardingColumn: "key1",
			linkToRule: &BaseRule{
				db:             "db1",
				table:          "table2",
				shardingColumn: "key1",
			},
		}
		rules["db1"]["table1"] = &BaseRule{
			db:             "db1",
			table:          "table2",
			shardingColumn: "key1",
		}
		linkedRule, err := createLinkedRule(rules, shard)
		assert.NoError(t, err)
		assert.Equal(t, expectedLinkedRule, linkedRule)
	})

	t.Run("invalid type", func(t *testing.T) {
		shard.Type = "other"
		_, err := createLinkedRule(rules, shard)
		assert.Error(t, err)
	})

	t.Run("db not found", func(t *testing.T) {
		shard.DB = "unknown"
		_, err := createLinkedRule(rules, shard)
		assert.Error(t, err)
	})

	t.Run("parent table not found", func(t *testing.T) {
		shard.DB = "db1"
		shard.ParentTable = "unknown"
		_, err := createLinkedRule(rules, shard)
		assert.Error(t, err)
	})

	t.Run("cannot link to another linked rule", func(t *testing.T) {
		mockLinkToRule := &MockRule{}
		mockLinkToRule.GetTypeFunc = func() string {
			return LinkedTableRuleType
		}
		_, err := createLinkedRule(rules, shard)
		assert.Error(t, err)
	})

	t.Run("must link to a base rule", func(t *testing.T) {
		mockLinkToRule := &MockRule{}
		mockLinkToRule.GetTypeFunc = func() string {
			return "other"
		}
		_, err := createLinkedRule(rules, shard)
		assert.Error(t, err)
	})
}

func TestBaseRuleMethods(t *testing.T) {
	baseRule := &BaseRule{
		db:              "db1",
		table:           "table1",
		shardingColumn:  "column1",
		shard:           &MockShard{},
		slices:          []string{"slice1", "slice2"},
		tableToSlice:    map[int]int{0: 0, 1: 1},
		subTableIndexes: []int{0, 1},
		ruleType:        "base",
		mycatDatabases:  []string{"db1", "db2"},
		mycatDatabaseToTableIndexMap: map[string]int{
			"db1": 0,
			"db2": 1,
		},
	}

	t.Run("GetDB", func(t *testing.T) {
		assert.Equal(t, "db1", baseRule.GetDB())
	})

	t.Run("GetTable", func(t *testing.T) {
		assert.Equal(t, "table1", baseRule.GetTable())
	})

	t.Run("GetShardingColumn", func(t *testing.T) {
		assert.Equal(t, "column1", baseRule.GetShardingColumn())
	})

	t.Run("IsLinkedRule", func(t *testing.T) {
		assert.False(t, baseRule.IsLinkedRule())
	})

	t.Run("GetShard", func(t *testing.T) {
		assert.NotNil(t, baseRule.GetShard())
	})

	t.Run("FindTableIndex", func(t *testing.T) {
		mockShard := &MockShard{
			FindForKeyFunc: func(key interface{}) (int, error) {
				return 0, nil
			},
		}
		baseRule.shard = mockShard
		index, err := baseRule.FindTableIndex("key")
		assert.NoError(t, err)
		assert.Equal(t, 0, index)
	})

	t.Run("GetSlice", func(t *testing.T) {
		assert.Equal(t, "slice1", baseRule.GetSlice(0))
	})

	t.Run("GetSliceIndexFromTableIndex", func(t *testing.T) {
		assert.Equal(t, 0, baseRule.GetSliceIndexFromTableIndex(0))
	})

	t.Run("GetSlices", func(t *testing.T) {
		assert.Equal(t, []string{"slice1", "slice2"}, baseRule.GetSlices())
	})

	t.Run("GetSubTableIndexes", func(t *testing.T) {
		assert.Equal(t, []int{0, 1}, baseRule.GetSubTableIndexes())
	})

	t.Run("GetFirstTableIndex", func(t *testing.T) {
		assert.Equal(t, 0, baseRule.GetFirstTableIndex())
	})

	t.Run("GetLastTableIndex", func(t *testing.T) {
		assert.Equal(t, 1, baseRule.GetLastTableIndex())
	})

	t.Run("GetType", func(t *testing.T) {
		assert.Equal(t, "base", baseRule.GetType())
	})

	t.Run("GetDatabaseNameByTableIndex", func(t *testing.T) {
		name, err := baseRule.GetDatabaseNameByTableIndex(0)
		assert.NoError(t, err)
		assert.Equal(t, "db1", name)
	})

	t.Run("GetTableIndexByDatabaseName", func(t *testing.T) {
		index, ok := baseRule.GetTableIndexByDatabaseName("db1")
		assert.True(t, ok)
		assert.Equal(t, 0, index)
	})

	t.Run("GetDatabases", func(t *testing.T) {
		assert.Equal(t, []string{"db1", "db2"}, baseRule.GetDatabases())
	})
}

func TestLinkedRuleMethods(t *testing.T) {
	baseRule := &BaseRule{
		db:              "db1",
		table:           "table1",
		shardingColumn:  "column1",
		shard:           &MockShard{},
		slices:          []string{"slice1", "slice2"},
		tableToSlice:    map[int]int{0: 0, 1: 1},
		subTableIndexes: []int{0, 1},
		ruleType:        "base",
		mycatDatabases:  []string{"db1", "db2"},
		mycatDatabaseToTableIndexMap: map[string]int{
			"db1": 0,
			"db2": 1,
		},
	}
	linkedRule := &LinkedRule{
		db:             "db1",
		table:          "table2",
		shardingColumn: "column1",
		linkToRule:     baseRule,
	}

	t.Run("GetDB", func(t *testing.T) {
		assert.Equal(t, "db1", linkedRule.GetDB())
	})

	t.Run("GetTable", func(t *testing.T) {
		assert.Equal(t, "table2", linkedRule.GetTable())
	})

	t.Run("GetParentDB", func(t *testing.T) {
		assert.Equal(t, "db1", linkedRule.GetParentDB())
	})

	t.Run("GetParentTable", func(t *testing.T) {
		assert.Equal(t, "table1", linkedRule.GetParentTable())
	})

	t.Run("GetShardingColumn", func(t *testing.T) {
		assert.Equal(t, "column1", linkedRule.GetShardingColumn())
	})

	t.Run("IsLinkedRule", func(t *testing.T) {
		assert.True(t, linkedRule.IsLinkedRule())
	})

	t.Run("GetShard", func(t *testing.T) {
		assert.NotNil(t, linkedRule.GetShard())
	})

	t.Run("FindTableIndex", func(t *testing.T) {
		mockShard := &MockShard{
			FindForKeyFunc: func(key interface{}) (int, error) {
				return 0, nil
			},
		}
		baseRule.shard = mockShard
		index, err := linkedRule.FindTableIndex("key")
		assert.NoError(t, err)
		assert.Equal(t, 0, index)
	})

	t.Run("GetFirstTableIndex", func(t *testing.T) {
		assert.Equal(t, 0, linkedRule.GetFirstTableIndex())
	})

	t.Run("GetLastTableIndex", func(t *testing.T) {
		assert.Equal(t, 1, linkedRule.GetLastTableIndex())
	})

	t.Run("GetSlice", func(t *testing.T) {
		assert.Equal(t, "slice1", linkedRule.GetSlice(0))
	})

	t.Run("GetSliceIndexFromTableIndex", func(t *testing.T) {
		assert.Equal(t, 0, linkedRule.GetSliceIndexFromTableIndex(0))
	})

	t.Run("GetSlices", func(t *testing.T) {
		assert.Equal(t, []string{"slice1", "slice2"}, linkedRule.GetSlices())
	})

	t.Run("GetSubTableIndexes", func(t *testing.T) {
		assert.Equal(t, []int{0, 1}, linkedRule.GetSubTableIndexes())
	})

	t.Run("GetType", func(t *testing.T) {
		assert.Equal(t, "base", linkedRule.GetType())
	})

	t.Run("GetDatabaseNameByTableIndex", func(t *testing.T) {
		name, err := linkedRule.GetDatabaseNameByTableIndex(0)
		assert.NoError(t, err)
		assert.Equal(t, "db1", name)
	})

	t.Run("GetDatabases", func(t *testing.T) {
		assert.Equal(t, []string{"db1", "db2"}, linkedRule.GetDatabases())
	})

	t.Run("GetTableIndexByDatabaseName", func(t *testing.T) {
		index, ok := linkedRule.GetTableIndexByDatabaseName("db1")
		assert.True(t, ok)
		assert.Equal(t, 0, index)
	})
}
