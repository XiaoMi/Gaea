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
package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func defaultNamespace() *Namespace {
	n := &Namespace{
		IsEncrypt:        false,
		Name:             "default",
		Online:           true,
		ReadOnly:         true,
		AllowedDBS:       make(map[string]bool),
		DefaultPhyDBS:    make(map[string]string),
		SlowSQLTime:      "",
		BlackSQL:         nil,
		AllowedIP:        nil,
		Slices:           make([]*Slice, 0),
		ShardRules:       make([]*Shard, 0),
		Users:            make([]*User, 0),
		DefaultSlice:     "",
		GlobalSequences:  nil,
		DefaultCharset:   "",
		DefaultCollation: "",
	}
	return n
}

func TestNamespaceEncode(t *testing.T) {
	var namespace = &Namespace{Name: "gaea_namespace_1", Online: true, ReadOnly: true, AllowedDBS: make(map[string]bool), Slices: make([]*Slice, 0), ShardRules: make([]*Shard, 0), Users: make([]*User, 0), DefaultSlice: "slice-0"}

	slice0 := &Slice{Name: "slice-0", UserName: "root", Password: "root", Master: "127.0.0.1:3306", Slaves: []string{"127.0.0.1:3306", "127.0.0.1:3306"}, Capacity: 128, MaxCapacity: 128, IdleTimeout: 120}
	slice1 := &Slice{Name: "slice-1", UserName: "root", Password: "root", Master: "127.0.0.1:3306", Slaves: []string{"127.0.0.1:3306", "127.0.0.1:3306"}, Capacity: 128, MaxCapacity: 128, IdleTimeout: 120}
	namespace.Slices = append(namespace.Slices, slice0)
	namespace.Slices = append(namespace.Slices, slice1)

	namespace.AllowedDBS["db1"] = true
	namespace.AllowedDBS["db2"] = true

	shard1 := &Shard{DB: "gaea", Table: "test_shard_hash", Type: "hash", Key: "id", Locations: []int{1, 1}, Slices: []string{"slice-0", "slice-1"}}
	shard2 := &Shard{DB: "gaea", Table: "test_shard_range", Type: "range", Key: "id", Locations: []int{1, 1}, Slices: []string{"slice-0", "slice-1"}, TableRowLimit: 10000}
	namespace.ShardRules = append(namespace.ShardRules, shard1)
	namespace.ShardRules = append(namespace.ShardRules, shard2)

	user1 := &User{UserName: "test1", Password: "test1", Namespace: "gaea_namespace_1", RWFlag: 2, RWSplit: 1}
	namespace.Users = append(namespace.Users, user1)

	t.Logf(string(namespace.Encode()))
}

func TestEncrypt(t *testing.T) {
	key := "1234abcd5678efg*"
	var namespace = &Namespace{Name: "gaea_namespace_1", Online: true, ReadOnly: true, AllowedDBS: make(map[string]bool), Slices: make([]*Slice, 0), ShardRules: make([]*Shard, 0), Users: make([]*User, 0), DefaultSlice: "slice-0"}
	slice0 := &Slice{Name: "slice-0", UserName: "test1", Password: "fdsafdsa23423sx*123", Master: "127.0.0.1:3306", Slaves: []string{"127.0.0.1:3306", "127.0.0.1:3306"}, Capacity: 128, MaxCapacity: 128, IdleTimeout: 120}
	slice1 := &Slice{Name: "slice-1", UserName: "test2", Password: "fasd14-43828284s*", Master: "127.0.0.1:3306", Slaves: []string{"127.0.0.1:3306", "127.0.0.1:3306"}, Capacity: 128, MaxCapacity: 128, IdleTimeout: 120}
	namespace.Slices = append(namespace.Slices, slice0)
	namespace.Slices = append(namespace.Slices, slice1)
	user1 := &User{UserName: "test1", Password: "testfadsfafdsla234231", Namespace: "gaea_namespace_1", RWFlag: 2, RWSplit: 1}
	user2 := &User{UserName: "test2", Password: "test2fdsafw5r3234", Namespace: "gaea_namespace_1", RWFlag: 2, RWSplit: 1}
	namespace.Users = append(namespace.Users, user1)
	namespace.Users = append(namespace.Users, user2)
	err := namespace.Encrypt(key)
	if err != nil {
		t.Errorf("test namespace encrypt failed, %v", err)
	}
	err = namespace.Decrypt(key)
	if err != nil {
		t.Errorf("test namespace failed, %v", err)
	}
	t.Logf(string(namespace.Encode()))
}

