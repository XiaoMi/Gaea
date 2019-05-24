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

//TODO YYYY-MM-DD HH:MM:SS,YYYY-MM-DD test
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
