// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

package router

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/models"
)

const (
	DefaultRuleType         = models.ShardDefault
	GlobalTableRuleType     = models.ShardGlobal
	LinkedTableRuleType     = models.ShardLinked // this type only exists in conf, then transfer to LinkedRule
	HashRuleType            = models.ShardHash
	RangeRuleType           = models.ShardRange
	ModRuleType             = models.ShardMod
	DateYearRuleType        = models.ShardYear
	DateMonthRuleType       = models.ShardMonth
	DateDayRuleType         = models.ShardDay
	MycatModRuleType        = models.ShardMycatMod
	MycatLongRuleType       = models.ShardMycatLong
	MycatStringRuleType     = models.ShardMycatString
	MycatMurmurRuleType     = models.ShardMycatMURMUR
	MycatPaddingModRuleType = models.ShardMycatPaddingMod

	MinMonthDaysCount = 28
	MaxMonthDaysCount = 31
	MonthsCount       = 12
)

type Rule interface {
	GetDB() string
	GetTable() string
	GetShardingColumn() string
	IsLinkedRule() bool
	GetShard() Shard
	FindTableIndex(key interface{}) (int, error)
	GetSlice(i int) string // i is slice index
	GetSliceIndexFromTableIndex(i int) int
	GetSlices() []string
	GetSubTableIndexes() []int
	GetFirstTableIndex() int
	GetLastTableIndex() int
	GetType() string
	GetDatabaseNameByTableIndex(index int) (string, error)
}

type MycatRule interface {
	Rule
	GetDatabases() []string
	GetTableIndexByDatabaseName(phyDB string) (int, bool)
}

type BaseRule struct {
	db             string
	table          string
	shardingColumn string

	ruleType        string
	slices          []string    // not the namespace slices
	subTableIndexes []int       //subTableIndexes store all the index of sharding sub-table
	tableToSlice    map[int]int //key is table index, and value is slice index
	shard           Shard

	// TODO: 目前全局表也借用这两个field存放默认分片的物理DB名
	mycatDatabases               []string
	mycatDatabaseToTableIndexMap map[string]int // key: phy db name, value: table index
}

type LinkedRule struct {
	db             string
	table          string
	shardingColumn string

	linkToRule *BaseRule
}

func NewDefaultRule(slice string) *BaseRule {
	var r *BaseRule = &BaseRule{
		ruleType:     DefaultRuleType,
		slices:       []string{slice},
		shard:        new(DefaultShard),
		tableToSlice: nil,
	}
	return r
}

func (r *BaseRule) GetDB() string {
	return r.db
}

func (r *BaseRule) GetTable() string {
	return r.table
}

func (r *BaseRule) GetShardingColumn() string {
	return r.shardingColumn
}

func (r *BaseRule) IsLinkedRule() bool {
	return false
}

func (r *BaseRule) GetShard() Shard {
	return r.shard
}

func (r *BaseRule) FindTableIndex(key interface{}) (int, error) {
	return r.shard.FindForKey(key)
}

// The confs should be verified before use to avoid panic.
func (r *BaseRule) GetSlice(i int) string {
	return r.slices[i]
}

func (r *BaseRule) GetSliceIndexFromTableIndex(i int) int {
	sliceIndex, ok := r.tableToSlice[i]
	if !ok {
		return -1
	}
	return sliceIndex
}

// This is dangerous since the caller can change the value in slices.
// It is better to return a iterator instead of exposing the origin slices.
func (r *BaseRule) GetSlices() []string {
	return r.slices
}

func (r *BaseRule) GetSubTableIndexes() []int {
	return r.subTableIndexes
}

func (r *BaseRule) GetFirstTableIndex() int {
	return r.subTableIndexes[0]
}

func (r *BaseRule) GetLastTableIndex() int {
	return r.subTableIndexes[len(r.subTableIndexes)-1]
}

func (r *BaseRule) GetType() string {
	return r.ruleType
}

func (r *BaseRule) GetDatabaseNameByTableIndex(index int) (string, error) {
	if IsMycatShardingRule(r.ruleType) || r.ruleType == GlobalTableRuleType {
		if index > len(r.subTableIndexes) {
			return "", errors.ErrInvalidArgument
		}
		return r.mycatDatabases[index], nil
	}
	return r.db, nil
}

func (r *BaseRule) GetTableIndexByDatabaseName(phyDB string) (int, bool) {
	idx, ok := r.mycatDatabaseToTableIndexMap[phyDB]
	return idx, ok
}