func TestVerifyName_Success(t *testing.T) {
	n := defaultNamespace()
	if err := n.verifyName(); err != nil {
		t.Errorf("test verifyName failed, %v", err)
	}
}

func TestVerifyName_Error(t *testing.T) {
	nf := defaultNamespace()
	nf.Name = ""
	if err := nf.verifyName(); err == nil {
		t.Errorf("test verifyName should fail but pass, name: %v", nf.Name)
	}
}

func TestVerifyAllowDBS_Success(t *testing.T) {
	n := defaultNamespace()
	n.AllowedDBS["db1"] = true
	if err := n.verifyAllowDBS(); err != nil {
		t.Errorf("test verifyAllowDBS failed, %v", err)
	}

}
func TestVerifyAllowDBS_Error(t *testing.T) {
	nf := defaultNamespace()
	if err := nf.verifyAllowDBS(); err == nil {
		t.Errorf("test verifyAllowDBS should fail but pass, name: %v", nf.Name)
	}
}

func TestVerifyUsers_Success(t *testing.T) {
	n := defaultNamespace()
	u1 := &User{UserName: "u1", Namespace: n.Name, Password: "pw1", RWFlag: ReadOnly, RWSplit: NoReadWriteSplit, OtherProperty: 0}
	u2 := &User{UserName: "u2", Namespace: n.Name, Password: "pw2", RWFlag: ReadWrite, RWSplit: ReadWriteSplit, OtherProperty: StatisticUser}
	n.Users = append(n.Users, u1)
	n.Users = append(n.Users, u2)

	if err := n.verifyUsers(); err != nil {
		t.Errorf("test verifyUsers failed, %v", err)
	}
}

