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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/XiaoMi/Gaea/log"
)

// NewRemoteClientConfig creates a new RemoteClientConfig based on the given parameters.
// namespaceStoragePath can be an absolute path or a relative path relative to the current execution directory
func NewLocalClient(namespaceStoragePath string, prefix string) (*LocalClient, error) {
	if len(namespaceStoragePath) == 0 {
		return nil, fmt.Errorf("namespaceStoragePath must be a valid path")
	}

	return &LocalClient{
		storagePath: namespaceStoragePath,
		Prefix:      prefix,
	}, nil
}

// LocalClient implements the Client interface for local storage using the file system.
type LocalClient struct {
	storagePath string
	Prefix      string
}

// Create creates a file with the given data at the specified path.
func (lc *LocalClient) Create(path string, data []byte) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}
	fullPath := filepath.Join(lc.storagePath, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write namespace data to file: %w", err)
	}
	return nil
}

// Update updates a file with the given data at the specified path.
func (lc *LocalClient) Update(path string, data []byte) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}
	fullPath := filepath.Join(lc.storagePath, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		log.Warn("local update node %s failed: %s", path, err)
		return fmt.Errorf("failed to write namespace data to file: %w", err)
	}

	return nil
}

// UpdateWithTTL updates a file with the given data and schedules its deletion after TTL expires.
func (lc *LocalClient) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}
	// Update the file with the given data.
	if err := lc.Update(path, data); err != nil {
		return err
	}

	// Schedule deletion after TTL expires to maintain data consistency.
	go func() {
		time.Sleep(ttl)
		if err := lc.Delete(path); err != nil {
			log.Warn("Failed to delete file %s after TTL expired: %v", path, err)
		}
	}()

	return nil
}

// Delete deletes the file at the specified path.
func (lc *LocalClient) Delete(path string) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}
	localPath := filepath.Join(lc.storagePath, path)
	// Check if the file exists before attempting to delete.
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// If the file does not exist, return nil to keep the behavior consistent with ETCD.
		return nil
	}
	// File exists, proceed to delete.
	err := os.Remove(localPath)
	if err != nil {
		log.Debug("local delete node %s failed: %s", path, err)
		return err
	}
	return nil
}

// Read reads data from the file at the specified path.
func (lc *LocalClient) Read(path string) ([]byte, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}
	localPath := filepath.Join(lc.storagePath, path)
	// Check if the file exists before attempting to read.
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// If the file does not exist, return nil, consistent with ETCD.
		return nil, nil
	}
	return os.ReadFile(localPath)
}

// List lists all files under the given path.
func (lc *LocalClient) List(path string) ([]string, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}
	localDir := filepath.Join(lc.storagePath, path)
	// Check if the directory exists before attempting to list.
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		// If the directory does not exist, return nil to be consistent with EtcdClientV3.
		return nil, nil
	}
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

// ListWithValues lists all key-value pairs under the given path.
func (lc *LocalClient) ListWithValues(path string) (map[string]string, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}
	localDir := filepath.Join(lc.storagePath, path)
	// Check if the directory exists
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		// If the directory does not exist, return an empty map and a nil error ,Be consistent with ETCD behavior
		return nil, nil
	}
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			localPath := filepath.Join(localDir, entry.Name())
			data, err := os.ReadFile(localPath)
			if err != nil {
				return nil, err
			}
			key := filepath.Join(path, entry.Name())
			values[key] = string(data)
		}
	}
	return values, nil
}

// Close is a no-op for LocalClient.
func (lc *LocalClient) Close() error {
	// No resources to close for local storage.
	return nil
}

// BasePrefix returns the base prefix for the local client.
func (lc *LocalClient) BasePrefix() string {
	return lc.Prefix
}

// Clean deletes all files in the specified storage directory, but keeps the directory itself.
func (lc *LocalClient) Clean(path string) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil
	}
	dir := filepath.Join(lc.storagePath, path)
	// Check if the storage directory exists.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// If the directory does not exist, return without doing anything.
		return nil
	}

	// Read all entries in the directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to clean directory %s: %w", dir, err)
	}
	// Iterate over all entries in the directory.
	for _, entry := range entries {
		if !entry.IsDir() {
			// Construct the full path for each file.
			localPath := filepath.Join(dir, entry.Name())
			// Attempt to remove the file.
			if err := os.Remove(localPath); err != nil {
				log.Warn("failed to remove file %s, err: %v", localPath, err)
			}
		}
	}
	return nil
}