func (r *BaseRule) GetDatabases() []string {
	return r.mycatDatabases
}

func (l *LinkedRule) GetDB() string {
	return l.db
}

func (l *LinkedRule) GetTable() string {
	return l.table
}

func (l *LinkedRule) GetParentDB() string {
	return l.linkToRule.GetDB()
}

func (l *LinkedRule) GetParentTable() string {
	return l.linkToRule.GetTable()
}

func (l *LinkedRule) GetShardingColumn() string {
	return l.shardingColumn
}

func (l *LinkedRule) IsLinkedRule() bool {
	return true
}

func (l *LinkedRule) GetShard() Shard {
	return l.linkToRule.GetShard()
}

func (l *LinkedRule) FindTableIndex(key interface{}) (int, error) {
	return l.linkToRule.FindTableIndex(key)
}

func (l *LinkedRule) GetFirstTableIndex() int {
	return l.linkToRule.GetFirstTableIndex()
}

func (l *LinkedRule) GetLastTableIndex() int {
	return l.linkToRule.GetLastTableIndex()
}

func (l *LinkedRule) GetSlice(i int) string {
	return l.linkToRule.GetSlice(i)
}

func (l *LinkedRule) GetSliceIndexFromTableIndex(i int) int {
	return l.linkToRule.GetSliceIndexFromTableIndex(i)
}

func (l *LinkedRule) GetSlices() []string {
	return l.linkToRule.GetSlices()
}

func (l *LinkedRule) GetSubTableIndexes() []int {
	return l.linkToRule.GetSubTableIndexes()
}

func (l *LinkedRule) GetType() string {
	return l.linkToRule.GetType()
}

func (l *LinkedRule) GetDatabaseNameByTableIndex(index int) (string, error) {
	return l.linkToRule.GetDatabaseNameByTableIndex(index)
}

func (l *LinkedRule) GetDatabases() []string {
	return l.linkToRule.GetDatabases()
}

func (l *LinkedRule) GetTableIndexByDatabaseName(phyDB string) (int, bool) {
	return l.linkToRule.GetTableIndexByDatabaseName(phyDB)
}

func createLinkedRule(rules map[string]map[string]Rule, shard *models.Shard) (*LinkedRule, error) {
	if shard.Type != LinkedTableRuleType {
		return nil, fmt.Errorf("LinkedRule type is not linked: %v", shard)
	}

	tableRules, ok := rules[shard.DB]
	if !ok {
		return nil, fmt.Errorf("db of LinkedRule is not found in parent rules")
	}
	dbRule, ok := tableRules[strings.ToLower(shard.ParentTable)]
	if !ok {
		return nil, fmt.Errorf("parent table of LinkedRule is not found in parent rules")
	}
	if dbRule.GetType() == LinkedTableRuleType {
		return nil, fmt.Errorf("LinkedRule cannot link to another LinkedRule")
	}
	linkToRule, ok := dbRule.(*BaseRule)
	if !ok {
		return nil, fmt.Errorf("LinkedRule must link to a BaseRule")
	}

	linkedRule := &LinkedRule{
		db:             shard.DB,
		table:          strings.ToLower(shard.Table),
		shardingColumn: strings.ToLower(shard.Key),
		linkToRule:     linkToRule,
	}

	return linkedRule, nil
}

func parseRule(cfg *models.Shard) (*BaseRule, error) {
	r := new(BaseRule)
	r.db = cfg.DB
	r.table = strings.ToLower(cfg.Table)
	r.shardingColumn = strings.ToLower(cfg.Key) //ignore case
	r.ruleType = cfg.Type
	r.slices = cfg.Slices //将rule model中的slices赋值给rule
	r.mycatDatabaseToTableIndexMap = make(map[string]int)

	subTableIndexs, tableToSlice, shard, err := parseRuleSliceInfos(cfg)
	if err != nil {
		return nil, err
	}

	r.subTableIndexes = subTableIndexs
	r.tableToSlice = tableToSlice
	r.shard = shard

	if IsMycatShardingRule(cfg.Type) {
		r.mycatDatabases, err = getRealDatabases(cfg.Databases)
		if err != nil {
			return nil, err
		}
		for i, db := range r.mycatDatabases {
			r.mycatDatabaseToTableIndexMap[db] = i
		}
	}

	if cfg.Type == GlobalTableRuleType {
		// 如果全局表指定了物理库名, 则使用mycatDatabases存储这一信息, 否则使用逻辑库名作为物理库名.
		if len(cfg.Databases) != 0 {
			r.mycatDatabases, err = getRealDatabases(cfg.Databases)
			if err != nil {
				return nil, err
			}
		} else {
			for i := 0; i < len(r.subTableIndexes); i++ {
				r.mycatDatabases = append(r.mycatDatabases, r.db)
			}
		}
		for i, db := range r.mycatDatabases {
			r.mycatDatabaseToTableIndexMap[db] = i
		}
	}

	return r, nil
}

