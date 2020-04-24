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

package proxy

import (
	"encoding/json"

	"github.com/XiaoMi/Gaea/util/requests"
)

// APIClient api client
type APIClient struct {
	addr     string
	user     string
	password string
}

// NewAPIClient create api client
func NewAPIClient(addr, user, password string) *APIClient {
	return &APIClient{addr: addr, user: user, password: password}
}

// PrepareConfig send prepare config
func (c *APIClient) PrepareConfig(name string) error {
	url := c.encodeURL("/api/proxy/config/prepare/%s", name)
	return requests.SendPut(url, c.user, c.password)
}

// CommitConfig send commit config
func (c *APIClient) CommitConfig(name string) error {
	url := c.encodeURL("/api/proxy/config/commit/%s", name)
	return requests.SendPut(url, c.user, c.password)
}

// DelNamespace send delete namespace to proxy
func (c *APIClient) DelNamespace(name string) error {
	url := c.encodeURL("/api/proxy/namespace/delete/%s", name)
	return requests.SendPut(url, c.user, c.password)
}

// GetNamespaceSQLFingerprint return sql fingerprint of specific namespace
func (c *APIClient) GetNamespaceSQLFingerprint(name string) (*SQLFingerprint, error) {
	var reply SQLFingerprint
	url := c.encodeURL("/api/proxy/stats/sessionsqlfingerprint/%s", name)
	resp, err := requests.SendGet(url, c.user, c.password)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		json.Unmarshal(resp.Body, &reply)
	}
	return &reply, err
}

func (c *APIClient) proxyConfigFingerprint() (string, error) {
	r := ""
	url := c.encodeURL("/api/proxy/config/fingerprint")
	resp, err := requests.SendGet(url, c.user, c.password)
	if err != nil {
		return r, err
	}
	if resp != nil && resp.Body != nil {
		json.Unmarshal(resp.Body, &r)
	}
	return r, err
}

// Ping ping proxy
func (c *APIClient) Ping() error {
	url := c.encodeURL("/api/proxy/ping")
	_, err := requests.SendGet(url, c.user, c.password)
	if err != nil {
		return err
	}
	return nil
}

func (c *APIClient) encodeURL(format string, args ...interface{}) string {
	return requests.EncodeURL(c.addr, format, args...)
}
