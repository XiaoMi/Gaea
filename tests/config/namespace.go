package config

import (
	"github.com/XiaoMi/Gaea/models"
)

type Option func(p *models.Namespace)

func WithName(name string) Option {
	return func(n *models.Namespace) {
		n.Name = name
	}
}

func WithOnline(online bool) Option {
	return func(n *models.Namespace) {
		n.Online = online
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(n *models.Namespace) {
		n.ReadOnly = readOnly
	}
}

func WithAllowedDBS(allowedDBS map[string]bool) Option {
	return func(n *models.Namespace) {
		n.AllowedDBS = allowedDBS
	}
}

func WithSlowSqlTime(slowSQLTime string) Option {
	return func(n *models.Namespace) {
		n.SlowSQLTime = slowSQLTime
	}
}

func WithBlackSql(blackSqls []string) Option {
	return func(n *models.Namespace) {
		n.BlackSQL = blackSqls
	}
}

func WithAllowedIp(allowedIp []string) Option {
	return func(n *models.Namespace) {
		n.AllowedIP = allowedIp
	}
}

func WithSlices(slices []*models.Slice) Option {
	return func(n *models.Namespace) {
		n.Slices = slices
	}
}

func WithShardRules(shard []*models.Shard) Option {
	return func(n *models.Namespace) {
		n.ShardRules = shard
	}
}

func WithUsers(users []*models.User) Option {
	return func(n *models.Namespace) {
		n.Users = users
	}
}

func WithDefaultSlice(defaultSlice string) Option {
	return func(n *models.Namespace) {
		n.DefaultSlice = defaultSlice
	}
}

func WithGlobalSequences(globalSequence []*models.GlobalSequence) Option {
	return func(n *models.Namespace) {
		n.GlobalSequences = globalSequence
	}
}

func WithMaxSqlExecuteTime(maxSqlExecuteTime int) Option {
	return func(n *models.Namespace) {
		n.MaxSqlExecuteTime = maxSqlExecuteTime
	}
}

func WithOpenGeneralLog(openGeneralLog bool) Option {
	return func(n *models.Namespace) {
		n.OpenGeneralLog = openGeneralLog
	}
}

func WithIsEncrypt(isEncrypt bool) Option {
	return func(n *models.Namespace) {
		n.IsEncrypt = isEncrypt
	}
}

func WithDefaultPhyDBs(phyDB map[string]string) Option {
	return func(n *models.Namespace) {
		n.DefaultPhyDBS = phyDB
	}
}

func NewNamespace(options ...Option) *models.Namespace {
	n := &models.Namespace{
		IsEncrypt:         true,
		Online:            true,
		ReadOnly:          false,
		AllowedDBS:        make(map[string]bool),
		SlowSQLTime:       "1000",
		BlackSQL:          []string{},
		Slices:            []*models.Slice{},
		ShardRules:        []*models.Shard{},
		Users:             []*models.User{},
		MaxSqlExecuteTime: 0,
		DefaultSlice:      "slice-0",
	}
	for _, option := range options {
		option(n)
	}
	return n
}
