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
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
)

// Stats proxy stats
type Stats struct {
	Host     string `json:"host"`
	Closed   bool   `json:"closed"`
	Error    string `json:"error"`
	UnixTime int64  `json:"unixtime"`
	Timeout  bool   `json:"timeout"`
}

// SQLFingerprint sql指纹
type SQLFingerprint struct {
	SlowSQL  map[string]string `json:"slow_sql"`
	ErrorSQL map[string]string `json:"error_sql"`
}

// GetStats return proxy status
func GetStats(p *models.ProxyMonitorMetric, cfg *models.CCConfig, timeout time.Duration) *Stats {
	fmt.Println(string(p.Encode()))
	var ch = make(chan struct{})
	var host = p.IP + ":" + p.AdminPort
	fmt.Println(host)
	stats := &Stats{}

	go func(host string) {
		defer close(ch)
		stats.Host = host
		_, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
		if err != nil {
			stats.Error = err.Error()
			stats.Closed = true
		} else {
			stats.Closed = false
		}
	}(host)

	select {
	case <-ch:
		return stats
	case <-time.After(timeout):
		return &Stats{Host: host, Timeout: true}
	}
}

func newProxyClient(host, user, password string) (*APIClient, error) {
	log.Debug("call rpc xping to proxy %s", host)
	c := NewAPIClient(host, user, password)
	if err := c.Ping(); err != nil {
		log.Fatal("call rpc xping to proxy failed")
		return c, err
	}
	log.Debug("call rpc xping OK")

	return c, nil
}

// PrepareConfig prepare phase of config change
func PrepareConfig(host, name string, cfg *models.CCConfig) error {
	c, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
	if err != nil {
		log.Fatal("create proxy client failed, %v", err)
		return err
	}

	err = c.PrepareConfig(name)
	if err != nil {
		log.Fatal("prepare proxy config failed, %v", err)
		return err
	}
	return nil
}

// CommitConfig commit phase of config change
func CommitConfig(host, name string, cfg *models.CCConfig) error {
	c, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
	if err != nil {
		log.Fatal("create proxy client failed, %v", err)
		return err
	}
	err = c.CommitConfig(name)
	if err != nil {
		log.Fatal("commit proxy config failed, %v", err)
		return err
	}
	return nil
}

// DelNamespace delete namespace
func DelNamespace(host, name string, cfg *models.CCConfig) error {
	c, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
	if err != nil {
		log.Fatal("create proxy client failed, %v", err)
		return err
	}
	err = c.DelNamespace(name)
	if err != nil {
		log.Warn("delete schema %s in proxy %s failed, %s", name, host, err.Error())
		return err
	}
	return nil
}

// QueryNamespaceSQLFingerprint return sql fingerprint
func QueryNamespaceSQLFingerprint(host, name string, cfg *models.CCConfig) (*SQLFingerprint, error) {
	c, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
	if err != nil {
		log.Fatal("create proxy client failed, %v", err)
		return nil, err
	}

	ret, err := c.GetNamespaceSQLFingerprint(name)
	return ret, err
}

// QueryProxyConfigFingerprint return config fingerprint of proxy
func QueryProxyConfigFingerprint(host string, cfg *models.CCConfig) (string, error) {
	c, err := newProxyClient(host, cfg.ProxyUserName, cfg.ProxyPassword)
	if err != nil {
		return "", err
	}
	configFingerprint, err := c.proxyConfigFingerprint()
	return configFingerprint, err
}
