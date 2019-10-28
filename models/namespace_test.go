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
		DefaultPhyDBS:    nil,
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
