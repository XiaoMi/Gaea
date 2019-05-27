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

func newTestStore() *Store {
	c := NewClient(ConfigEtcd, "127.0.0.1:2381", "test", "test", "")
	return NewStore(c)
}

func TestNewStore(t *testing.T) {
	c := NewClient(ConfigEtcd, "127.0.0.1:2381", "test", "test", "")
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
