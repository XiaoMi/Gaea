package config

import (
	"fmt"

	"github.com/XiaoMi/Gaea/models"
)

func GetDefaultKingshardHashRules(db string, table string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "hash",
		Key:   "id",
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
	}
}

func GetDefaultKingshardModRules(db string, table string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "mod",
		Key:   "id",
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
	}
}

func GetDefaultKingshardRangeRules(db string, table string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "range",
		Key:   "id",
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		TableRowLimit: 3,
	}
}

func GetDefaultKingshardDateYearRules(db string, table string, key string, dateRange []string) *models.Shard {
	return &models.Shard{

		DB:    db,
		Table: table,
		Type:  "date_year",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		DateRange: dateRange,
	}
}

func GetDefaultKingshardDateMonthRules(db string, table string, key string, dataRange []string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "date_month",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		DateRange: dataRange,
	}
}

func GetDefaultKingshardDateDayRules(db string, table string, key string, dataRange []string) *models.Shard {
	return &models.Shard{

		DB:    db,
		Table: table,
		Type:  "date_day",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		DateRange:     dataRange,
		TableRowLimit: 3,
	}
}

func GetDefaultMycatshardModRules(db string, table string, key string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "mycat_mod",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		Databases: []string{
			fmt.Sprintf("%s_[0-3]", db),
		},
	}
}

func GetDefaultMycatshardLongRules(db string, table string, key string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "mycat_long",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		Databases: []string{
			fmt.Sprintf("%s_[0-3]", db),
		},
		PartitionCount:  "4",
		PartitionLength: "256",
	}
}

func GetDefaultMycatshardMurmurRules(db string, table string, key string) *models.Shard {
	return &models.Shard{

		DB:    db,
		Table: table,
		Type:  "mycat_murmur",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		Databases: []string{
			fmt.Sprintf("%s_[0-3]", db),
		},
		Seed:               "0",
		VirtualBucketTimes: "160",
	}
}

func GetDefaultMycatshardStringRules(db string, table string, key string) *models.Shard {
	return &models.Shard{
		DB:    db,
		Table: table,
		Type:  "mycat_string",
		Key:   key,
		Locations: []int{
			2,
			2,
		},
		Slices: []string{
			"slice-0",
			"slice-1",
		},
		Databases: []string{
			fmt.Sprintf("%s_[0-3]", db),
		},
		PartitionCount:  "4",
		PartitionLength: "256",
		HashSlice:       ":",
	}
}