func parseRuleSliceInfos(cfg *models.Shard) ([]int, map[int]int, Shard, error) {
	switch cfg.Type {
	case HashRuleType:
		subTableIndexs, tableToSlice, err := parseHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := &HashShard{ShardNum: len(tableToSlice)}
		return subTableIndexs, tableToSlice, shard, nil
	case ModRuleType:
		subTableIndexs, tableToSlice, err := parseHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := &ModShard{ShardNum: len(tableToSlice)}
		return subTableIndexs, tableToSlice, shard, nil
	case RangeRuleType:
		subTableIndexs, tableToSlice, err := parseHashRuleSliceInfos(cfg.Locations, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		rs, err := ParseNumSharding(cfg.Locations, cfg.TableRowLimit)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(rs) != len(tableToSlice) {
			return nil, nil, nil, fmt.Errorf("range space %d not equal tables %d", len(rs), len(tableToSlice))
		}
		shard := &NumRangeShard{Shards: rs}
		return subTableIndexs, tableToSlice, shard, nil
	case DateDayRuleType:
		subTableIndexs, tableToSlice, err := parseDateDayRuleSliceInfos(cfg.DateRange, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := &DateDayShard{}
		return subTableIndexs, tableToSlice, shard, nil
	case DateMonthRuleType:
		subTableIndexs, tableToSlice, err := parseDateMonthRuleSliceInfos(cfg.DateRange, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := &DateMonthShard{}
		return subTableIndexs, tableToSlice, shard, nil
	case DateYearRuleType:
		subTableIndexs, tableToSlice, err := parseDateYearRuleSliceInfos(cfg.DateRange, cfg.Slices)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := &DateYearShard{}
		return subTableIndexs, tableToSlice, shard, nil
	case MycatModRuleType:
		subTableIndexs, tableToSlice, err := parseMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := NewMycatPartitionModShard(len(tableToSlice))
		return subTableIndexs, tableToSlice, shard, nil
	case MycatLongRuleType:
		subTableIndexs, tableToSlice, err := parseMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := NewMycatPartitionLongShard(len(tableToSlice), cfg.PartitionCount, cfg.PartitionLength)
		if err = shard.Init(); err != nil {
			return nil, nil, nil, err
		}
		return subTableIndexs, tableToSlice, shard, nil
	case MycatStringRuleType:
		subTableIndexs, tableToSlice, err := parseMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := NewMycatPartitionStringShard(len(tableToSlice), cfg.PartitionCount, cfg.PartitionLength, cfg.HashSlice)
		if err = shard.Init(); err != nil {
			return nil, nil, nil, err
		}
		return subTableIndexs, tableToSlice, shard, nil
	case MycatMurmurRuleType:
		subTableIndexs, tableToSlice, err := parseMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}

		shard, err := NewMycatPartitionMurmurHashShard(cfg.Seed, cfg.VirtualBucketTimes, len(tableToSlice))
		if err != nil {
			return nil, nil, nil, err
		}
		if err = shard.Init(); err != nil {
			return nil, nil, nil, err
		}
		return subTableIndexs, tableToSlice, shard, nil
	case MycatPaddingModRuleType:
		subTableIndexs, tableToSlice, err := parseMycatHashRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}

		shard, err := GetMycatPartitionPaddingModShard(cfg.PadFrom, cfg.PadLength, cfg.ModBegin, cfg.ModEnd, len(tableToSlice))
		if err != nil {
			return nil, nil, nil, err
		}
		if err = shard.Init(); err != nil {
			return nil, nil, nil, err
		}
		return subTableIndexs, tableToSlice, shard, nil
	case GlobalTableRuleType:
		subTableIndexs, tableToSlice, err := parseGlobalTableRuleSliceInfos(cfg.Locations, cfg.Slices, cfg.Databases)
		if err != nil {
			return nil, nil, nil, err
		}
		shard := NewGlobalTableShard()
		return subTableIndexs, tableToSlice, shard, nil
	default:
		return nil, nil, nil, errors.ErrUnknownRuleType
	}
}

func parseHashRuleSliceInfos(locations []int, slices []string) ([]int, map[int]int, error) {
	var sumTables int
	var subTableIndexs []int
	tableToSlice := make(map[int]int, 0)

	if len(locations) != len(slices) {
		return nil, nil, errors.ErrLocationsCount
	}
	for i := 0; i < len(locations); i++ {
		for j := 0; j < locations[i]; j++ {
			subTableIndexs = append(subTableIndexs, j+sumTables)
			tableToSlice[j+sumTables] = i
		}
		sumTables += locations[i]
	}
	return subTableIndexs, tableToSlice, nil
}

func parseMycatHashRuleSliceInfos(locations []int, slices []string, databases []string) ([]int, map[int]int, error) {
	subTableIndexs, tableToSlice, err := parseHashRuleSliceInfos(locations, slices)
	if err != nil {
		return nil, nil, err
	}

	realDatabaseList, err := getRealDatabases(databases)
	if err != nil {
		return nil, nil, err
	}

	if len(tableToSlice) != len(realDatabaseList) {
		return nil, nil, errors.ErrLocationsCount
	}

	return subTableIndexs, tableToSlice, nil
}

func parseDateDayRuleSliceInfos(dateRange []string, slices []string) ([]int, map[int]int, error) {
	var subTableIndexs []int
	tableToSlice := make(map[int]int, 0)

	if len(dateRange) != len(slices) {
		return nil, nil, errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		dayNumbers, err := ParseDayRange(dateRange[i])
		if err != nil {
			return nil, nil, err
		}
		if len(subTableIndexs) > 0 && dayNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return nil, nil, errors.ErrDateRangeOverlap
		}
		for _, v := range dayNumbers {
			subTableIndexs = append(subTableIndexs, v)
			tableToSlice[v] = i
		}
	}
	return subTableIndexs, tableToSlice, nil
}

