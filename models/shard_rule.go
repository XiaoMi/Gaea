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

package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/core/errors"
)

var ruleVerifyFuncMapping = map[string]func(shard *Shard) error{
	ShardHash:            verifyHashRule,
	ShardMod:             verifyModRule,
	ShardRange:           verifyRangeRule,
	ShardDay:             verifyDayRule,
	ShardMonth:           verifyMonthRule,
	ShardYear:            verifyYearRule,
	ShardMycatMod:        verifyMycatModRule,
	ShardMycatLong:       verifyMycatLongRule,
	ShardMycatString:     verifyMycatStringRule,
	ShardMycatMURMUR:     verifyMycatMURMURRule,
	ShardMycatPaddingMod: verifyMycatPaddingRule,
	ShardGlobal:          verifyGlobalRule,
}

func verifyHashRule(s *Shard) error {
	if _, err := verifyHashRuleSliceInfos(s.Locations, s.Slices); err != nil {
		return err
	}
	return nil
}

func verifyModRule(s *Shard) error {
	if _, err := verifyHashRuleSliceInfos(s.Locations, s.Slices); err != nil {
		return err
	}
	return nil
}

func verifyRangeRule(s *Shard) error {
	tableToSlice, err := verifyHashRuleSliceInfos(s.Locations, s.Slices)
	if err != nil {
		return err
	}

	tableCount, err := ParseNumSharding(s.Locations, s.TableRowLimit)
	if err != nil {
		return err
	}
	if len(tableCount) != len(tableToSlice) {
		return fmt.Errorf("range space %d not equal tables %d", tableCount, len(tableToSlice))
	}
	return nil
}

func verifyDayRule(s *Shard) error {
	if err := verifyDateDayRuleSliceInfos(s.DateRange, s.Slices); err != nil {
		return err
	}
	return nil
}

func verifyMonthRule(s *Shard) error {
	err := verifyDateMonthRuleSliceInfos(s.DateRange, s.Slices)
	if err != nil {
		return err
	}
	return nil
}

func verifyYearRule(s *Shard) error {
	err := verifyDateYearRuleSliceInfos(s.DateRange, s.Slices)
	if err != nil {
		return err
	}
	return nil
}

func verifyMycatModRule(s *Shard) error {
	if _, err := verifyMycatHashRuleSliceInfos(s.Locations, s.Slices, s.Databases); err != nil {
		return err
	}
	return nil
}

func verifyMycatLongRule(s *Shard) error {
	tableToSlice, err := verifyMycatHashRuleSliceInfos(s.Locations, s.Slices, s.Databases)
	if err != nil {
		return err
	}
	if err := verifyMycatPatitionLongShard(len(tableToSlice), s.PartitionCount, s.PartitionLength); err != nil {
		return err
	}
	return nil
}

func verifyMycatStringRule(s *Shard) error {
	tableToSlice, err := verifyMycatHashRuleSliceInfos(s.Locations, s.Slices, s.Databases)
	if err != nil {
		return err
	}
	if err := verifyMycatPartitionStringShard(len(tableToSlice), s.PartitionCount, s.PartitionLength, s.HashSlice); err != nil {
		return err
	}
	return nil
}

func verifyMycatMURMURRule(s *Shard) error {
	tableToSlice, err := verifyMycatHashRuleSliceInfos(s.Locations, s.Slices, s.Databases)
	if err != nil {
		return err
	}
	if err := verifyMycatPartitionMurmurHashShard(s.Seed, s.VirtualBucketTimes, len(tableToSlice)); err != nil {
		return err
	}
	return nil
}

func verifyMycatPaddingRule(s *Shard) error {
	tableToSlice, err := verifyMycatHashRuleSliceInfos(s.Locations, s.Slices, s.Databases)
	if err != nil {
		return err
	}
	if err := verifyMycatPartitionPaddingModShard(s.PadFrom, s.PadLength, s.ModBegin, s.ModEnd, len(tableToSlice)); err != nil {
		return err
	}
	return nil
}

func verifyGlobalRule(s *Shard) error {
	if err := verifyGlobalTableRuleSliceInfos(s.Locations, s.Slices, s.Databases); err != nil {
		return err
	}
	return nil
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
	var subTableIndexs []int
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		dayNumbers, err := ParseDayRange(dateRange[i])
		if err != nil {
			return err
		}
		if len(subTableIndexs) > 0 && dayNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return errors.ErrDateRangeOverlap
		}
		for _, v := range dayNumbers {
			subTableIndexs = append(subTableIndexs, v)
		}
	}
	return nil
}

func verifyDateMonthRuleSliceInfos(dateRange []string, slices []string) error {
	var subTableIndexs []int
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		monthNumbers, err := ParseMonthRange(dateRange[i])
		if err != nil {
			return err
		}
		if len(subTableIndexs) > 0 && monthNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return errors.ErrDateRangeOverlap
		}
		for _, v := range monthNumbers {
			subTableIndexs = append(subTableIndexs, v)
		}
	}
	return nil
}

func verifyDateYearRuleSliceInfos(dateRange []string, slices []string) error {
	var subTableIndexs []int
	if len(dateRange) != len(slices) {
		return errors.ErrDateRangeCount
	}
	for i := 0; i < len(dateRange); i++ {
		yearNumbers, err := ParseYearRange(dateRange[i])
		if err != nil {
			return err
		}
		if len(subTableIndexs) > 0 && yearNumbers[0] <= subTableIndexs[len(subTableIndexs)-1] {
			return errors.ErrDateRangeOverlap
		}
		for _, v := range yearNumbers {
			subTableIndexs = append(subTableIndexs, v)
		}
	}
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
