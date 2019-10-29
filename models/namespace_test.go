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

func TestFunc_VerifyName(t *testing.T) {
	n := defaultNamespace()
	if err := n.verifyName(); err != nil {
		t.Errorf("test verifyName failed, %v", err)
	}

	nf := defaultNamespace()
	nf.Name = ""
	if err := nf.verifyName(); err == nil {
		t.Errorf("test verifyName failed, should fail but pass, name: %v", nf.Name)
	}
}

func TestFunc_VerifyAllowDBS(t *testing.T) {
	n := defaultNamespace()
	n.AllowedDBS["db1"] = true
	if err := n.verifyAllowDBS(); err != nil {
		t.Errorf("test verifyAllowDBS failed, %v", err)
	}

	nf := defaultNamespace()
	if err := nf.verifyAllowDBS(); err == nil {
		t.Errorf("test verifyAllowDBS failed, should fail but pass, name: %v", nf.Name)
	}
}

func TestFunc_VerifyUsers(t *testing.T) {
	n := defaultNamespace()
	u1 := &User{UserName: "u1", Namespace: n.Name, Password: "pw1", RWFlag: ReadOnly, RWSplit: NoReadWriteSplit, OtherProperty: 0}
	u2 := &User{UserName: "u2", Namespace: n.Name, Password: "pw2", RWFlag: ReadWrite, RWSplit: ReadWriteSplit, OtherProperty: StatisticUser}
	n.Users = append(n.Users, u1)
	n.Users = append(n.Users, u2)

	if err := n.verifyUsers(); err != nil {
		t.Errorf("test verifyUsers failed, %v", err)
	}

	nf := defaultNamespace()
	uf1 := &User{UserName: "u1", Namespace: "someone", Password: "pw1", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	uf2 := &User{UserName: "u1", Namespace: n.Name, Password: "pw2", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	uf3 := &User{UserName: "", Namespace: "", Password: "", RWFlag: -1, RWSplit: -1, OtherProperty: -1}
	nf.Users = append(nf.Users, uf1)
	nf.Users = append(nf.Users, uf2)
	nf.Users = append(nf.Users, uf3)

	if err := nf.verifyUsers(); err == nil {
		t.Errorf("test verifyUsers failed, should fail but pass, users: %s", JSONEncode(nf.Users))
	}
}

func TestFunc_VerifySlowSQLTime(t *testing.T) {
	n := defaultNamespace()
	ssts := []string{"", "10"}
	for _, sst := range ssts {
		n.SlowSQLTime = sst
		if err := n.verifySlowSQLTime(); err != nil {
			t.Errorf("test verifySlowSQLTime failed, %v", err)
		}
	}

	sstfs := []string{"-1", "10.0", "test"}
	for _, sst := range sstfs {
		n.SlowSQLTime = sst
		if err := n.verifySlowSQLTime(); err == nil {
			t.Errorf("test verifySlowSQLTime failed, should fail but pass, sst: %v", n.SlowSQLTime)
		}
	}
}

func TestFunc_VerifyDBs(t *testing.T) {
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

	nf := defaultNamespace()
	// logic database mode
	nf.AllowedDBS["test1"] = true
	nf.DefaultPhyDBS["test2"] = ""
	if err := nf.verifyDBs(); err == nil {
		t.Errorf("test verifyDBs should fail but pass, allowedDBS: %v, defaultPhyDBS: %v", nf.AllowedDBS, nf.DefaultPhyDBS)
	}
}

func TestFunc_VerifyAllowIps(t *testing.T) {
	n := defaultNamespace()
	n.AllowedIP = append(n.AllowedIP, "  ")
	n.AllowedIP = append(n.AllowedIP, "10.221.163.82")
	if err := n.verifyAllowIps(); err != nil {
		t.Errorf("test verifyAllowIps failed, %v", err)
	}

	nf := defaultNamespace()
	var ipfs = []string{"test", "1.1.1"}
	for _, ipf := range ipfs {
		nf.AllowedIP = []string{ipf}
		if err := nf.verifyAllowIps(); err == nil {
			t.Errorf("test verifyAllowIps should fail but pass, %v", nf.AllowedIP)
		}
	}
}

func TestFunc_VerifyCharset(t *testing.T) {
	n := defaultNamespace()
	var ccs = [][]string{[]string{"", ""}, []string{"big5", ""}, []string{"big5", "big5_chinese_ci"}}
	for _, cc := range ccs {
		n.DefaultCharset = cc[0]
		n.DefaultCollation = cc[1]
		if err := n.verifyCharset(); err != nil {
			t.Errorf("test verifyCharset failed, %v", err)
		}
	}

	var ccfs = [][]string{[]string{"", "test"}, []string{"test", ""}, []string{"big5", "test"}, []string{"big5", "latin2_czech_cs"}}
	for _, ccf := range ccfs {
		n.DefaultCharset = ccf[0]
		n.DefaultCollation = ccf[1]
		if err := n.verifyCharset(); err == nil {
			t.Errorf("test verifyCharset should fail but pass, charset: %s, collation: %s", n.DefaultCharset, n.DefaultCollation)
		}
	}
}