func TestVerifyUsers_Error(t *testing.T) {
	nf := defaultNamespace()
	uf1 := &User{UserName: "u1", Namespace: "someone", Password: "pw1", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	uf2 := &User{UserName: "u1", Namespace: nf.Name, Password: "pw2", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	uf3 := &User{UserName: "", Namespace: "", Password: "", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	nf.Users = append(nf.Users, uf1)
	nf.Users = append(nf.Users, uf2)
	nf.Users = append(nf.Users, uf3)

	if err := nf.verifyUsers(); err == nil {
		t.Errorf("test verifyUsers should fail but pass, users: %s", JSONEncode(nf.Users))
	}
}

func TestVerifySlowSQLTime_Success(t *testing.T) {
	n := defaultNamespace()
	ssts := []string{"", "10"}
	for _, sst := range ssts {
		n.SlowSQLTime = sst
		if err := n.verifySlowSQLTime(); err != nil {
			t.Errorf("test verifySlowSQLTime failed, %v", err)
		}
	}
}

func TestVerifySlowSQLTime_Error(t *testing.T) {
	nf := defaultNamespace()
	sstfs := []string{"-1", "10.0", "test"}
	for _, sst := range sstfs {
		nf.SlowSQLTime = sst
		if err := nf.verifySlowSQLTime(); err == nil {
			t.Errorf("test verifySlowSQLTime should fail but pass, sst: %v", nf.SlowSQLTime)
		}
	}
}

func TestVerifyDBs_Success(t *testing.T) {
	n := defaultNamespace()
	// no logic database mode
	if err := n.verifyDBs(); err != nil {
		t.Errorf("test verifyDBs failed, %v", err)
	}

	// logic database mode
	n.AllowedDBS["test1"] = true
	n.DefaultPhyDBS["test1"] = ""
	if err := n.verifyDBs(); err != nil {
		t.Errorf("test verifyDBs failed, %v", err)
	}
}

func TestVerifyDBs_Error(t *testing.T) {
	nf := defaultNamespace()
	// logic database mode
	nf.AllowedDBS["test1"] = true
	nf.DefaultPhyDBS["test2"] = ""
	if err := nf.verifyDBs(); err == nil {
		t.Errorf("test verifyDBs should fail but pass, allowedDBS: %v, defaultPhyDBS: %v", nf.AllowedDBS, nf.DefaultPhyDBS)
	}
}

func TestVerifyAllowIps_Success(t *testing.T) {
	n := defaultNamespace()
	n.AllowedIP = append(n.AllowedIP, "  ")
	n.AllowedIP = append(n.AllowedIP, "10.221.163.82")
	if err := n.verifyAllowIps(); err != nil {
		t.Errorf("test verifyAllowIps failed, %v", err)
	}
}

func TestVerifyAllowIps_Error(t *testing.T) {
	nf := defaultNamespace()
	var ipfs = []string{"test", "1.1.1"}
	for _, ipf := range ipfs {
		nf.AllowedIP = []string{ipf}
		if err := nf.verifyAllowIps(); err == nil {
			t.Errorf("test verifyAllowIps should fail but pass, %v", nf.AllowedIP)
		}
	}
}

func TestVerifyCharset_Success(t *testing.T) {
	n := defaultNamespace()
	var ccs = [][]string{[]string{"", ""}, []string{"big5", ""}, []string{"big5", "big5_chinese_ci"}}
	for _, cc := range ccs {
		n.DefaultCharset = cc[0]
		n.DefaultCollation = cc[1]
		if err := n.verifyCharset(); err != nil {
			t.Errorf("test verifyCharset failed, %v", err)
		}
	}
}

func TestVerifyCharset_Error(t *testing.T) {
	nf := defaultNamespace()
	var ccfs = [][]string{[]string{"", "test"}, []string{"test", ""}, []string{"big5", "test"}, []string{"big5", "latin2_czech_cs"}}
	for _, ccf := range ccfs {
		nf.DefaultCharset = ccf[0]
		nf.DefaultCollation = ccf[1]
		if err := nf.verifyCharset(); err == nil {
			t.Errorf("test verifyCharset should fail but pass, charset: %s, collation: %s", nf.DefaultCharset, nf.DefaultCollation)
		}
	}
}

func TestVerifySlices_Success(t *testing.T) {
	n := defaultNamespace()
	var slice1 = &Slice{Name: "slice1", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100}
	var slice2 = &Slice{Name: "slice2", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100}
	n.Slices = append(n.Slices, slice1)
	n.Slices = append(n.Slices, slice2)
	if err := n.verifySlices(); err != nil {
		t.Errorf("test verifySlices failed, %v", err)
	}
}

func TestVerifySlices_Error(t *testing.T) {
	nf := defaultNamespace()
	var slice1 = &Slice{Name: "slice1", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100}
	var slicefs = []*Slice{
		&Slice{Name: "", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "user", Password: "", Master: "", Slaves: []string{}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "user", Password: "", Master: "", Slaves: []string{""}, Capacity: 1, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 0, MaxCapacity: 1, IdleTimeout: 100},
		&Slice{Name: "slice1", UserName: "user", Password: "", Master: "1.1.1.1:1", Slaves: []string{"1.1.1.1:2"}, Capacity: 1, MaxCapacity: 0, IdleTimeout: 100},
	}
	for _, slicef := range slicefs {
		nf.Slices = append(nf.Slices, slice1)
		nf.Slices = append(nf.Slices, slicef)
		if err := nf.verifySlices(); err == nil {
			t.Errorf("test verifySlices should fail but pass, slices: %s", JSONEncode(nf.Slices))
		}
	}
}

func TestVerifyDefaultSlice_Success(t *testing.T) {
	n := defaultNamespace()
	n.Slices = append(n.Slices, &Slice{Name: "slice1"})
	var dss = []string{"", "slice1"}
	for _, ds := range dss {
		n.DefaultSlice = ds
		if err := n.verifyDefaultSlice(); err != nil {
			t.Errorf("test verifyDefaultSlice failed, %v", err)
		}
	}
}

func TestVerifyDefaultSlice_Error(t *testing.T) {
	nf := defaultNamespace()
	nf.Slices = append(nf.Slices, &Slice{Name: "slice1"})
	nf.DefaultSlice = "slice2"
	if err := nf.verifyDefaultSlice(); err == nil {
		t.Errorf("test verifyDefaultSlice should fail but pass, defaultSlice: %s", nf.DefaultSlice)
	}
}

func TestVerifyShardRules_Success(t *testing.T) {
	n := defaultNamespace()
	n.Slices = []*Slice{
		&Slice{Name: "slice-0", UserName: "root", Password: "root", Master: "127.0.0.1:3306", Capacity: 64, MaxCapacity: 128, IdleTimeout: 3600},
		&Slice{Name: "slice-1", UserName: "root", Password: "root", Master: "127.0.0.1:3307", Capacity: 64, MaxCapacity: 128, IdleTimeout: 3600},
	}
	n.ShardRules = []*Shard{
		&Shard{DB: "db_ks", Table: "tbl_ks", Type: "mod", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}},
		&Shard{DB: "db_ks", Table: "tbl_ks_child", Type: "linked", Key: "id", ParentTable: "tbl_ks"},
		&Shard{DB: "db_ks", Table: "tbl_ks_user_child", Type: "linked", Key: "user_id", ParentTable: "tbl_ks"},
		&Shard{DB: "db_ks", Table: "tbl_ks_global_one", Type: "global", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}},
		&Shard{DB: "db_ks", Table: "tbl_ks_global_two", Type: "global", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}},
		&Shard{DB: "db_ks", Table: "tbl_ks_range", Type: "range", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, TableRowLimit: 100},
		&Shard{DB: "db_ks", Table: "tbl_ks_year", Type: "date_year", Key: "create_time", Slices: []string{"slice-0", "slice-1"}, DateRange: []string{"2014-2017", "2018-2019"}},
		&Shard{DB: "db_ks", Table: "tbl_ks_month", Type: "date_month", Key: "create_time", Slices: []string{"slice-0", "slice-1"}, DateRange: []string{"201405-201406", "201408-201409"}},
		&Shard{DB: "db_ks", Table: "tbl_ks_day", Type: "date_day", Key: "create_time", Slices: []string{"slice-0", "slice-1"}, DateRange: []string{"20140901-20140905", "20140907-20140908"}},
		&Shard{DB: "db_mycat", Table: "tbl_mycat", Type: "mycat_mod", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_[0-3]"}},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_child", Type: "linked", ParentTable: "tbl_mycat", Key: "id"},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_user_child", Type: "linked", ParentTable: "tbl_mycat", Key: "user_id"},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_murmur", Type: "mycat_murmur", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_0", "db_mycat_1", "db_mycat_2", "db_mycat_3"}, Seed: "0", VirtualBucketTimes: "160"},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_long", Type: "mycat_long", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_[0-3]"}, PartitionCount: "4", PartitionLength: "256"},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_global_one", Type: "global", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_[0-3]"}},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_global_two", Type: "global", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_[0-3]"}},
		&Shard{DB: "db_mycat", Table: "tbl_mycat_string", Type: "mycat_string", Key: "id", Locations: []int{2, 2}, Slices: []string{"slice-0", "slice-1"}, Databases: []string{"db_mycat_[0-3]"}, PartitionCount: "4", PartitionLength: "256", HashSlice: "20"},
	}
	if err := n.verifyShardRules(); err != nil {
		t.Errorf("test verifyShardRules failed, %v", err)
	}
}

