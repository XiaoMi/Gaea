// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdclient

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEtcdClient_Mkdir(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	dir := "/test/dir"

	// Mock the Set call for creating a directory
	mockAPI.On("Set", mock.Anything, dir, "", &client.SetOptions{Dir: true, PrevExist: client.PrevNoExist}).
		Return(&client.Response{}, nil).Once()

	err := etcdClient.Mkdir(dir)
	assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_Create(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/key"
	data := []byte("value")

	// Mock the Set call for creating a key-value pair
	mockAPI.On("Set", mock.Anything, path, "value", &client.SetOptions{PrevExist: client.PrevNoExist}).
		Return(&client.Response{}, nil).Once()

	err := etcdClient.Create(path, data)
	assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_Update(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/key"
	data := []byte("updated_value")

	// Mock the Set call for updating a key-value pair
	mockAPI.On("Set", mock.Anything, path, "updated_value", &client.SetOptions{PrevExist: client.PrevIgnore}).
		Return(&client.Response{}, nil).Once()

	err := etcdClient.Update(path, data)
	assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_UpdateWithTTL(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/key"
	data := []byte("value_with_ttl")
	ttl := time.Second * 10

	// Mock the Set call for updating a key-value pair with TTL
	mockAPI.On("Set", mock.Anything, path, "value_with_ttl", &client.SetOptions{PrevExist: client.PrevIgnore, TTL: ttl}).
		Return(&client.Response{}, nil).Once()

	err := etcdClient.UpdateWithTTL(path, data, ttl)
	assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_Delete(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/key"

	// Mock the Delete call
	mockAPI.On("Delete", mock.Anything, path, (*client.DeleteOptions)(nil)).
		Return(&client.Response{}, nil).Once()

	err := etcdClient.Delete(path)
	assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_Read(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/key"
	expectedValue := "value"

	// Create the specific context type
	ctx := context.Background()

	// Mock the Get call
	mockAPI.On("Get", ctx, path, (*client.GetOptions)(nil)).
		Return(&client.Response{
			Node: &client.Node{Value: expectedValue},
		}, nil).Once()

	data, err := etcdClient.Read(path)
	assert.NoError(t, err)
	assert.Equal(t, []byte(expectedValue), data)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_List(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/dir"
	childNodes := []*client.Node{
		{Key: "/test/dir/key1"},
		{Key: "/test/dir/key2"},
	}

	// 使用 mock.AnythingOfType 允许参数类型匹配
	mockAPI.On("Get", mock.Anything, path, (*client.GetOptions)(nil)).
		Return(&client.Response{
			Node: &client.Node{Dir: true, Nodes: childNodes},
		}, nil).Once()

	files, err := etcdClient.List(path)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"/test/dir/key1", "/test/dir/key2"}, files)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_ListWithValues(t *testing.T) {
	mockAPI := new(MockKeysAPI)
	etcdClient := &EtcdClient{kapi: mockAPI}

	path := "/test/dir"
	childNodes := []*client.Node{
		{Key: "/test/dir/key1", Value: "value1"},
		{Key: "/test/dir/key2", Value: "value2"},
	}

	// Mock the Get call
	mockAPI.On("Get", mock.Anything, path, &client.GetOptions{Recursive: true}).
		Return(&client.Response{
			Node: &client.Node{Dir: true, Nodes: childNodes},
		}, nil).Once()

	files, err := etcdClient.ListWithValues(path)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"/test/dir/key1": "value1",
		"/test/dir/key2": "value2",
	}, files)

	mockAPI.AssertExpectations(t)
}

func TestEtcdClient_Close(t *testing.T) {
	etcdClient := &EtcdClient{}

	// Close the client
	err := etcdClient.Close()
	assert.NoError(t, err)

	// Attempt to close again
	err = etcdClient.Close()
	assert.NoError(t, err)
}

func Test_isErrNoNode(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeKeyNotFound
	if !IsErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if IsErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
}

func Test_isErrNodeExists(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeNodeExist
	if !isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
}
