package config

import (
	"github.com/XiaoMi/Gaea/models"
)

type ShardOption func(*models.Shard)

func NewShard(opts ...ShardOption) *models.Shard {
	shard := &models.Shard{}
	for _, opt := range opts {
		opt(shard)
	}
	return shard
}

func WithDB(db string) ShardOption {
	return func(s *models.Shard) {
		s.DB = db
	}
}

func WithTable(table string) ShardOption {
	return func(s *models.Shard) {
		s.Table = table
	}
}

func WithParentTable(table string) ShardOption {
	return func(s *models.Shard) {
		s.ParentTable = table
	}
}

func WithType(t string) ShardOption {
	return func(s *models.Shard) {
		s.Type = t
	}
}

func WithKey(key string) ShardOption {
	return func(s *models.Shard) {
		s.Key = key
	}
}

func WithLocations(locations []int) ShardOption {
	return func(s *models.Shard) {
		s.Locations = locations
	}
}

func WithTableRowLimit(tableRowLimit int) ShardOption {
	return func(s *models.Shard) {
		s.TableRowLimit = tableRowLimit
	}
}

func WithDatabases(datavases []string) ShardOption {
	return func(s *models.Shard) {
		s.Databases = datavases
	}
}

func WithPartitionLength(l string) ShardOption {
	return func(s *models.Shard) {
		s.PartitionLength = l
	}
}

func WithPartitionCount(c string) ShardOption {
	return func(s *models.Shard) {
		s.PartitionCount = c
	}
}

func WithHashSlice(slice string) ShardOption {
	return func(s *models.Shard) {
		s.HashSlice = slice
	}
}

func WithVirtualBucketTimes(v string) ShardOption {
	return func(s *models.Shard) {
		s.VirtualBucketTimes = v
	}
}

func WithSeed(seed string) ShardOption {
	return func(s *models.Shard) {
		s.Seed = seed
	}
}
func WithShardSlice(slice []string) ShardOption {
	return func(s *models.Shard) {
		s.Slices = slice
	}
}

func WithDateRange(daterange []string) ShardOption {
	return func(s *models.Shard) {
		s.DateRange = daterange
	}
}