func TestVerifyShardRule_Error_common(t *testing.T) {
	nf := defaultNamespace()
	// slices not match
	nf.ShardRules = []*Shard{&Shard{Slices: []string{"slice"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}

	// db duplicate
	nf.ShardRules = []*Shard{
		&Shard{DB: "db", Table: "table", Type: ShardMod},
		&Shard{DB: "db", Table: "table", Type: ShardMod},
	}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardDefault(t *testing.T) {
	nf := defaultNamespace()
	nf.ShardRules = []*Shard{&Shard{Type: ShardDefault}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardLinked(t *testing.T) {
	nf := defaultNamespace()
	// without parent rules
	nf.ShardRules = []*Shard{
		&Shard{Type: ShardLinked, DB: "db1", Table: "table1", ParentTable: "table2"},
	}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}

	// link to another linkedRule
	nf.ShardRules = []*Shard{
		&Shard{Type: ShardLinked, DB: "db1", Table: "table1", ParentTable: "table2"},
		&Shard{Type: ShardLinked, DB: "db1", Table: "table2", ParentTable: "table2"},
	}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardHash(t *testing.T) {
	nf := defaultNamespace()
	// locations count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardHash, Locations: []int{1}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardMod(t *testing.T) {
	nf := defaultNamespace()
	// locations count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardMod, Locations: []int{1}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardRange(t *testing.T) {
	nf := defaultNamespace()
	// locations count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardRange, Locations: []int{1}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, shardRule: %s", JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardDay(t *testing.T) {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// dateRange count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"1"}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// dateRange not match
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"1"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"abc"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"22222222"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"201910301-2019103"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"22222222-33333333"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"11111111-22222222"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// date range overlapped
	nf.Slices = []*Slice{&Slice{Name: "slice1"}, &Slice{Name: "slice2"}}
	nf.ShardRules = []*Shard{&Shard{Type: ShardDay, DateRange: []string{"20181030", "20181001"}, Slices: []string{"slice1", "slice2"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardMonth(t *testing.T) {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// dateRange count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"1"}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// dateRange not match
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"1"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"abc"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"222222"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"2019101-20191"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"222222-333333"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"111111-222222"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// date range overlapped
	nf.Slices = []*Slice{&Slice{Name: "slice1"}, &Slice{Name: "slice2"}}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMonth, DateRange: []string{"201810", "201809"}, Slices: []string{"slice1", "slice2"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardYear(t *testing.T) {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// dateRange count is not equal
	nf.ShardRules = []*Shard{&Shard{Type: ShardYear, DateRange: []string{"1"}, Slices: []string{}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// dateRange not match
	nf.ShardRules = []*Shard{&Shard{Type: ShardYear, DateRange: []string{"1"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardYear, DateRange: []string{"abc"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardYear, DateRange: []string{"20191-201"}, Slices: []string{"slice1"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// date range overlapped
	nf.Slices = []*Slice{&Slice{Name: "slice1"}, &Slice{Name: "slice2"}}
	nf.ShardRules = []*Shard{&Shard{Type: ShardYear, DateRange: []string{"2018", "2017"}, Slices: []string{"slice1", "slice2"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardMycatMod(t *testing.T) {
	if err := testVerifyShardRules_Error_ShardMycatMod(ShardMycatMod); err != nil {
		t.Error(err)
	}
}

func testVerifyShardRules_Error_ShardMycatMod(t string) error {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// location count is not equal of slice
	nf.ShardRules = []*Shard{&Shard{Type: t, Locations: []int{1}, Slices: []string{}, Databases: []string{""}}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// db bound value invalid
	nf.ShardRules = []*Shard{&Shard{Type: t, Locations: []int{1}, Slices: []string{"slice1"}, Databases: []string{"db[1-1]"}}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// location count is not equal of db
	nf.ShardRules = []*Shard{&Shard{Type: t, Locations: []int{1}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	return nil
}

func TestVerifyShardRules_Error_ShardMycatLong(t *testing.T) {
	if err := testVerifyShardRules_Error_ShardMycatMod(ShardMycatLong); err != nil {
		t.Error(err)
	}

	if err := testVerifyShardRules_Error_ShardMycatLong(ShardMycatLong); err != nil {
		t.Error(err)
	}
}

func testVerifyShardRules_Error_ShardMycatLong(t string) error {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// patitionCount to int array
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatLong, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "test", PartitionLength: ""}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// patitionLength to int array
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatLong, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "4", PartitionLength: "test"}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// partitionScope not match
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatLong, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "2", PartitionLength: "256"}}
	if err := nf.verifyShardRules(); err == nil {
		return fmt.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	return nil
}

func TestVerifyShardRules_Error_ShardMycatString(t *testing.T) {
	if err := testVerifyShardRules_Error_ShardMycatMod(ShardMycatString); err != nil {
		t.Error(err)
	}

	if err := testVerifyShardRules_Error_ShardMycatLong(ShardMycatString); err != nil {
		t.Error(err)
	}

	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// verify hashSlice
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatString, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "2", PartitionLength: "512", HashSlice: "test"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatString, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "2", PartitionLength: "512", HashSlice: "test:"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatString, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "2", PartitionLength: "512", HashSlice: ":test"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatString, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PartitionCount: "2", PartitionLength: "512", HashSlice: "a:b:c"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardMycatMURMUR(t *testing.T) {
	if err := testVerifyShardRules_Error_ShardMycatMod(ShardMycatMURMUR); err != nil {
		t.Error(err)
	}

	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// verify seed and virtualBucketTimes
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatMURMUR, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, Seed: "test", VirtualBucketTimes: ""}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatMURMUR, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, Seed: "5", VirtualBucketTimes: "test"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardMycatPaddingMod(t *testing.T) {
	if err := testVerifyShardRules_Error_ShardMycatMod(ShardMycatPaddingMod); err != nil {
		t.Error(err)
	}

	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "test", PadLength: "", ModBegin: "", ModEnd: ""}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "test", ModBegin: "", ModEnd: ""}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "1", ModBegin: "test", ModEnd: ""}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "1", ModBegin: "1", ModEnd: "test"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "2", PadLength: "1", ModBegin: "1", ModEnd: "1"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{1}, Slices: []string{"slice1"}, Databases: []string{"db0"}, PadFrom: "1", PadLength: "1", ModBegin: "1", ModEnd: "1"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "1", ModBegin: "-1", ModEnd: "4"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "1", ModBegin: "4", ModEnd: "4"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "0", ModBegin: "0", ModEnd: "4"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "-1", ModBegin: "0", ModEnd: "4"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	nf.ShardRules = []*Shard{&Shard{Type: ShardMycatPaddingMod, Locations: []int{2}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}, PadFrom: "1", PadLength: "2", ModBegin: "0", ModEnd: "4"}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestVerifyShardRules_Error_ShardGlobal(t *testing.T) {
	nf := defaultNamespace()
	nf.Slices = []*Slice{&Slice{Name: "slice1"}}
	// location count is not equal of slice
	nf.ShardRules = []*Shard{&Shard{Type: ShardGlobal, Locations: []int{1}, Slices: []string{}, Databases: []string{""}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// db bound value invalid
	nf.ShardRules = []*Shard{&Shard{Type: ShardGlobal, Locations: []int{1}, Slices: []string{"slice1"}, Databases: []string{"db[1-1]"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
	// location count is not equal of db
	nf.ShardRules = []*Shard{&Shard{Type: ShardGlobal, Locations: []int{1}, Slices: []string{"slice1"}, Databases: []string{"db[0-1]"}}}
	if err := nf.verifyShardRules(); err == nil {
		t.Errorf("test verifyShardRules should fail but pass, slices: %s, shardRule: %s", JSONEncode(nf.Slices), JSONEncode(nf.ShardRules))
	}
}

func TestNamespace_Verify(t *testing.T) {
	nsStr := `
{
    "name": "gaea_namespace_1",
    "online": true,
    "read_only": true,
    "allowed_dbs": {
        "db_ks": true,
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_ks": "db_ks",
        "db_mycat": "db_mycat_0"
    },
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3306",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3307",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        }
    ],
    "shard_rules": [
        {
            "db": "db_ks",
            "table": "tbl_ks",
            "type": "mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_child",
            "type": "linked",
            "key": "id",
            "parent_table": "tbl_ks"
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_user_child",
            "type": "linked",
            "key": "user_id",
            "parent_table": "tbl_ks"
        },
		{
            "db": "db_ks",
            "table": "tbl_ks_global_one",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
		{
            "db": "db_ks",
            "table": "tbl_ks_global_two",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
		{
			"db": "db_ks",
            "table": "tbl_ks_range",
            "type": "range",
			"key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"table_row_limit": 100
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_year",
            "type": "date_year",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"2014-2017",
				"2018-2019"
			]
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_month",
            "type": "date_month",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"201405-201406",
				"201408-201409"
			]
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_day",
            "type": "date_day",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"20140901-20140905",
				"20140907-20140908"
			]
		},
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_child",
            "type": "linked",
            "parent_table": "tbl_mycat",
            "key": "id"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_user_child",
            "type": "linked",
            "parent_table": "tbl_mycat",
            "key": "user_id"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_murmur",
            "type": "mycat_murmur",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_0","db_mycat_1","db_mycat_2","db_mycat_3"
            ],
			"seed": "0",
			"virtual_bucket_times": "160"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_long",
            "type": "mycat_long",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
			"partition_count": "4",
			"partition_length": "256"
        },
		{
            "db": "db_mycat",
            "table": "tbl_mycat_global_one",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
		{
            "db": "db_mycat",
            "table": "tbl_mycat_global_two",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_string",
            "type": "mycat_string",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
			"partition_count": "4",
			"partition_length": "256",
			"hash_slice": "20"
        }
    ],
	"global_sequences": [
		{
			"db": "db_mycat",
			"table": "tbl_mycat",
			"type": "test",
			"pk_name": "id"
		},
		{
			"db": "db_ks",
			"table": "tbl_ks",
			"type": "test",
			"pk_name": "user_id"
		}
	],
    "users": [
        {
            "user_name": "test_shard_hash",
            "password": "test_shard_hash",
            "namespace": "gaea_namespace_1",
            "rw_flag": 2,
            "rw_split": 1
        }
    ],
    "default_slice": "slice-0"
}`
	ns := &Namespace{}
	if err := json.Unmarshal([]byte(nsStr), ns); err != nil {
		t.Errorf("namespace unmarshal failed, err: %v", err)
	}

	if err := ns.Verify(); err != nil {
		t.Errorf("namespace verify failed, err: %v", err)
	}
}
