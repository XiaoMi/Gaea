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
	"github.com/go-ini/ini"
)

// CCConfig means gaea cc config
type CCConfig struct {
	Addr          string `ini:"addr"`
	AdminUserName string `ini:"admin_username"`
	AdminPassword string `ini:"admin_password"`
	ProxyUserName string `ini:"proxy_username"`
	ProxyPassword string `ini:"proxy_password"`
	// etcd 相关配置
	CoordinatorAddr string `ini:"coordinator_addr"`
	UserName        string `ini:"username"`
	Password        string `ini:"password"`

	LogPath     string `ini:"log_path"`
	LogLevel    string `ini:"log_level"`
	LogFileName string `ini:"log_filename"`
	LogOutput   string `ini:"log_output"`
}

// ParseCCConfig parser gaea cc config from file
func ParseCCConfig(cfgFile string) (*CCConfig, error) {
	cfg, err := ini.Load(cfgFile)

	if err != nil {
		return nil, err
	}

	ccConfig := new(CCConfig)
	err = cfg.MapTo(ccConfig)
	return ccConfig, err
}

// Verify verify cc config
func (cc *CCConfig) Verify() error {
	return nil
}
