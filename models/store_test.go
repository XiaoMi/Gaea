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
	"testing"

	"github.com/bytedance/mockey"
)

func newTestStore() *Store {
	c, _ := NewClient(ConfigEtcd, "127.0.0.1:2381", "test", "test", "")

	return NewStore(c)
}

func TestNewStore(t *testing.T) {
	c, _ := NewClient(ConfigEtcd, "127.0.0.1:2381", "test", "test", "")
	store := NewStore(c)
	defer store.Close()
	if store == nil {
		t.Fatalf("test NewStore failed")
	}
}

func TestNamespaceBase(t *testing.T) {
	store := newTestStore()
	defer store.Close()
	base := store.NamespaceBase()
	if base != "/gaea/namespace" {
		t.Fatalf("test NamespaceBase failed, %v", base)
	}
}

func TestNamespacePath(t *testing.T) {
	store := newTestStore()
	defer store.Close()
	path := store.NamespacePath("test")
	if path != "/gaea/namespace/test" {
		t.Fatalf("test NamespacePath failed, %v", path)
	}
}

func TestLoadNamespace(t *testing.T) {
	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*Store).LoadNamespace).Return(nil, nil).Build()
	})
}

func TestLoadNamespaces(t *testing.T) {
	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*Store).LoadNamespaces).Return(nil, nil).Build()
	})
}

func TestLoadOriginNamespace(t *testing.T) {
	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*Store).LoadOriginNamespace).Return(nil, nil).Build()
	})

}

func TestLoadOriginNamespaces(t *testing.T) {
	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*Store).LoadOriginNamespaces).Return(nil, nil).Build()
	})

}

func TestDecryptNamespaces(t *testing.T) {
	ns := `{
				"open_general_log": true,
				"is_encrypt": true,
				"name": "test_namespace",
				"online": true,
				"read_only": false,
				"allowed_dbs": {
					"db_e2e_test": true
				},
				"default_phy_dbs": null,
				"slow_sql_time": "1000",
				"black_sql": [],
				"allowed_ip": null,
				"slices": [
					{
						"name": "slice-0",
						"user_name": "bB8gaOiGowx+mJHV/rT21w==",
						"password": "bB8gaOiGowx+mJHV/rT21w==",
						"master": "127.0.0.1:3319",
						"slaves": [],
						"statistic_slaves": [],
						"capacity": 12,
						"max_capacity": 24,
						"idle_timeout": 3600,
						"capability": 238087,
						"init_connect": "",
						"health_check_sql": ""
					}
				],
				"shard_rules": [],
				"users": [
					{
						"user_name": "bB8gaOiGowx+mJHV/rT21w==",
						"password": "bB8gaOiGowx+mJHV/rT21w==",
						"namespace": "test_namespace",
						"rw_flag": 2,
						"rw_split": 1,
						"other_property": 0
					}
				],
				"default_slice": "slice-0",
				"global_sequences": null,
				"default_charset": "",
				"default_collation": "",
				"max_sql_execute_time": 0,
				"max_sql_result_size": 0,
				"max_client_connections": 100000,
				"down_after_no_alive": 32,
				"seconds_behind_master": 32,
				"check_select_lock": false,
				"support_multi_query": false,
				"local_slave_read_priority": 0,
				"set_for_keep_session": false,
				"client_qps_limit": 0,
				"support_limit_transaction": false,
				"allowed_session_variables": {}
			}`

	// Convert JSON string to Namespace struct
	var namespace Namespace
	if err := json.Unmarshal([]byte(ns), &namespace); err != nil {
		t.Fatalf("failed to unmarshal namespace: %v", err)
	}

	// Create a map with namespace as expected by DecryptNamespaces function
	originNamespaces := map[string]*Namespace{
		"test_namespace": &namespace,
	}

	// DecryptNamespaces function call
	decryptedNamespaces, err := DecryptNamespaces(originNamespaces, "1234abcd5678efg*")
	if err != nil {
		t.Errorf("unexpected error in DecryptNamespaces: %v", err)
	}
	// Assert the output
	if decryptedNamespaces["test_namespace"].Users[0].UserName != "superroot" {
		t.Errorf("expected user name to be 'test_namespace', got '%s'", decryptedNamespaces["test_namespace"].Users[0].UserName)
	}
}

func TestDecryptNamespacesError(t *testing.T) {
	ns := `{
				"open_general_log": true,
				"is_encrypt": true,
				"name": "",
				"online": true,
				"read_only": false,
				"allowed_dbs": {
					"db_e2e_test": true
				},
				"default_phy_dbs": null,
				"slow_sql_time": "1000",
				"black_sql": [],
				"allowed_ip": null,
				"slices": [
					{
						"name": "slice-0",
						"user_name": "bB8gaOiGowx+mJHV/rT21w==",
						"password": "bB8gaOiGowx+mJHV/rT21w==",
						"master": "127.0.0.1:3319",
						"slaves": [],
						"statistic_slaves": [],
						"capacity": 12,
						"max_capacity": 24,
						"idle_timeout": 3600,
						"capability": 238087,
						"init_connect": "",
						"health_check_sql": ""
					}
				],
				"shard_rules": [],
				"users": [
					{
						"user_name": "bB8gaOiGowx+mJHV/rT21w==",
						"password": "bB8gaOiGowx+mJHV/rT21w==",
						"namespace": "test_namespace",
						"rw_flag": 2,
						"rw_split": 1,
						"other_property": 0
					}
				],
				"default_slice": "slice-0",
				"global_sequences": null,
				"default_charset": "",
				"default_collation": "",
				"max_sql_execute_time": 0,
				"max_sql_result_size": 0,
				"max_client_connections": 100000,
				"down_after_no_alive": 32,
				"seconds_behind_master": 32,
				"check_select_lock": false,
				"support_multi_query": false,
				"local_slave_read_priority": 0,
				"set_for_keep_session": false,
				"client_qps_limit": 0,
				"support_limit_transaction": false,
				"allowed_session_variables": {}
			}`

	// Convert JSON string to Namespace struct
	var namespace Namespace
	if err := json.Unmarshal([]byte(ns), &namespace); err != nil {
		t.Fatalf("failed to unmarshal namespace: %v", err)
	}

	// Create a map with namespace as expected by DecryptNamespaces function
	originNamespaces := map[string]*Namespace{
		"test_namespace": &namespace,
	}

	// DecryptNamespaces function call
	_, err := DecryptNamespaces(originNamespaces, "1234abcd5678efg*")
	if err == nil {
		t.Errorf("expected decrypt error : %v", err)
	}
}
