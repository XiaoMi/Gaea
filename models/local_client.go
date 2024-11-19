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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/log"
)

// NewLocalClient creates and initializes a new LocalClient instance.
// It takes namespaceStoragePath as the path for storing namespace-related data, and prefix as the prefix for identifying data.
// If namespaceStoragePath is empty, it returns an error.
// This function is responsible for converting the storage path to an absolute path and creating the storage directory if it does not exist.
func NewLocalClient(namespaceStoragePath string, prefix string) (*LocalClient, error) {
	if len(namespaceStoragePath) == 0 {
		return nil, fmt.Errorf("namespaceStoragePath must be a valid path")
	}
	// Convert namespaceStoragePath to an absolute path
	absStoragePath, err := filepath.Abs(namespaceStoragePath)
	if err != nil {
		log.Warn("Failed to get absolute path: %s: %v", namespaceStoragePath, err)
		return nil, fmt.Errorf("failed to get absolute path: %s: %w", namespaceStoragePath, err)
	}

	// Try to create the storage directory, ignore errors
	os.MkdirAll(absStoragePath, 0755)

	// Check if the path exists and is a directory
	info, err := os.Stat(absStoragePath)
	if err != nil {
		log.Warn("Failed to stat path: %s: %v", absStoragePath, err)
		return nil, fmt.Errorf("failed to stat path: %s: %w", absStoragePath, err)
	}
	if !info.IsDir() {
		log.Warn("Path exists but is not a directory: %s", absStoragePath)
		return nil, fmt.Errorf("path exists but is not a directory: %s", absStoragePath)
	}

	return &LocalClient{
		storagePath: absStoragePath,
		Prefix:      prefix,
		FileSuffix:  ".json",
	}, nil
}

// LocalClient implements the Client interface for local storage using the file system.
type LocalClient struct {
	storagePath string
	Prefix      string
	FileSuffix  string
}

// Create is a method of the LocalClient type, used to create a file at the specified path and write data into it.
// If no storage path is configured, the method will not perform any operations.
// The method first converts the specified path into a path relative to the storage path, then attempts to create the necessary directories and files.
func (lc *LocalClient) Create(path string, data []byte) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}

	// Get the full path of the file.
	fullPath, err := lc.FullNamespacePath(path)
	if err != nil {
		return err
	}
	// Create a storage directory
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Warn("Failed to create directories for %s: %v", dir, err)
		return fmt.Errorf("failed to create directories for %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		log.Warn("Failed to write data to file %s: %v", fullPath, err)
		return fmt.Errorf("failed to write data to file %s: %w", fullPath, err)
	}
	return nil
}

// Update updates the file at the specified path with the given data.
// This function first checks if the storage path is configured, and if not, it does nothing.
// It then constructs the full path relative to the storage path and attempts to create any necessary directories.
// Finally, it writes the provided data to the file at the computed full path.
func (lc *LocalClient) Update(path string, data []byte) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}

	// Get the full path of the file.
	fullPath, err := lc.FullNamespacePath(path)
	if err != nil {
		return err
	}

	// Create a storage directory
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Warn("Failed to create directories for %s: %v", dir, err)
		return fmt.Errorf("failed to create directories for %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		log.Warn("Failed to write data to file %s: %v", fullPath, err)
		return fmt.Errorf("failed to write data to file %s: %w", fullPath, err)
	}
	return nil
}

// UpdateWithTTL updates the file at the specified path with the given data and sets a time-to-live (TTL) for automatic deletion.
// This function is useful for temporarily storing data that needs to be deleted after a certain period of time to prevent data accumulation.
func (lc *LocalClient) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	// If no storage path is configured, do nothing.
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
		// Attempt to delete the file after the TTL has expired.
		// If deletion fails, log a warning.
		// Note: Logging failures here because the function has already returned, and there is no way to handle errors further.
		if err := lc.Delete(path); err != nil {
			log.Warn("Failed to delete file %s after TTL expired: %v", path, err)
		}
	}()

	// Return nil if the update was successful, indicating that the file will be deleted after the TTL expires.
	return nil
}

// Delete is a method of the LocalClient type, used to delete a file at a specified path.
// If no storage path is configured, the operation ends immediately without performing any action.
// If the path parameter is invalid or the file does not exist, the method handles these cases gracefully.
func (lc *LocalClient) Delete(path string) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, do nothing.
		return nil
	}

	// Get the full path of the file.
	fullPath, err := lc.FullNamespacePath(path)
	if err != nil {
		return err
	}

	// Check if the file exists before attempting to delete.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// If the file does not exist, return nil to keep the behavior consistent with ETCD.
		return nil
	}

	if err := os.Remove(fullPath); err != nil {
		log.Warn("Failed to delete file %s: %v", fullPath, err)
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}
	return nil
}

// Read reads data from a file at a specified path.
// This function first checks if the storage path is configured, then constructs the full path based on the provided path,
// and attempts to read the file. If the file does not exist or an error occurs, it returns the corresponding error.
func (lc *LocalClient) Read(path string) ([]byte, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}

	// Get the full path of the file.
	fullPath, err := lc.FullNamespacePath(path)
	if err != nil {
		return nil, err
	}

	// Check if the file exists before attempting to read.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// If the file does not exist, return nil, consistent with ETCD.
		return nil, nil
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		log.Warn("Failed to read file %s: %v", fullPath, err)
		return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}
	return data, nil
}

