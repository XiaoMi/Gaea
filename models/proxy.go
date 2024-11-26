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
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/log/zap"
	"github.com/XiaoMi/Gaea/mysql"

	"github.com/go-ini/ini"
)

const (
	defaultGaeaCluster = "gaea"
)

// Proxy means proxy structure of proxy config
type Proxy struct {
	// config type
	ConfigType string `ini:"config_type"`

	// 文件配置类型内容
	FileConfigPath string `ini:"file_config_path"`

	// etcd 相关配置
	CoordinatorAddr string `ini:"coordinator_addr"`
	CoordinatorRoot string `ini:"coordinator_root"`
	UserName        string `ini:"username"`
	Password        string `ini:"password"`

	// 服务相关信息
	Environ string `ini:"environ"`
	Service string `ini:"service_name"`
	Cluster string `ini:"cluster_name"`

	LocalNamespaceStoragePath string `ini:"local_namespace_storage_path"`

	LogPath       string `ini:"log_path"`
	LogLevel      string `ini:"log_level"`
	LogFileName   string `ini:"log_filename"`
	LogOutput     string `ini:"log_output"`
	LogKeepDays   int    `ini:"log_keep_days"`
	LogKeepCounts int    `ini:"log_keep_counts"`

	ProtoType      string `ini:"proto_type"`
	ProxyAddr      string `ini:"proxy_addr"`
	AdminAddr      string `ini:"admin_addr"`
	AdminUser      string `ini:"admin_user"`
	AdminPassword  string `ini:"admin_password"`
	SlowSQLTime    int64  `ini:"slow_sql_time"`
	SessionTimeout int    `ini:"session_timeout"`

	// 监控配置
	StatsEnabled  string `ini:"stats_enabled"`  // set true to enable stats
	StatsInterval int    `ini:"stats_interval"` // set stats interval of connect pool

	EncryptKey string `ini:"encrypt_key"`

	ServerVersion  string `ini:"server_version"`
	AuthPlugin     string `ini:"auth_plugin"`
	NumCPU         int    `ini:"num_cpu"`
	NetBufferSize  int    `ini:"net_buffer_size"`
	DefaultCharset string `ini:"default_charset"`
	ConfigFile     string
}

// ParseProxyConfigFromFile parser proxy config from file
func ParseProxyConfigFromFile(cfgFile string) (*Proxy, error) {
	cfg, err := ini.Load(cfgFile)

	if err != nil {
		return nil, err
	}

	var proxyConfig = &Proxy{}
	err = cfg.MapTo(proxyConfig)
	// default config type: etcd
	if proxyConfig.ConfigType == "" {
		proxyConfig.ConfigType = ConfigEtcd
	}
	if proxyConfig.Cluster == "" && proxyConfig.CoordinatorRoot == "" {
		proxyConfig.Cluster = defaultGaeaCluster
	} else if proxyConfig.Cluster == "" && proxyConfig.CoordinatorRoot != "" {
		proxyConfig.Cluster = strings.TrimPrefix(proxyConfig.CoordinatorRoot, "/")
	} else if proxyConfig.Cluster != "" {
		proxyConfig.CoordinatorRoot = "/" + proxyConfig.Cluster
	}

	proxyConfig.ConfigFile = cfgFile

	if proxyConfig.NetBufferSize > 0 {
		mysql.InitNetBufferSize(proxyConfig.NetBufferSize)
	}

	if err := proxyConfig.Verify(); err != nil {
		return nil, err
	}
	return proxyConfig, err
}

// Verify verify proxy config
func (p *Proxy) Verify() (err error) {

	//first check ConfigType
	switch p.ConfigType {
	case ConfigFile, ConfigEtcd, ConfigEtcdV3:
	default:
		return fmt.Errorf("unsupport config_type: %s", p.ConfigType)
	}

	// check statstics
	if _, err = strconv.ParseBool(p.StatsEnabled); err != nil {
		return fmt.Errorf("StatsEnabled should be a bool value: current: %s, "+
			"error: %s", p.StatsEnabled, err.Error())
	}
	if p.StatsInterval < 0 {
		return fmt.Errorf("stats_interval should be >= 0: %d", p.StatsInterval)
	}

	// check gloal slow query and session timeout
	if p.SlowSQLTime < 0 {
		return fmt.Errorf("slow_sql_time should be >= 0: %d", p.SlowSQLTime)
	}
	if p.SessionTimeout < 0 {
		return fmt.Errorf("session_timeout should be >= 0: %d", p.SlowSQLTime)
	}

	switch p.AuthPlugin {
	case "", mysql.MysqlNativePassword, mysql.CachingSHA2Password:
	default:
		return fmt.Errorf("unsupport auth_plugin: %s", p.AuthPlugin)
	}
	return
}

// ProxyInfo for report proxy information
type ProxyInfo struct {
	Token     string `json:"token"`
	StartTime string `json:"start_time"`

	IP        string `json:"ip"`
	ProtoType string `json:"proto_type"`
	ProxyPort string `json:"proxy_port"`
	AdminPort string `json:"admin_port"`

	Pid int    `json:"pid"`
	Pwd string `json:"pwd"`
	Sys string `json:"sys"`
}

// Encode encode proxy info
func (p *ProxyInfo) Encode() []byte {
	return JSONEncode(p)
}

func InitXLog(output, path, filename, level, service string, logKeepDays int, logKeepCounts int) error {
	cfg := make(map[string]string)
	cfg["path"] = path
	cfg["filename"] = filename
	cfg["level"] = level
	cfg["service"] = service
	cfg["skip"] = "5" // 设置xlog打印方法堆栈需要跳过的层数, 5目前为调用log.Debug()等方法的方法名, 比xlog默认值多一层.
	cfg["log_keep_days"] = strconv.Itoa(log.DefaultLogKeepDays)
	if logKeepDays != 0 {
		cfg["log_keep_days"] = strconv.Itoa(logKeepDays)
	}

	cfg["log_keep_counts"] = strconv.Itoa(log.DefaultLogKeepCounts)
	if logKeepCounts != 0 {
		cfg["log_keep_counts"] = strconv.Itoa(logKeepCounts)
	}
	logger, err := zap.CreateLogManager(cfg)
	if err != nil {
		return err
	}

	log.SetGlobalLogger(logger)
	return nil
}
