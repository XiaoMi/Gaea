// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

package etcdclient

import (
	"context"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/mock"
)

type MockKeysAPI struct {
	mock.Mock
}

func (m *MockKeysAPI) Get(ctx context.Context, key string, opts *client.GetOptions) (*client.Response, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) Set(ctx context.Context, key, value string, opts *client.SetOptions) (*client.Response, error) {
	args := m.Called(ctx, key, value, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) Delete(ctx context.Context, key string, opts *client.DeleteOptions) (*client.Response, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) Create(ctx context.Context, key, value string) (*client.Response, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) CreateInOrder(ctx context.Context, dir, value string, opts *client.CreateInOrderOptions) (*client.Response, error) {
	args := m.Called(ctx, dir, value, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) Update(ctx context.Context, key, value string) (*client.Response, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *MockKeysAPI) Watcher(key string, opts *client.WatcherOptions) client.Watcher {
	args := m.Called(key, opts)
	return args.Get(0).(client.Watcher)
}
