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

package file

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/log"
)

const (
	defaultFilePath = "./etc/file"
)

// Client used to test with config from file
type Client struct {
	Prefix string
}

// New constructor of EtcdClient
func New(path string) (*Client, error) {
	if strings.TrimSpace(path) == "" {
		path = defaultFilePath
	}
	if err := checkDir(path); err != nil {
		log.Warn("check file config directory failed, %v", err)
		return nil, err
	}
	return &Client{Prefix: path}, nil
}

func checkDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("invalid path")
	}
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("invalid path, should be a directory")
	}

	return nil
}

// Close do nothing
func (c *Client) Close() error {
	return nil
}

// Create do nothing
func (c *Client) Create(path string, data []byte) error {
	return nil
}

// Update do nothing
func (c *Client) Update(path string, data []byte) error {
	return nil
}

// UpdateWithTTL update path with data and ttl
func (c *Client) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	return nil
}

// Delete delete path
func (c *Client) Delete(path string) error {
	return nil
}

// Read read file data
func (c *Client) Read(file string) ([]byte, error) {
	value, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// List list path, return slice of all files
func (c *Client) List(path string) ([]string, error) {
	r := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return r, err
	}

	for _, f := range files {
		r = append(r, f.Name())
	}

	return r, nil
}

// BasePrefix return base prefix
func (c *Client) BasePrefix() string {
	return c.Prefix
}
