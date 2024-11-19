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
)

func TestNewLocalClient(t *testing.T) {
	// Test absolute path
	absPath := "/tmp/gaea_local_storage_abs"
	lcAbs, err := NewLocalClient(absPath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient with absolute path: %v", err)
	}
	if lcAbs.storagePath != absPath {
		t.Fatalf("Expected storagePath %s, got %s", absPath, lcAbs.storagePath)
	}
	// Check if the directory is created
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Fatalf("Storage directory not created for absolute path")
	}

	// Test relative path
	relPath := "./gaea_local_storage_rel"
	absRelPath, _ := filepath.Abs(relPath)
	lcRel, err := NewLocalClient(relPath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient with relative path: %v", err)
	}
	if lcRel.storagePath != absRelPath {
		t.Fatalf("Expected storagePath %s, got %s", absRelPath, lcRel.storagePath)
	}
	// Check if the directory is created
	if _, err := os.Stat(absRelPath); os.IsNotExist(err) {
		t.Fatalf("Storage directory not created for relative path")
	}

	// Clean up the test directory
	os.RemoveAll(absPath)
	os.RemoveAll(absRelPath)
}

func TestLocalClient_Create(t *testing.T) {
	// Test file creation
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Test file creation
	path := "/gaea_default_cluster/namespace/test_namespace"
	data := []byte("test data")
	err = lc.Create(path, data)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check if the directory is created
	relPath, _ := lc.safeJoinPath(path)
	fullPath := filepath.Join(lc.storagePath, relPath) + lc.FileSuffix
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Fatalf("File not created at %s", fullPath)
	}
}

func TestLocalClient_Read(t *testing.T) {
	// Test file creation
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Test file creation
	path := "/gaea_default_cluster/namespace/test_namespace"
	data := []byte("test data")
	lc.Create(path, data)

	// Test read
	readData, err := lc.Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if string(readData) != string(data) {
		t.Fatalf("Read data mismatch: expected '%s', got '%s'", data, readData)
	}
}

func TestLocalClient_Update(t *testing.T) {
	// Initialize LocalClient
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Create the initial file
	path := "/gaea_default_cluster/namespace/test_namespace"
	data := []byte("initial data")
	lc.Create(path, data)

	// Update file content
	updatedData := []byte("updated data")
	err = lc.Update(path, updatedData)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Read the updated content
	readData, err := lc.Read(path)
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if string(readData) != string(updatedData) {
		t.Fatalf("Updated data mismatch: expected '%s', got '%s'", updatedData, readData)
	}
}

func TestLocalClient_Delete(t *testing.T) {
	// Initialize LocalClient
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Create a file for deletion
	path := "/gaea_default_cluster/namespace/test_namespace"
	data := []byte("test data")
	lc.Create(path, data)

	// Delete the file
	err = lc.Delete(path)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Delete the file
	readData, err := lc.Read(path)
	if err != nil {
		t.Fatalf("Read after delete failed: %v", err)
	}
	if readData != nil {
		t.Fatalf("Data should be nil after delete, got: '%s'", readData)
	}
}

func TestLocalClient_List(t *testing.T) {
	// Initialize LocalClient
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Create multiple files
	lc.Create("/gaea_default_cluster/namespace/test_namespace1", []byte("data1"))
	lc.Create("/gaea_default_cluster/namespace/test_namespace2", []byte("data2"))
	lc.Create("/gaea_default_cluster/namespace/test_namespace3", []byte("data3"))

	// List files
	items, err := lc.List("/gaea_default_cluster/namespace/")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	expectedItems := map[string]bool{
		"test_namespace1": true,
		"test_namespace2": true,
		"test_namespace3": true,
	}

	if len(items) != len(expectedItems) {
		t.Fatalf("Expected %d items, got %d", len(expectedItems), len(items))
	}

	for _, item := range items {
		if !expectedItems[item] {
			t.Fatalf("Unexpected item in list: %s", item)
		}
		delete(expectedItems, item)
	}

	if len(expectedItems) != 0 {
		t.Fatalf("Missing items in list: %v", expectedItems)
	}
}

