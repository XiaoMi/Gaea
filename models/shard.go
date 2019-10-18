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
	"regexp"
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/core/errors"
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

func (p *Shard) verify() error {

	// get index of linked table config and handle it later
	switch p.Type {
	case ShardLinked:
		return nil
	case ShardDefault:
		return fmt.Errorf("[default-rule] duplicate, must only one")
	}

	if err := verifyRuleSliceInfos(p); err != nil {
		return err
	}

	if IsMycatShardingRule(p.Type) {
		if _, err := getRealDatabases(p.Databases); err != nil {
			return err
		}
	}

	if p.Type == ShardGlobal {
		// 如果全局表指定了物理库名, 则使用mycatDatabases存储这一信息, 否则使用逻辑库名作为物理库名.
		if len(p.Databases) != 0 {
			if _, err := getRealDatabases(p.Databases); err != nil {
				return err
			}
		}
	}

	return nil
}

// Encode encode json
func (p *Shard) Encode() []byte {
	return JSONEncode(p)
}

func verifyRuleSliceInfos(cfg *Shard) error {
	switch cfg.Type {
	case ShardHash:
		_, err := verifyHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return err
		}
		return nil
	case ShardMod:
		_, err := verifyHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return err
		}
		return nil
	case ShardRange:
		tableToSlice, err := verifyHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return err
		}
		tablecount := 0
		for _, location := range cfg.Locations {
			tablecount += location
		}
		if tablecount != len(tableToSlice) {
			return fmt.Errorf("range space %d not equal tables %d", tablecount, len(tableToSlice))
		}
		return nil
	case ShardDay:
		if err := verifyDateDayRuleSliceInfos(cfg.DateRange, cfg.Slices); err != nil {
			return err
		}
		return nil
	case ShardMonth:
		err := verifyDateMonthRuleSliceInfos(cfg.DateRange, cfg.Slices)
		if err != nil {
			return err
		}
		return nil
	case ShardYear:
		err := verifyDateYearRuleSliceInfos(cfg.DateRange, cfg.Slices)
		if err != nil {
			return err
		}
		return nil
	case ShardMycatMod:
		if _, err := verifyMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases); err != nil {
			return err
		}
		return nil
	case ShardMycatLong:
		tableToSlice, err := verifyMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return err
		}
		if err := verifyMycatPatitionLongShard(len(tableToSlice), cfg.PartitionCount, cfg.PartitionLength); err != nil {
			return err
		}
		return nil
	case ShardMycatString:
		tableToSlice, err := verifyMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return err
		}
		if err := verifyMycatPartitionStringShard(len(tableToSlice), cfg.PartitionCount, cfg.PartitionLength, cfg.HashSlice); err != nil {
			return err
		}
		return nil
	case ShardMycatMURMUR:
		tableToSlice, err := verifyMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return err
		}
		if err := verifyMycatPartitionMurmurHashShard(cfg.Seed, cfg.VirtualBucketTimes, len(tableToSlice)); err != nil {
			return err
		}
		return nil
	case ShardMycatPaddingMod:
		tableToSlice, err := verifyMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return err
		}
		if err := verifyMycatPartitionPaddingModShard(cfg.PadFrom, cfg.PadLength, cfg.ModBegin, cfg.ModEnd, len(tableToSlice)); err != nil {
			return err
		}
		return nil
	case ShardGlobal:
		if err := verifyGlobalTableRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases); err != nil {
			return err
		}
		return nil
	default:
		return errors.ErrUnknownRuleType
	}
}

func verifyHashRuleSliceInfos(locations []int, slices []string) (map[int]int, error) {
	var sumTables int
	tableToSlice := make(map[int]int, 0)

	if len(locations) != len(slices) {
		return nil, errors.ErrLocationsCount
	}
	for i := 0; i < len(locations); i++ {
		for j := 0; j < locations[i]; j++ {
			tableToSlice[j+sumTables] = i
		}
		sumTables += locations[i]
	}
	return tableToSlice, nil
}

func verifyMycatHashRuleSliceInfos(locations []int, slices []string, databases []string) (map[int]int, error) {
	tableToSlice, err := verifyHashRuleSliceInfos(locations, slices)
	if err != nil {
		return nil, err
	}

	realDatabaseList, err := getRealDatabases(databases)
	if err != nil {
		return nil, err
	}

	if len(tableToSlice) != len(realDatabaseList) {
		return nil, errors.ErrLocationsCount
	}

	return tableToSlice, nil
}

func verifyDateDayRuleSliceInfos(dateRange []string, slices []string) error {
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	/*check dateRange*/
	return nil
}

func verifyDateMonthRuleSliceInfos(dateRange []string, slices []string) error {
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	/*check dateRange*/
	return nil
}

func verifyDateYearRuleSliceInfos(dateRange []string, slices []string) error {
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	/*check dateRange*/
	return nil
}

