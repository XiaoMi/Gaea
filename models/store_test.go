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

func TestNewStore(t *testing.T) {
	// 创建 Mock 客户端
	mockClient := new(MockClient)

	// 设置 Mock 客户端的行为
	mockClient.On("BasePrefix").Return("/gaea_default_cluster")

	// 使用 Mock 客户端创建 Store
	store := NewStore(mockClient)

	// 检查 Store 的 prefix 是否正确
	expectedPrefix := "/gaea_default_cluster"
	actualPrefix := store.prefix
	if actualPrefix != expectedPrefix {
		t.Fatalf("Expected prefix '%s', got '%s'", expectedPrefix, actualPrefix)
	}

	// 确保所有预期的调用都被执行
	mockClient.AssertExpectations(t)
}

func TestNamespaceBase(t *testing.T) {
	// 创建 Mock 客户端
	mockClient := new(MockClient)

	// 设置 Mock 客户端的行为
	mockClient.On("BasePrefix").Return("/gaea_default_cluster")

	// 使用 Mock 客户端创建 Store
	store := NewStore(mockClient)

	// 测试 NamespaceBase 方法
	expectedBase := "/gaea_default_cluster/namespace"
	actualBase := store.NamespaceBase()
	if actualBase != expectedBase {
		t.Fatalf("Expected base '%s', got '%s'", expectedBase, actualBase)
	}

	// 确保所有预期的调用都被执行
	mockClient.AssertExpectations(t)
}

func TestNamespacePath(t *testing.T) {
	// 创建 Mock 客户端
	mockClient := new(MockClient)

	// 设置 Mock 客户端的行为
	mockClient.On("BasePrefix").Return("/gaea_default_cluster")

	// 使用 Mock 客户端创建 Store
	store := NewStore(mockClient)

	// 测试 NamespacePath 方法
	expectedPath := "/gaea_default_cluster/namespace/test_namespace"
	actualPath := store.NamespacePath("test_namespace")
	if actualPath != expectedPath {
		t.Fatalf("Expected path '%s', got '%s'", expectedPath, actualPath)
	}

	// 确保所有预期的调用都被执行
	mockClient.AssertExpectations(t)
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