func TestLocalClient_ListWithValues(t *testing.T) {
	// Initialize LocalClient
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Create multiple files
	lc.Create("/gaea_default_cluster/namespace/test_namespace1", []byte("data1"))
	lc.Create("/gaea_default_cluster/namespace/test_namespace2", []byte("data2"))
	lc.Create("/gaea_default_cluster/namespace/test_namespace3", []byte("data3"))

	// List files
	values, err := lc.ListWithValues("/gaea_default_cluster/namespace/")
	if err != nil {
		t.Fatalf("ListWithValues failed: %v", err)
	}

	expectedValues := map[string]string{
		"/gaea_default_cluster/namespace/test_namespace1": "data1",
		"/gaea_default_cluster/namespace/test_namespace2": "data2",
		"/gaea_default_cluster/namespace/test_namespace3": "data3",
	}

	if len(values) != len(expectedValues) {
		t.Fatalf("Expected %d items, got %d", len(expectedValues), len(values))
	}

	for key, value := range values {
		expectedValue, exists := expectedValues[key]
		if !exists {
			t.Fatalf("Unexpected key in values: %s", key)
		}
		if value != expectedValue {
			t.Fatalf("Value mismatch for key %s: expected '%s', got '%s'", key, expectedValue, value)
		}
		delete(expectedValues, key)
	}

	if len(expectedValues) != 0 {
		t.Fatalf("Missing keys in values: %v", expectedValues)
	}
}
func TestLocalClient_UpdateWithTTL(t *testing.T) {
	// Initialize LocalClient
	storagePath := "/tmp/gaea_local_storage_abs"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Update the file and set TTL
	path := "/gaea_default_cluster/namespace/test_namespace"
	data := []byte("temporary data")
	ttl := 1 * time.Second
	err = lc.UpdateWithTTL(path, data, ttl)
	if err != nil {
		t.Fatalf("UpdateWithTTL failed: %v", err)
	}

	//Confirm that the file exists
	readData, err := lc.Read(path)
	if err != nil {
		t.Fatalf("Read after UpdateWithTTL failed: %v", err)
	}
	if string(readData) != string(data) {
		t.Fatalf("Data mismatch: expected '%s', got '%s'", data, readData)
	}

	// Wait for TTL to expire
	time.Sleep(ttl + time.Second)

	// Confirm that the file has been deleted
	readData, err = lc.Read(path)
	if err != nil {
		t.Fatalf("Read after TTL expired failed: %v", err)
	}
	if readData != nil {
		t.Fatalf("Data should be nil after TTL expired, got: '%s'", readData)
	}
}

func TestLocalClient_Clean(t *testing.T) {
	// Initialize LocalClient
	storagePath := "./test_storage_clean"
	lc, err := NewLocalClient(storagePath, "/gaea")
	if err != nil {
		t.Fatalf("Failed to initialize LocalClient: %v", err)
	}
	defer os.RemoveAll(storagePath)

	// Create multiple files
	lc.Create("/namespace/item1", []byte("data1"))
	lc.Create("/namespace/item2", []byte("data2"))

	// Clean up the directory
	err = lc.Clean("/namespace")
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}

	// Check if the directory is empty
	items, err := lc.List("/namespace")
	if err != nil {
		t.Fatalf("List after clean failed: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("Items should be empty after clean, got: %v", items)
	}
}

func TestLocalClient_safeJoinPath(t *testing.T) {
	lc := &LocalClient{}

	// Test absolute path
	absPath := "/namespace/test"
	relPath, err := lc.safeJoinPath(absPath)
	if err != nil {
		t.Fatalf("safeJoinPath failed for absolute path: %v", err)
	}
	expectedRelPath := "namespace/test"
	if relPath != expectedRelPath {
		t.Fatalf("Expected relative path '%s', got '%s'", expectedRelPath, relPath)
	}

	// Test relative path
	relInputPath := "namespace/test"
	relPath, err = lc.safeJoinPath(relInputPath)
	if err != nil {
		t.Fatalf("safeJoinPath failed for relative path: %v", err)
	}
	if relPath != expectedRelPath {
		t.Fatalf("Expected relative path '%s', got '%s'", expectedRelPath, relPath)
	}

	// Test for illegal paths
	invalidPath := "../etc/passwd"
	_, err = lc.safeJoinPath(invalidPath)
	if err == nil {
		t.Fatalf("Expected error for invalid path '%s', got nil", invalidPath)
	}
}