func verifyGlobalTableRuleSliceInfos(locations []int, slices []string, databases []string) error {
	tableToSlice, err := verifyHashRuleSliceInfos(locations, slices)
	if err != nil {
		return err
	}

	if len(databases) != 0 {
		realDatabaseList, err := getRealDatabases(databases)
		if err != nil {
			return err
		}
		if len(tableToSlice) != len(realDatabaseList) {
			return errors.ErrLocationsCount
		}
	}

	return nil
}

func includeSlice(slices []string, sliceName string) bool {
	for _, s := range slices {
		if s == sliceName {
			return true
		}
	}
	return false
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

func toIntArray(str string) ([]int, error) {
	str = strings.Replace(str, " ", "", -1)
	strList := strings.Split(str, ",")
	ret := make([]int, 0, len(strList))
	for _, s := range strList {
		num, err := strconv.Atoi(s)
		if err != nil {
			return ret, err
		}
		ret = append(ret, num)
	}
	return ret, nil
}

func verifyMycatPatitionLongShard(shardNum int, partitionCount, partitionLength string) error {
	countList, err := toIntArray(partitionCount)
	if err != nil {
		return err
	}
	lengthList, err := toIntArray(partitionLength)
	if err != nil {
		return err
	}
	countSize := len(countList)
	lengthSize := len(lengthList)
	if countSize != lengthSize {
		return fmt.Errorf("error, check your scope & scopeLength definition")
	}

	segmentLength := 0
	for i := 0; i < countSize; i++ {
		segmentLength += countList[i]
	}

	if segmentLength != shardNum {
		return fmt.Errorf("segmentLength is not equal to shardNum")
	}

	ai := make([]int, segmentLength+1)

	index := 0
	for i := 0; i < countSize; i++ {
		for j := 0; j < countList[i]; j++ {
			ai[index+1] = ai[index] + lengthList[i]
			index++
		}
	}
	if ai[len(ai)-1] != PartitionLength {
		return fmt.Errorf("error, check your partitionScope definition")
	}
	return nil
}

func verifyMycatPartitionStringShard(shardNum int, partitionCount, partitionLength string, hashSliceStr string) error {
	if err := verifyMycatPatitionLongShard(shardNum, partitionCount, partitionLength); err != nil {
		return err
	}
	if err := verifyHashSliceStartEnd(hashSliceStr); err != nil {
		return err
	}
	return nil
}

func verifyHashSliceStartEnd(hashSliceStr string) error {
	hashSliceStr = strings.TrimSpace(hashSliceStr)
	strs := strings.Split(hashSliceStr, ":")

	if len(strs) == 1 {
		_, err := strconv.Atoi(strs[0])
		if err != nil {
			return err
		}
		return nil
	}

	if len(strs) == 2 {
		if err := verifyHashSliceValue(strs[0]); err != nil {
			return fmt.Errorf("parse hash slice start error: %v", err)
		}
		if err := verifyHashSliceValue(strs[1]); err != nil {
			return fmt.Errorf("parse hash slice end error: %v", err)
		}
		return nil
	}

	return fmt.Errorf("invalid hash slice str")
}

func verifyHashSliceValue(str string) error {
	if str == "" {
		return nil
	}
	_, err := strconv.Atoi(str)
	return err
}

func verifyMycatPartitionMurmurHashShard(seedStr, virtualBucketTimesStr string, count int) error {
	_, err := strconv.Atoi(seedStr)
	if err != nil {
		return err
	}
	if virtualBucketTimesStr == "" {
		virtualBucketTimesStr = "160"
	}
	if _, err := strconv.Atoi(virtualBucketTimesStr); err != nil {
		return err
	}
	return nil
}

func verifyMycatPartitionPaddingModShard(padFromStr, padLengthStr, modBeginStr, modEndStr string, mod int) error {
	padFrom, err := strconv.Atoi(padFromStr)
	if err != nil {
		return err
	}
	padLength, err := strconv.Atoi(padLengthStr)
	if err != nil {
		return err
	}
	modBegin, err := strconv.Atoi(modBeginStr)
	if err != nil {
		return err
	}
	modEnd, err := strconv.Atoi(modEndStr)
	if err != nil {
		return err
	}
	if padFrom != PaddingModLeftEnd && padFrom != PaddingModRightEnd {
		return fmt.Errorf("invalid padding mod mode: %d", padFrom)
	}
	if mod < PaddingModDefaultMod {
		return fmt.Errorf("invalid padding mod number: %d", mod)
	}
	if modBegin < 0 || modBegin >= modEnd {
		return fmt.Errorf("invalid padding modBegin or modEnd: %d, %d", modBegin, modEnd)
	}
	if padLength <= 0 {
		return fmt.Errorf("invalid padding mod padLength: %d", padLength)
	}
	if padLength < (modEnd - modBegin) {
		return fmt.Errorf("invalid padding mod, padLength is less than modBegin - modEnd: %d, %d, %d", padLength, modBegin, modEnd)
	}
	return nil
}

func IsMycatShardingRule(ruleType string) bool {
	return ruleType == ShardMod || ruleType == ShardMycatLong || ruleType == ShardMycatMURMUR || ruleType == ShardMycatPaddingMod || ruleType == ShardMycatString
}
