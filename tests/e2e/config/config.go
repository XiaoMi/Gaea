package config

import (
	"fmt"
	"github.com/XiaoMi/Gaea/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

// TestConfig is Global Test configuration.
var TestConfig = NewDefaultConfig("config/config.yaml")

// NewDefaultConfig creates a default configuration.
func NewDefaultConfig(file string) *Config {
	var cfg *Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	return cfg
}

type Config struct {
	GaeaHost     string   `yaml:"gaea_host"`
	GaeaPort     int      `yaml:"gaea_port"`
	GaeaUser     string   `yaml:"gaea_user"`
	GaeaPassword string   `yaml:"gaea_password"`
	Database     string   `yaml:"database"`
	Tables       []string `yaml:"tables"`

	MasterHost     string `yaml:"master_host"`
	MasterPort     int    `yaml:"master_port"`
	MasterUser     string `yaml:"master_user"`
	MasterPassword string `yaml:"master_password"`

	SlaveHost     string `yaml:"slave_host"`
	SlavePort     int    `yaml:"slave_port"`
	SlaveUser     string `yaml:"slave_user"`
	SlavePassword string `yaml:"slave_password"`
}

func (c *Config) GenCreateTableSqls() []string {
	var sqls []string
	for _, t := range c.Tables {
		s := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
  id int(11) NOT NULL,
  name varchar(20) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, c.Database, t)
		sqls = append(sqls, s)
	}
	return sqls
}

func (c *Config) GenDropTableSqls() []string {
	var sqls []string
	for _, t := range c.Tables {
		s := fmt.Sprintf(`DROP TABLE IF EXISTS %s.%s `, c.Database, t)
		sqls = append(sqls, s)
	}
	return sqls
}

func (c *Config) GenDropDatabaseSqls() string {
	return fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, c.Database)
}

func (c *Config) GenCreateDatabaseSqls() string {
	return fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, c.Database)
}

func NewGaeaNsConfig(c *Config) *models.Namespace {
	return &models.Namespace{
		Name:     "test_namespace_1",
		Online:   false,
		ReadOnly: false,
		AllowedDBS: map[string]bool{
			"sbtest1":   true,
			"sbtest2":   true,
			"test_gaea": true,
		},
		SlowSQLTime: "0",
		AllowedIP:   nil,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        c.MasterUser,
				Password:        c.MasterPassword,
				Master:          c.MasterHost + ":" + strconv.Itoa(c.MasterPort),
				Slaves:          []string{c.SlaveHost + ":" + strconv.Itoa(c.SlavePort)},
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		ShardRules: nil,
		Users: []*models.User{
			{
				UserName:      c.GaeaUser,
				Password:      c.GaeaPassword,
				Namespace:     "test_namespace_1",
				RWFlag:        2,
				RWSplit:       1,
				OtherProperty: 0,
			},
		},
		DefaultSlice:      "slice-0",
		GlobalSequences:   nil,
		MaxSqlExecuteTime: 0,
	}
}

var GaeaIniConfig = `
; config type, etcd/file/etcdv3, you can test gaea with file type, you shoud use etcd/etcdv3 in production
; 请指定设定方式为 file 或 etcd 或 etcdv3
config_type=file
;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./test/e2e/config/namespace

;coordinator addr
;coordinator_addr=http://127.0.0.1:2379
;etcd user config
;username=root
;password=root

;environ
environ=local
;service name
service_name=gaea_proxy
;gaea_proxy cluster name
cluster_name=gaea_default_cluster

;log config
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file

;admin addr
admin_addr=0.0.0.0:13307
; basic auth
admin_user=admin
admin_password=admin

;proxy addr
proto_type=tcp4
proxy_addr=0.0.0.0:13306
proxy_charset=utf8
;slow sql time, when execute time is higher than this, log it, unit: ms
slow_sql_time=100
;close session after session timeout, unit: seconds
session_timeout=3600

;stats conf
stats_enabled=true
;stats interval
stats_interval=10

;encrypt key
encrypt_key=1234abcd5678efg*

;server_version
server_version=5.6.20-gaea

;auth plugin mysql_native_password or caching_sha2_password or ''
auth_plugin=mysql_native_password
`