// List lists the files or directories under a specified path.
// This function first checks if the storagePath is configured; if not, it returns nil.
func (lc *LocalClient) List(path string) ([]string, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}

	// Get the full path of the directory
	fullPath, err := lc.FullDirPath(path)
	if err != nil {
		return nil, err
	}
	// Check if the directory exists before attempting to list.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// If the directory does not exist, return nil to be consistent with EtcdClientV3.
		return nil, nil
	}
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		log.Warn("Failed to read directory %s: %v", fullPath, err)
		return nil, fmt.Errorf("failed to read directory %s: %w", fullPath, err)
	}

	var files []string
	for _, entry := range entries {
		name := entry.Name()
		// Remove the `.json` suffix
		name = strings.TrimSuffix(name, lc.FileSuffix)
		files = append(files, name)

	}
	return files, nil
}

// ListWithValues lists all key-value pairs under the given path.
func (lc *LocalClient) ListWithValues(path string) (map[string]string, error) {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil, nil
	}

	// Get the full path of the directory
	fullPath, err := lc.FullDirPath(path)
	if err != nil {
		return nil, err
	}
	// Check if the directory exists before attempting to list.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// If the directory does not exist, return nil to be consistent with EtcdClientV3.
		return nil, nil
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		log.Warn("Failed to read directory %s: %v", fullPath, err)
		return nil, err
	}
	values := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			// Remove the `.json` suffix
			nameWithoutExt := strings.TrimSuffix(name, lc.FileSuffix)
			localPath := filepath.Join(fullPath, name)
			data, err := os.ReadFile(localPath)
			if err != nil {
				log.Warn("Failed to read file %s: %v", localPath, err)
				return nil, err
			}
			key := filepath.Join(path, nameWithoutExt)
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

// Clean removes all files from the specified directory within the storage path.
func (lc *LocalClient) Clean(path string) error {
	if len(lc.storagePath) == 0 {
		// If no storage path is configured, return nil.
		return nil
	}

	relPath, err := lc.safeJoinPath(path)
	if err != nil {
		log.Warn("Failed to join path: %s: %v", path, err)
		return err
	}
	dir := filepath.Join(lc.storagePath, relPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// If the directory does not exist, return without doing anything.
		return nil
	}

	// Read all entries in the directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Warn("Failed to read directory %s: %v", dir, err)
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var failedFiles []string

	for _, entry := range entries {
		if !entry.IsDir() {
			// Construct the full path for each file.
			localPath := filepath.Join(dir, entry.Name())
			// Attempt to remove the file.
			if err := os.Remove(localPath); err != nil {
				log.Warn("Failed to remove file %s: %v", localPath, err)
				failedFiles = append(failedFiles, localPath)
			}
		}
	}

	if len(failedFiles) > 0 {
		errMsg := fmt.Sprintf("failed to remove files: %v", failedFiles)
		log.Warn(errMsg)
		return fmt.Errorf("failed to clean: %s", errMsg)
	}

	return nil
}

// safeJoinPath ensures the provided path is safe and returns a cleaned version of it.
// If the path is an absolute path, it converts it to a relative path based on the root directory.
// It normalizes the path to eliminate elements like `..`.
// It prevents directory traversal attacks by ensuring the path does not start with `..` or contain `../`.
// It checks for forbidden characters in the path.
// It limits the path length to prevent resource exhaustion due to excessively long paths.
func (lc *LocalClient) safeJoinPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path cannot be empty")
	}

	// If the path is an absolute path, convert it to a relative path based on the root directory.
	if filepath.IsAbs(path) {
		var err error
		path, err = filepath.Rel("/", path)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
	}

	// Normalize the path to eliminate elements like `..`.
	cleanPath := filepath.Clean(path)

	// Prevent directory traversal attacks by ensuring the path does not start with `..` or contain `../`.
	if cleanPath == ".." || strings.HasPrefix(cleanPath, "../") || strings.Contains(cleanPath, "/../") {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	// Check for forbidden characters in the path.
	if strings.ContainsAny(cleanPath, "<>\"|?*") {
		return "", fmt.Errorf("invalid characters in path: %s", path)
	}

	// Limit the path length to prevent resource exhaustion due to excessively long paths.
	if len(cleanPath) > 1024 {
		return "", fmt.Errorf("path too long: %s", path)
	}

	// Return the cleaned and validated path.
	return cleanPath, nil
}

// FullNamespacePath returns the full path of a file or directory within the namespace.
// it joins the relative path with the storagePath of the LocalClient instance and appends the file suffix to construct the full path.
func (lc *LocalClient) FullNamespacePath(path string) (string, error) {
	// Convert path to a relative path relative to storagePath
	relPath, err := lc.safeJoinPath(path)
	if err != nil {
		log.Warn("Failed to join path: %s: %v", path, err)
		return "", fmt.Errorf("failed to join path: %s: %w", path, err)
	}
	fullPath := filepath.Join(lc.storagePath, relPath) + lc.FileSuffix
	return fullPath, nil
}

// FullDirPath returns the full directory path for a given relative path.
// This function ensures that the provided path is safely appended to the storagePath,
func (lc *LocalClient) FullDirPath(path string) (string, error) {
	// Convert path to a relative path relative to storagePath
	relPath, err := lc.safeJoinPath(path)
	if err != nil {
		// Log and return an error if the path combination fails
		log.Warn("Failed to join path: %s: %v", path, err)
		return "", fmt.Errorf("failed to join path: %s: %w", path, err)
	}
	// Safely combine storagePath and relPath to get the full path
	fullPath := filepath.Join(lc.storagePath, relPath)
	// Return the full path
	return fullPath, nil
}
