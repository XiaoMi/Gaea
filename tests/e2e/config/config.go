package config

import (
	"github.com/XiaoMi/Gaea/models"
)

type DBConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Url          string `yaml:"url"`
	Timeout      int    `yaml:"timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	ParseTime    bool   `yaml:"parse_time"`
	Loc          string `yaml:"loc"`
	Charset      string `yaml:"charset"`
}

type DBOption func(*DBConfig)

func NewDBConfig(options ...DBOption) *DBConfig {
	config := &DBConfig{
		Timeout:      30,
		ReadTimeout:  30,
		WriteTimeout: 30,
		ParseTime:    true,
		Loc:          "Local",
		Charset:      "utf8mb4",
	}
	for _, option := range options {
		option(config)
	}
	return config
}

func WithUser(user string) DBOption {
	return func(config *DBConfig) {
		config.User = user
	}
}

func WithPassword(password string) DBOption {
	return func(config *DBConfig) {
		config.Password = password
	}
}

func WithUrl(url string) DBOption {
	return func(config *DBConfig) {
		config.Url = url
	}
}

func WithTimeOut(out int) DBOption {
	return func(config *DBConfig) {
		config.Timeout = out
	}
}

func WithReadTimeOut(out int) DBOption {
	return func(config *DBConfig) {
		config.ReadTimeout = out
	}
}

func WithWriteTimeOut(out int) DBOption {
	return func(config *DBConfig) {
		config.WriteTimeout = out
	}
}

func WithLoc(loc string) DBOption {
	return func(config *DBConfig) {
		config.Loc = loc
	}
}

func WithCharSet(charset string) DBOption {
	return func(config *DBConfig) {
		config.Charset = charset
	}
}

func WithParseTime(prase bool) DBOption {
	return func(config *DBConfig) {
		config.ParseTime = prase
	}
}

type MysqlClusterManager struct {
	MysqlSlices map[string]MysqlSlice
}

type MysqlSlice struct {
	Master *DBConfig
	Slaves []*DBConfig
}

func NewMysqlClusterManager(slices []*models.Slice) MysqlClusterManager {
	m := MysqlClusterManager{
		MysqlSlices: map[string]MysqlSlice{},
	}
	for _, s := range slices {
		if s != nil {
			master := NewDBConfig(
				WithUser(s.UserName),
				WithPassword(s.Password),
				WithUrl(s.Master),
			)
			slaves := []*DBConfig{}
			for _, l := range s.Slaves {
				db := NewDBConfig(
					WithUser(s.UserName),
					WithPassword(s.Password),
					WithUrl(l),
				)
				slaves = append(slaves, db)
			}
			m.MysqlSlices[s.Name] = MysqlSlice{
				Master: master,
				Slaves: slaves,
			}
		}
	}
	return m

}

func (m MysqlClusterManager) GetMysqlSliceByName(name string) MysqlSlice {
	if _, ok := m.MysqlSlices[name]; ok {
		return m.MysqlSlices[name]
	}
	return MysqlSlice{}
}

func NewGaeaDBConfig(users []*models.User, url string) []*DBConfig {
	res := []*DBConfig{}
	for _, v := range users {
		if v != nil {
			db := NewDBConfig(
				WithUser(v.UserName),
				WithPassword(v.Password),
				WithUrl(url),
			)
			res = append(res, db)
		}
	}
	return res
}
