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
	"github.com/XiaoMi/Gaea/core/errors"
	"regexp"
	"strconv"
)

// constants of shard type
const (
	ShardDefault         = "default"
	ShardGlobal          = "global"
	ShardLinked          = "linked"
	ShardMod             = "mod"
	ShardHash            = "hash"
	ShardRange           = "range"
	ShardYear            = "date_year"
	ShardMonth           = "date_month"
	ShardDay             = "date_day"
	ShardMycatMod        = "mycat_mod"
	ShardMycatLong       = "mycat_long"
	ShardMycatString     = "mycat_string"
	ShardMycatMURMUR     = "mycat_murmur"
	ShardMycatPaddingMod = "mycat_padding_mod"

	// PartitionLength length of partition
	PartitionLength = 1024

	// mod padding
	PaddingModLeftEnd  = 0
	PaddingModRightEnd = 1

	PaddingModDefaultPadFrom   = PaddingModRightEnd
	PaddingModDefaultPadLength = 18
	PaddingModDefaultModBegin  = 10
	PaddingModDefaultModEnd    = 16
	PaddingModDefaultMod       = 2
)

// Shard means shard model in etcd
type Shard struct {
	DB            string   `json:"db"`
	Table         string   `json:"table"`
	ParentTable   string   `json:"parent_table"`
	Type          string   `json:"type"` // 表类型: 包括分表如hash/range/data,关联表如: linked 全局表如: global等
	Key           string   `json:"key"`
	Locations     []int    `json:"locations"`
	Slices        []string `json:"slices"`
	DateRange     []string `json:"date_range"`
	TableRowLimit int      `json:"table_row_limit"`

	// only used in mycat logic database (schema)
	Databases []string `json:"databases"`

	// used in mycat partition long shard and partition string shard
	PartitionCount  string `json:"partition_count"`
	PartitionLength string `json:"partition_length"`

	// used in mycat partition string shard
	HashSlice string `json:"hash_slice"`

	// used in mycat murmur hash shard
	Seed               string `json:"seed"`
	VirtualBucketTimes string `json:"virtual_bucket_times"`

	// used in mycat padding mod shard
	PadFrom   string `json:"pad_from"`
	PadLength string `json:"pad_length"`
	ModBegin  string `json:"mod_begin"`
	ModEnd    string `json:"mod_end"`
}

func (s *Shard) verify() error {
	if err := s.verifyRuleSliceInfos(); err != nil {
		return err
	}
	return nil
}

func (s *Shard) verifyRuleSliceInfos() error {
	f, ok := ruleVerifyFuncMapping[s.Type]
	if !ok {
		return errors.ErrUnknownRuleType
	}
	return f(s)
}

// Encode encode json
func (s *Shard) Encode() []byte {
	return JSONEncode(s)
}

func IsMycatShardingRule(ruleType string) bool {
	return ruleType == ShardMod || ruleType == ShardMycatLong || ruleType == ShardMycatMURMUR || ruleType == ShardMycatPaddingMod || ruleType == ShardMycatString
}

var rangeDatabaseRegex = regexp.MustCompile(`^(\S+?)\[(\d+)-(\d+)\]$`)

// if a dbname is a database list, then parse the real dbnames and add to the result.
// the range contains left bound and right bound, which means [left, right].
func getRealDatabases(dbs []string) ([]string, error) {
	var ret []string
	for _, db := range dbs {
		if rangeDatabaseRegex.MatchString(db) {
			matches := rangeDatabaseRegex.FindStringSubmatch(db)
			if len(matches) != 4 {
				return nil, fmt.Errorf("invalid database list: %s", db)
			}
			dbPrefix := matches[1]
			leftBoundStr := matches[2]
			rightBoundStr := matches[3]
			leftBound, err := strconv.Atoi(leftBoundStr)
			if err != nil {
				return nil, fmt.Errorf("invalid left bound value of database list: %s", db)
			}
			rightBound, err := strconv.Atoi(rightBoundStr)
			if err != nil {
				return nil, fmt.Errorf("invalid right bound value of database list: %s", db)
			}
			if rightBound <= leftBound {
				return nil, fmt.Errorf("invalid bound value of database list: %s", db)
			}
			for i := leftBound; i <= rightBound; i++ {
				realDB := dbPrefix + strconv.Itoa(i)
				ret = append(ret, realDB)
			}
		} else {
			ret = append(ret, db)
		}
	}
	return ret, nil
}
