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

package models

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "test_storage_*")
	if err != nil {
		t.Fatal("Failed to create temporary directory for tests")
	}
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	return tempDir
}

func TestLocalClient_Create(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	data := []byte("test data")

	err := client.Create(path, data)
	assert.NoError(t, err)

	// Verify file is created with correct data
	fullPath := filepath.Join(client.storagePath, path)
	writtenData, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, data, writtenData)
}

func TestLocalClient_Update(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	data := []byte("updated data")

	err := client.Update(path, data)
	assert.NoError(t, err)

	// Verify file is updated with correct data
	fullPath := filepath.Join(client.storagePath, path)
	writtenData, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, data, writtenData)
}

func TestLocalClient_UpdateWithTTL(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	data := []byte("data with TTL")
	ttl := 100 * time.Millisecond

	err := client.UpdateWithTTL(path, data, ttl)
	assert.NoError(t, err)

	// Verify file is created with correct data
	fullPath := filepath.Join(client.storagePath, path)
	writtenData, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, data, writtenData)

	// Wait for TTL to expire and verify file is deleted
	time.Sleep(200 * time.Millisecond)
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalClient_Delete(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	data := []byte("to be deleted")

	err := client.Create(path, data)
	assert.NoError(t, err)

	// Delete file
	err = client.Delete(path)
	assert.NoError(t, err)

	// Verify file is deleted
	fullPath := filepath.Join(client.storagePath, path)
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalClient_Read(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	data := []byte("read this data")

	err := client.Create(path, data)
	assert.NoError(t, err)

	// Read file
	readData, err := client.Read(path)
	assert.NoError(t, err)
	assert.Equal(t, data, readData)
}

func TestLocalClient_List(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	// Create multiple files
	paths := []string{"gaea_default_cluster/namespace/test_namespace1", "gaea_default_cluster/namespace/test_namespace2", "gaea_default_cluster/namespace/test_namespace3"}
	for _, path := range paths {
		assert.NoError(t, client.Create(path, []byte("test data")))
	}

	// List files
	files, err := client.List("gaea_default_cluster/namespace/")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"test_namespace1", "test_namespace2", "test_namespace3"}, files)
}

func TestLocalClient_ListWithValues(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	// Create multiple files with values
	filesData := map[string]string{
		"/gaea_default_cluster/namespace/test_namespace1": "data1",
		"/gaea_default_cluster/namespace/test_namespace2": "data2",
		"/gaea_default_cluster/namespace/test_namespace3": "data3",
	}
	for path, data := range filesData {
		assert.NoError(t, client.Create(path, []byte(data)))
	}

	// List files with values
	values, err := client.ListWithValues("/gaea_default_cluster/namespace")
	assert.NoError(t, err)

	// Ensure all expected values are found
	for key, expectedValue := range filesData {
		actualValue, ok := values[key]
		if !ok {
			t.Errorf("Key %s not found in the values map", key)
		} else {
			assert.Equal(t, expectedValue, actualValue)
		}
	}
}

func TestLocalClient_UpdateAndDelete(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	path := "gaea_default_cluster/namespace/test_namespace"
	initialData := []byte("initial data")
	updateData := []byte("updated data")

	// Initialize file creation
	err := client.Create(path, initialData)
	assert.NoError(t, err)

	err = client.Update(path, updateData)
	assert.NoError(t, err)

	err = client.Delete(path)
	assert.NoError(t, err)
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalClient_Clean(t *testing.T) {
	var client = LocalClient{}
	client.storagePath = setupTestDir(t)

	// Create multiple files with values
	filesData := map[string]string{
		"/gaea_default_cluster/namespace/test_namespace1": "data1",
		"/gaea_default_cluster/namespace/test_namespace2": "data2",
		"/gaea_default_cluster/namespace/test_namespace3": "data3",
	}
	for path, data := range filesData {
		assert.NoError(t, client.Create(path, []byte(data)))
	}

	// clean storage
	err := client.Clean("/gaea_default_cluster/namespace")
	assert.NoError(t, err)

	for _, path := range filesData {
		fullPath := filepath.Join(client.storagePath, path)
		_, err = os.Stat(fullPath)
		assert.True(t, os.IsNotExist(err), "expected %s to be deleted", fullPath)
	}

}
