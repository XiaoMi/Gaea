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
package server

import (
	"testing"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
)

type userinfo struct {
	username string
	password string
}

type usercase struct {
	username  string
	password  string
	namespace string
}

func TestCreateUserManagerFromNamespaceConfigs(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	_, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserManager_CheckUser(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		username string
		valid    bool
	}{
		{username: "user1", valid: true},
		{username: "user2", valid: true},
		{username: "user3", valid: false},
		{username: "", valid: false},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualValid := userManager.CheckUser(test.username)
			if actualValid != test.valid {
				t.Errorf("CheckUser error, username: %s, expect: %t, actual: %t", test.username, test.valid, actualValid)
			}
		})
	}
}

func TestUserManager_GetNamespaceByUser(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: "namespace1"},
		{username: "user1", password: "pwd2", namespace: "namespace1"},
		{username: "user2", password: "pwd1", namespace: "namespace1"},
		{username: "user2", password: "pwd2", namespace: "namespace2"},
		{username: "user2", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd4", namespace: ""},
		{username: "user2", password: "pwd4", namespace: ""},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_ClearNamespaceUsers_Namespace0(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}
	userManager.ClearNamespaceUsers("namespace0")

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: "namespace1"},
		{username: "user1", password: "pwd2", namespace: "namespace1"},
		{username: "user2", password: "pwd1", namespace: "namespace1"},
		{username: "user2", password: "pwd2", namespace: "namespace2"},
		{username: "user2", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd4", namespace: ""},
		{username: "user2", password: "pwd4", namespace: ""},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_ClearNamespaceUsers_Namespace1(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}
	userManager.ClearNamespaceUsers("namespace1")

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: ""},
		{username: "user1", password: "pwd2", namespace: ""},
		{username: "user2", password: "pwd1", namespace: ""},
		{username: "user2", password: "pwd2", namespace: "namespace2"},
		{username: "user2", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd4", namespace: ""},
		{username: "user2", password: "pwd4", namespace: ""},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_ClearNamespaceUsers_Namespace2(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}
	userManager.ClearNamespaceUsers("namespace2")

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: "namespace1"},
		{username: "user1", password: "pwd2", namespace: "namespace1"},
		{username: "user2", password: "pwd1", namespace: "namespace1"},
		{username: "user2", password: "pwd2", namespace: ""},
		{username: "user2", password: "pwd3", namespace: ""},
		{username: "user1", password: "pwd3", namespace: ""},
		{username: "user1", password: "pwd4", namespace: ""},
		{username: "user2", password: "pwd4", namespace: ""},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_RebuildNamespaceUsers_Namespace1(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}

	ns := "namespace1"
	user1 := &userinfo{username: "user1", password: "pwd1"}
	user2 := &userinfo{username: "user1", password: "pwd4"}
	user3 := &userinfo{username: "user3", password: "pwd1"}

	newNamespace := createNamespaceUsers(ns, []*userinfo{user1, user2, user3})
	userManager.RebuildNamespaceUsers(newNamespace)

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: "namespace1"},
		{username: "user1", password: "pwd2", namespace: ""},
		{username: "user1", password: "pwd4", namespace: "namespace1"},
		{username: "user3", password: "pwd1", namespace: "namespace1"},
		{username: "user2", password: "pwd1", namespace: ""},
		{username: "user2", password: "pwd2", namespace: "namespace2"},
		{username: "user2", password: "pwd3", namespace: "namespace2"},
		{username: "user1", password: "pwd3", namespace: "namespace2"},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_RebuildNamespaceUsers_Namespace2(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}

	ns := "namespace2"
	user1 := &userinfo{username: "user5", password: "pwd1"}

	newNamespace := createNamespaceUsers(ns, []*userinfo{user1})
	userManager.RebuildNamespaceUsers(newNamespace)

	tests := []usercase{
		{username: "", password: "", namespace: ""},
		{username: "", password: "pwd", namespace: ""},
		{username: "user", password: "", namespace: ""},
		{username: "user", password: "pwd", namespace: ""},
		{username: "user1", password: "pwd1", namespace: "namespace1"},
		{username: "user1", password: "pwd2", namespace: "namespace1"},
		{username: "user2", password: "pwd1", namespace: "namespace1"},
		{username: "user2", password: "pwd2", namespace: ""},
		{username: "user2", password: "pwd3", namespace: ""},
		{username: "user1", password: "pwd3", namespace: ""},
		{username: "user5", password: "pwd1", namespace: "namespace2"},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			actualNamespace := userManager.GetNamespaceByUser(test.username, test.password)
			if actualNamespace != test.namespace {
				t.Errorf("GetNamespaceByUser error, username: %s, password: %s, expect: %s, actual: %s", test.username, test.password, test.namespace, actualNamespace)
			}
		})
	}
}

func TestUserManager_CheckPassword(t *testing.T) {
	nsCfg := prepareNamespaceUsers()
	userManager, err := CreateUserManager(nsCfg)
	if err != nil {
		t.Fatal(err)
	}
	salt := []byte("abcdefg_?!")

	tests := []struct {
		username string
		password string
		valid    bool
	}{
		{username: "user1", password: "pwd1", valid: true},
		{username: "user1", password: "pwd2", valid: true},
		{username: "user1", password: "pwd4", valid: false},
		{username: "user2", password: "pwd1", valid: true},
		{username: "user2", password: "pwd4", valid: false},
		{username: "user3", password: "pwd", valid: false},
	}
	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			auth := mysql.CalcPassword(salt, []byte(test.password))
			actualValid, actualPassword := userManager.CheckPassword(test.username, salt, auth)
			if actualValid == test.valid {
				if actualValid && (actualPassword != test.password) {
					t.Errorf("password not equal, expect: %v, acutal: %t, %s", test, actualValid, actualPassword)
				}
			} else {
				t.Errorf("valid not equal, expect: %v, acutal: %t, %s", test, actualValid, actualPassword)
			}
		})
	}
}

func prepareNamespaceUsers() map[string]*models.Namespace {
	nsMap := make(map[string]*models.Namespace)
	ns1 := "namespace1"
	ns1user1 := &userinfo{username: "user1", password: "pwd1"}
	ns1user2 := &userinfo{username: "user1", password: "pwd2"}
	ns1user3 := &userinfo{username: "user2", password: "pwd1"}
	namespace1 := createNamespaceUsers(ns1, []*userinfo{ns1user1, ns1user2, ns1user3})
	nsMap[ns1] = namespace1

	ns2 := "namespace2"
	ns2user1 := &userinfo{username: "user2", password: "pwd2"}
	ns2user2 := &userinfo{username: "user2", password: "pwd3"}
	ns2user3 := &userinfo{username: "user1", password: "pwd3"}
	namespace2 := createNamespaceUsers(ns2, []*userinfo{ns2user1, ns2user2, ns2user3})
	nsMap[ns2] = namespace2

	return nsMap
}

func createNamespaceUsers(ns string, users []*userinfo) *models.Namespace {
	var userList []*models.User
	for _, user := range users {
		u := &models.User{
			UserName: user.username,
			Password: user.password,
		}
		userList = append(userList, u)
	}
	return &models.Namespace{
		Name:  ns,
		Users: userList,
	}
}