func parseDateMonthRuleSliceInfos(dateRange []string, slices []string) ([]int, map[int]int, error) {
	var subTableIndexs []int
	tableToSlice := make(map[int]int, 0)

	if len(dateRange) != len(slices) {
		return nil, nil, errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		monthNumbers, err := ParseMonthRange(dateRange[i])
		if err != nil {
			return nil, nil, err
		}
		if len(subTableIndexs) > 0 && monthNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return nil, nil, errors.ErrDateRangeOverlap
		}
		for _, v := range monthNumbers {
			subTableIndexs = append(subTableIndexs, v)
			tableToSlice[v] = i
		}
	}
	return subTableIndexs, tableToSlice, nil
}

func parseDateYearRuleSliceInfos(dateRange []string, slices []string) ([]int, map[int]int, error) {
	var subTableIndexs []int
	tableToSlice := make(map[int]int, 0)

	if len(dateRange) != len(slices) {
		return nil, nil, errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		yearNumbers, err := ParseYearRange(dateRange[i])
		if err != nil {
			return nil, nil, err
		}
		if len(subTableIndexs) > 0 && yearNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return nil, nil, errors.ErrDateRangeOverlap
		}
		for _, v := range yearNumbers {
			tableToSlice[v] = i
			subTableIndexs = append(subTableIndexs, v)
		}
	}
	return subTableIndexs, tableToSlice, nil
}

func parseGlobalTableRuleSliceInfos(locations []int, slices []string, databases []string) ([]int, map[int]int, error) {
	subTableIndexs, tableToSlice, err := parseHashRuleSliceInfos(locations, slices)
	if err != nil {
		return nil, nil, err
	}

	if len(databases) != 0 {
		realDatabaseList, err := getRealDatabases(databases)
		if err != nil {
			return nil, nil, err
		}
		if len(tableToSlice) != len(realDatabaseList) {
			return nil, nil, errors.ErrLocationsCount
		}
	}

	return subTableIndexs, tableToSlice, nil
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

func IsMycatShardingRule(ruleType string) bool {
	return ruleType == MycatModRuleType || ruleType == MycatLongRuleType || ruleType == MycatMurmurRuleType || ruleType == MycatPaddingModRuleType || ruleType == MycatStringRuleType
}
