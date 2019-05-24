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
	"errors"
	"testing"
)

// all the test cases' expect results are calculated from mycat rule function

func Test_MycatPartitionModShard_FindForKey(t *testing.T) {
	tests := []struct {
		ShardNum    int
		Key         interface{}
		ExpectSlice int
		Error       error
	}{
		{1, 0, 0, nil},
		{1, 1, 0, nil},
		{1, 2, 0, nil},
		{1, -1, 0, nil},

		{2, 0, 0, nil},
		{2, 1, 1, nil},
		{2, 2, 0, nil},
		{2, 3, 1, nil},
		{2, -1, 1, nil},
		{2, -2, 0, nil},

		{3, 0, 0, nil},
		{3, 1, 1, nil},
		{3, 2, 2, nil},
		{3, 3, 0, nil},
		{3, 4, 1, nil},
		{3, -1, 1, nil},
		{3, -2, 2, nil},
		{3, -3, 0, nil},
		{3, -4, 1, nil},
	}
	for _, test := range tests {
		shard := NewMycatPartitionModShard(test.ShardNum)
		actualSlice, err := shard.FindForKey(test.Key)
		if err != nil {
			t.Errorf("case: %v, err: %s", test, err.Error())
		}
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %v, expect: %d, actual: %d", test, test.ExpectSlice, actualSlice)
		}
	}
}

var (
	err1 = errors.New("error, check your scope & scopeLength definition")
	err2 = errors.New("error, check your partitionScope definition")
)

func Test_MycatPartitionLongShard_Init(t *testing.T) {
	tests := []struct {
		shardNum        int
		partitionCount  string
		partitionLength string
		err             error
	}{
		{1, "1", "1024", nil},
		{2, "2", "512", nil},
		{4, "4", "256", nil},
		{8, "8", "128", nil},
		{2, "1,1", "512,512", nil},
		{3, "1,2", "512,256", nil},
		{5, "1,4", "512,128", nil},
		{11, "1,2,8", "256,128,64", nil},
		{10, "1,4,5", "512,128", err1},
		{3, "1,2", "128,512", err2},
		{3, "1,2", "128,128", err2},
	}
	for _, test := range tests {
		shard := NewMycatPartitionLongShard(test.shardNum, test.partitionCount, test.partitionLength)
		actualErr := shard.Init()
		if test.err == nil {
			if actualErr != nil {
				t.Errorf("not equal, case: %+v, expect: %v, actual: %v", test, test.err, actualErr)
			}
		} else {
			if actualErr == nil || test.err.Error() != actualErr.Error() {
				t.Errorf("not equal, case: %+v, expect: %v, actual: %v", test, test.err, actualErr)
			}
		}
	}
}

func Test_MycatPartitionLongShard_FindForKey_BalanceLength_1(t *testing.T) {
	partitionCount := "4"
	partitionLength := "256"
	shard := NewMycatPartitionLongShard(4, partitionCount, partitionLength)
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"-1", 3},
		{"0", 0},
		{"1", 0},
		{"2", 0},
		{"3", 0},
		{"64", 0},
		{"128", 0},
		{"256", 1},
		{"512", 2},
		{"768", 3},
		{"1024", 0},
		{"1025", 0},
		{"1280", 1},
		{"2048", 0},
		{"-9223372036854775808", 0},
		{"9223372036854775807", 3},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_MycatPartitionLongShard_FindForKey_BalanceLength_2(t *testing.T) {
	partitionCount := "2,2"
	partitionLength := "256,256"
	shard := NewMycatPartitionLongShard(4, partitionCount, partitionLength)
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"-1", 3},
		{"0", 0},
		{"1", 0},
		{"2", 0},
		{"3", 0},
		{"64", 0},
		{"128", 0},
		{"256", 1},
		{"512", 2},
		{"768", 3},
		{"1024", 0},
		{"1025", 0},
		{"1280", 1},
		{"2048", 0},
		{"-9223372036854775808", 0},
		{"9223372036854775807", 3},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_MycatPartitionLongShard_FindForKey_InalanceLength_3(t *testing.T) {
	partitionCount := "1,1,4"
	partitionLength := "512,256,64"
	shard := NewMycatPartitionLongShard(6, partitionCount, partitionLength)
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{uint64(0), 0},
		{[]byte("0"), 0},
		{"-1", 5},
		{"0", 0},
		{"1", 0},
		{"2", 0},
		{"3", 0},
		{"64", 0},
		{"128", 0},
		{"256", 0},
		{"512", 1},
		{"768", 2},
		{"1024", 0},
		{"1025", 0},
		{"1280", 0},
		{"2048", 0},
		{"-9223372036854775808", 0},
		{"9223372036854775807", 5},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

// take care: count=0, mycat will occur java.util.NoSuchElementException, but here return the same result with count=1
// besides, count=0 is a invalid conf, should be checked before call the Init() method
func Test_MycatPartitionMurmurHashShard_Seed0_Count1(t *testing.T) {
	seed := "0"
	virtualBucketTimes := "160"
	count := 1
	shard, err := NewMycatPartitionMurmurHashShard(seed, virtualBucketTimes, count)
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	err = shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"", 0},
		{uint64(0), 0},
		{int64(0), 0},
		{[]byte("gaea"), 0},
		{"hello, world", 0},
		{"你好, 中国", 0},
		{"?!)_FFSD", 0},
		{"-50", 0},
		{"-49", 0},
		{"-48", 0},
		{"-47", 0},
		{"-46", 0},
		{"-45", 0},
		{"-44", 0},
		{"-43", 0},
		{"-42", 0},
		{"-41", 0},
		{"-40", 0},
		{"-39", 0},
		{"-38", 0},
		{"-37", 0},
		{"-36", 0},
		{"-35", 0},
		{"-34", 0},
		{"-33", 0},
		{"-32", 0},
		{"-31", 0},
		{"-30", 0},
		{"-29", 0},
		{"-28", 0},
		{"-27", 0},
		{"-26", 0},
		{"-25", 0},
		{"-24", 0},
		{"-23", 0},
		{"-22", 0},
		{"-21", 0},
		{"-20", 0},
		{"-19", 0},
		{"-18", 0},
		{"-17", 0},
		{"-16", 0},
		{"-15", 0},
		{"-14", 0},
		{"-13", 0},
		{"-12", 0},
		{"-11", 0},
		{"-10", 0},
		{"-9", 0},
		{"-8", 0},
		{"-7", 0},
		{"-6", 0},
		{"-5", 0},
		{"-4", 0},
		{"-3", 0},
		{"-2", 0},
		{"-1", 0},
		{"0", 0},
		{"1", 0},
		{"2", 0},
		{"3", 0},
		{"4", 0},
		{"5", 0},
		{"6", 0},
		{"7", 0},
		{"8", 0},
		{"9", 0},
		{"10", 0},
		{"11", 0},
		{"12", 0},
		{"13", 0},
		{"14", 0},
		{"15", 0},
		{"16", 0},
		{"17", 0},
		{"18", 0},
		{"19", 0},
		{"20", 0},
		{"21", 0},
		{"22", 0},
		{"23", 0},
		{"24", 0},
		{"25", 0},
		{"26", 0},
		{"27", 0},
		{"28", 0},
		{"29", 0},
		{"30", 0},
		{"31", 0},
		{"32", 0},
		{"33", 0},
		{"34", 0},
		{"35", 0},
		{"36", 0},
		{"37", 0},
		{"38", 0},
		{"39", 0},
		{"40", 0},
		{"41", 0},
		{"42", 0},
		{"43", 0},
		{"44", 0},
		{"45", 0},
		{"46", 0},
		{"47", 0},
		{"48", 0},
		{"49", 0},
		{"50", 0},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_MycatPartitionMurmurHashShard_Seed0_Count2(t *testing.T) {
	seed := "0"
	virtualBucketTimes := "160"
	count := 2
	shard, err := NewMycatPartitionMurmurHashShard(seed, virtualBucketTimes, count)
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	err = shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"", 0},
		{"hello, world", 0},
		{"你好, 中国", 0},
		{"?!)_FFSD", 1},
		{"-50", 0},
		{"-49", 0},
		{"-48", 0},
		{"-47", 0},
		{"-46", 1},
		{"-45", 1},
		{"-44", 0},
		{"-43", 1},
		{"-42", 0},
		{"-41", 0},
		{"-40", 0},
		{"-39", 1},
		{"-38", 0},
		{"-37", 0},
		{"-36", 1},
		{"-35", 0},
		{"-34", 1},
		{"-33", 0},
		{"-32", 0},
		{"-31", 0},
		{"-30", 1},
		{"-29", 1},
		{"-28", 1},
		{"-27", 1},
		{"-26", 0},
		{"-25", 1},
		{"-24", 0},
		{"-23", 0},
		{"-22", 0},
		{"-21", 0},
		{"-20", 1},
		{"-19", 1},
		{"-18", 1},
		{"-17", 1},
		{"-16", 0},
		{"-15", 0},
		{"-14", 1},
		{"-13", 1},
		{"-12", 0},
		{"-11", 0},
		{"-10", 0},
		{"-9", 1},
		{"-8", 1},
		{"-7", 0},
		{"-6", 0},
		{"-5", 0},
		{"-4", 0},
		{"-3", 1},
		{"-2", 0},
		{"-1", 0},
		{"0", 1},
		{"1", 1},
		{"2", 1},
		{"3", 1},
		{"4", 1},
		{"5", 0},
		{"6", 0},
		{"7", 1},
		{"8", 0},
		{"9", 1},
		{"10", 1},
		{"11", 1},
		{"12", 1},
		{"13", 0},
		{"14", 0},
		{"15", 1},
		{"16", 0},
		{"17", 0},
		{"18", 0},
		{"19", 1},
		{"20", 1},
		{"21", 0},
		{"22", 0},
		{"23", 1},
		{"24", 1},
		{"25", 0},
		{"26", 0},
		{"27", 1},
		{"28", 0},
		{"29", 1},
		{"30", 1},
		{"31", 1},
		{"32", 1},
		{"33", 1},
		{"34", 0},
		{"35", 0},
		{"36", 0},
		{"37", 1},
		{"38", 0},
		{"39", 0},
		{"40", 0},
		{"41", 1},
		{"42", 1},
		{"43", 1},
		{"44", 1},
		{"45", 1},
		{"46", 1},
		{"47", 1},
		{"48", 1},
		{"49", 0},
		{"50", 1},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_MycatPartitionMurmurHashShard_IntKey_Seed0_Count2(t *testing.T) {
	seed := "0"
	virtualBucketTimes := "160"
	count := 2
	shard, err := NewMycatPartitionMurmurHashShard(seed, virtualBucketTimes, count)
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	err = shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{-50, 0},
		{-49, 0},
		{-48, 0},
		{-47, 0},
		{-46, 1},
		{-45, 1},
		{-44, 0},
		{-43, 1},
		{-42, 0},
		{-41, 0},
		{-40, 0},
		{-39, 1},
		{-38, 0},
		{-37, 0},
		{-36, 1},
		{-35, 0},
		{-34, 1},
		{-33, 0},
		{-32, 0},
		{-31, 0},
		{-30, 1},
		{-29, 1},
		{-28, 1},
		{-27, 1},
		{-26, 0},
		{-25, 1},
		{-24, 0},
		{-23, 0},
		{-22, 0},
		{-21, 0},
		{-20, 1},
		{-19, 1},
		{-18, 1},
		{-17, 1},
		{-16, 0},
		{-15, 0},
		{-14, 1},
		{-13, 1},
		{-12, 0},
		{-11, 0},
		{-10, 0},
		{-9, 1},
		{-8, 1},
		{-7, 0},
		{-6, 0},
		{-5, 0},
		{-4, 0},
		{-3, 1},
		{-2, 0},
		{-1, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
		{5, 0},
		{6, 0},
		{7, 1},
		{8, 0},
		{9, 1},
		{10, 1},
		{11, 1},
		{12, 1},
		{13, 0},
		{14, 0},
		{15, 1},
		{16, 0},
		{17, 0},
		{18, 0},
		{19, 1},
		{20, 1},
		{21, 0},
		{22, 0},
		{23, 1},
		{24, 1},
		{25, 0},
		{26, 0},
		{27, 1},
		{28, 0},
		{29, 1},
		{30, 1},
		{31, 1},
		{32, 1},
		{33, 1},
		{34, 0},
		{35, 0},
		{36, 0},
		{37, 1},
		{38, 0},
		{39, 0},
		{40, 0},
		{41, 1},
		{42, 1},
		{43, 1},
		{44, 1},
		{45, 1},
		{46, 1},
		{47, 1},
		{48, 1},
		{49, 0},
		{50, 1},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_MycatPartitionMurmurHashShard_Seed1_Count4(t *testing.T) {
	seed := "1"
	virtualBucketTimes := "160"
	count := 4
	shard, err := NewMycatPartitionMurmurHashShard(seed, virtualBucketTimes, count)
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	err = shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"", 2},
		{"hello, world", 1},
		{"你好, 中国", 0},
		{"?!)_FFSD", 1},
		{"-50", 1},
		{"-49", 0},
		{"-48", 0},
		{"-47", 2},
		{"-46", 3},
		{"-45", 3},
		{"-44", 3},
		{"-43", 1},
		{"-42", 3},
		{"-41", 3},
		{"-40", 1},
		{"-39", 3},
		{"-38", 3},
		{"-37", 2},
		{"-36", 1},
		{"-35", 0},
		{"-34", 1},
		{"-33", 0},
		{"-32", 1},
		{"-31", 3},
		{"-30", 0},
		{"-29", 1},
		{"-28", 0},
		{"-27", 1},
		{"-26", 1},
		{"-25", 2},
		{"-24", 0},
		{"-23", 0},
		{"-22", 0},
		{"-21", 1},
		{"-20", 1},
		{"-19", 1},
		{"-18", 2},
		{"-17", 3},
		{"-16", 2},
		{"-15", 0},
		{"-14", 0},
		{"-13", 2},
		{"-12", 0},
		{"-11", 0},
		{"-10", 3},
		{"-9", 2},
		{"-8", 3},
		{"-7", 0},
		{"-6", 3},
		{"-5", 3},
		{"-4", 1},
		{"-3", 2},
		{"-2", 0},
		{"-1", 3},
		{"0", 0},
		{"1", 2},
		{"2", 2},
		{"3", 2},
		{"4", 3},
		{"5", 3},
		{"6", 1},
		{"7", 3},
		{"8", 2},
		{"9", 2},
		{"10", 0},
		{"11", 1},
		{"12", 0},
		{"13", 1},
		{"14", 0},
		{"15", 2},
		{"16", 2},
		{"17", 1},
		{"18", 2},
		{"19", 2},
		{"20", 0},
		{"21", 0},
		{"22", 0},
		{"23", 1},
		{"24", 2},
		{"25", 2},
		{"26", 3},
		{"27", 3},
		{"28", 2},
		{"29", 1},
		{"30", 1},
		{"31", 1},
		{"32", 0},
		{"33", 3},
		{"34", 2},
		{"35", 0},
		{"36", 1},
		{"37", 2},
		{"38", 1},
		{"39", 0},
		{"40", 2},
		{"41", 0},
		{"42", 2},
		{"43", 0},
		{"44", 2},
		{"45", 0},
		{"46", 3},
		{"47", 0},
		{"48", 3},
		{"49", 2},
		{"50", 0},
	}
	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_PartitionPaddingMod_Default_Success(t *testing.T) {
	shard := newDefaultMycatPartitionPaddingModShard()
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"1000000000000000", 0},
		{"1000000000000001", 1},
		{"1000000000000002", 0},
		{"1000000000000003", 1},
		{"1000000000000004", 0},
		{"1000000000000005", 1},
		{"1000000000000006", 0},
		{"1000000000000007", 1},
		{"1000000000000008", 0},
		{"1000000000000009", 1},
		{"1000000000000010", 0},
		{"1000000000000011", 1},
		{"1000000000000012", 0},
		{"1000000000000013", 1},
		{"1000000000000014", 0},
		{"1000000000000015", 1},
		{"1000000000000016", 0},
		{"1000000000000017", 1},
		{"1000000000000018", 0},
		{"1000000000000019", 1},
		{"1000000000000020", 0},
		{"1000000000000021", 1},
		{"1000000000000022", 0},
		{"1000000000000023", 1},
		{"1000000000000024", 0},
		{"1000000000000025", 1},
		{"1000000000000026", 0},
		{"1000000000000027", 1},
		{"1000000000000028", 0},
		{"1000000000000029", 1},
		{"1000000000000030", 0},
		{"1000000000000031", 1},
		{"1000000000000032", 0},
		{"1000000000000033", 1},
		{"1000000000000034", 0},
		{"1000000000000035", 1},
		{"1000000000000036", 0},
		{"1000000000000037", 1},
		{"1000000000000038", 0},
		{"1000000000000039", 1},
		{"1000000000000040", 0},
		{"1000000000000041", 1},
		{"1000000000000042", 0},
		{"1000000000000043", 1},
		{"1000000000000044", 0},
		{"1000000000000045", 1},
		{"1000000000000046", 0},
		{"1000000000000047", 1},
		{"1000000000000048", 0},
		{"1000000000000049", 1},
		{"1000000000000050", 0},
		{"1000000000000051", 1},
		{"1000000000000052", 0},
		{"1000000000000053", 1},
		{"1000000000000054", 0},
		{"1000000000000055", 1},
		{"1000000000000056", 0},
		{"1000000000000057", 1},
		{"1000000000000058", 0},
		{"1000000000000059", 1},
		{"1000000000000060", 0},
		{"1000000000000061", 1},
		{"1000000000000062", 0},
		{"1000000000000063", 1},
		{"1000000000000064", 0},
		{"1000000000000065", 1},
		{"1000000000000066", 0},
		{"1000000000000067", 1},
		{"1000000000000068", 0},
		{"1000000000000069", 1},
		{"1000000000000070", 0},
		{"1000000000000071", 1},
		{"1000000000000072", 0},
		{"1000000000000073", 1},
		{"1000000000000074", 0},
		{"1000000000000075", 1},
		{"1000000000000076", 0},
		{"1000000000000077", 1},
		{"1000000000000078", 0},
		{"1000000000000079", 1},
		{"1000000000000080", 0},
		{"1000000000000081", 1},
		{"1000000000000082", 0},
		{"1000000000000083", 1},
		{"1000000000000084", 0},
		{"1000000000000085", 1},
		{"1000000000000086", 0},
		{"1000000000000087", 1},
		{"1000000000000088", 0},
		{"1000000000000089", 1},
		{"1000000000000090", 0},
		{"1000000000000091", 1},
		{"1000000000000092", 0},
		{"1000000000000093", 1},
		{"1000000000000094", 0},
		{"1000000000000095", 1},
		{"1000000000000096", 0},
		{"1000000000000097", 1},
		{"1000000000000098", 0},
		{"1000000000000099", 1},
		{"1000000000000100", 0},
	}

	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_PartitionPaddingMod_Default_KeyIsNotNumber(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("expect panic, err: %s", err)
		}
	}()

	shard := newDefaultMycatPartitionPaddingModShard()
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}

	_, err = shard.FindForKey("hello")
}

func Test_PartitionPaddingMod_CustomParam(t *testing.T) {
	shard := newMycatPartitionPaddingModShard(1, 18, 14, 16, 32)
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"1000000000000000", 0},
		{"1000000000000001", 1},
		{"1000000000000002", 2},
		{"1000000000000003", 3},
		{"1000000000000004", 4},
		{"1000000000000005", 5},
		{"1000000000000006", 6},
		{"1000000000000007", 7},
		{"1000000000000008", 8},
		{"1000000000000009", 9},
		{"1000000000000010", 10},
		{"1000000000000011", 11},
		{"1000000000000012", 12},
		{"1000000000000013", 13},
		{"1000000000000014", 14},
		{"1000000000000015", 15},
		{"1000000000000016", 16},
		{"1000000000000017", 17},
		{"1000000000000018", 18},
		{"1000000000000019", 19},
		{"1000000000000020", 20},
		{"1000000000000021", 21},
		{"1000000000000022", 22},
		{"1000000000000023", 23},
		{"1000000000000024", 24},
		{"1000000000000025", 25},
		{"1000000000000026", 26},
		{"1000000000000027", 27},
		{"1000000000000028", 28},
		{"1000000000000029", 29},
		{"1000000000000030", 30},
		{"1000000000000031", 31},
		{"1000000000000032", 0},
		{"1000000000000033", 1},
		{"1000000000000034", 2},
		{"1000000000000035", 3},
		{"1000000000000036", 4},
		{"1000000000000037", 5},
		{"1000000000000038", 6},
		{"1000000000000039", 7},
		{"1000000000000040", 8},
		{"1000000000000041", 9},
		{"1000000000000042", 10},
		{"1000000000000043", 11},
		{"1000000000000044", 12},
		{"1000000000000045", 13},
		{"1000000000000046", 14},
		{"1000000000000047", 15},
		{"1000000000000048", 16},
		{"1000000000000049", 17},
		{"1000000000000050", 18},
		{"1000000000000051", 19},
		{"1000000000000052", 20},
		{"1000000000000053", 21},
		{"1000000000000054", 22},
		{"1000000000000055", 23},
		{"1000000000000056", 24},
		{"1000000000000057", 25},
		{"1000000000000058", 26},
		{"1000000000000059", 27},
		{"1000000000000060", 28},
		{"1000000000000061", 29},
		{"1000000000000062", 30},
		{"1000000000000063", 31},
		{"1000000000000064", 0},
		{"1000000000000065", 1},
		{"1000000000000066", 2},
		{"1000000000000067", 3},
		{"1000000000000068", 4},
		{"1000000000000069", 5},
		{"1000000000000070", 6},
		{"1000000000000071", 7},
		{"1000000000000072", 8},
		{"1000000000000073", 9},
		{"1000000000000074", 10},
		{"1000000000000075", 11},
		{"1000000000000076", 12},
		{"1000000000000077", 13},
		{"1000000000000078", 14},
		{"1000000000000079", 15},
		{"1000000000000080", 16},
		{"1000000000000081", 17},
		{"1000000000000082", 18},
		{"1000000000000083", 19},
		{"1000000000000084", 20},
		{"1000000000000085", 21},
		{"1000000000000086", 22},
		{"1000000000000087", 23},
		{"1000000000000088", 24},
		{"1000000000000089", 25},
		{"1000000000000090", 26},
		{"1000000000000091", 27},
		{"1000000000000092", 28},
		{"1000000000000093", 29},
		{"1000000000000094", 30},
		{"1000000000000095", 31},
		{"1000000000000096", 0},
		{"1000000000000097", 1},
		{"1000000000000098", 2},
		{"1000000000000099", 3},
		{"1000000000000100", 0},
		{"1000000000000101", 1},
		{"1000000000000102", 2},
		{"1000000000000103", 3},
		{"1000000000000104", 4},
		{"1000000000000105", 5},
		{"1000000000000106", 6},
		{"1000000000000107", 7},
		{"1000000000000108", 8},
		{"1000000000000109", 9},
		{"1000000000000110", 10},
		{"1000000000000111", 11},
		{"1000000000000112", 12},
		{"1000000000000113", 13},
		{"1000000000000114", 14},
		{"1000000000000115", 15},
		{"1000000000000116", 16},
		{"1000000000000117", 17},
		{"1000000000000118", 18},
		{"1000000000000119", 19},
		{"1000000000000120", 20},
		{"1000000000000121", 21},
		{"1000000000000122", 22},
		{"1000000000000123", 23},
		{"1000000000000124", 24},
		{"1000000000000125", 25},
		{"1000000000000126", 26},
		{"1000000000000127", 27},
		{"1000000000000128", 28},
		{"1000000000000129", 29},
		{"1000000000000130", 30},
		{"1000000000000131", 31},
		{"1000000000000132", 0},
		{"1000000000000133", 1},
		{"1000000000000134", 2},
		{"1000000000000135", 3},
		{"1000000000000136", 4},
		{"1000000000000137", 5},
		{"1000000000000138", 6},
		{"1000000000000139", 7},
		{"1000000000000140", 8},
		{"1000000000000141", 9},
		{"1000000000000142", 10},
		{"1000000000000143", 11},
		{"1000000000000144", 12},
		{"1000000000000145", 13},
		{"1000000000000146", 14},
		{"1000000000000147", 15},
		{"1000000000000148", 16},
		{"1000000000000149", 17},
		{"1000000000000150", 18},
		{"1000000000000151", 19},
		{"1000000000000152", 20},
		{"1000000000000153", 21},
		{"1000000000000154", 22},
		{"1000000000000155", 23},
		{"1000000000000156", 24},
		{"1000000000000157", 25},
		{"1000000000000158", 26},
		{"1000000000000159", 27},
		{"1000000000000160", 28},
		{"1000000000000161", 29},
		{"1000000000000162", 30},
		{"1000000000000163", 31},
		{"1000000000000164", 0},
		{"1000000000000165", 1},
		{"1000000000000166", 2},
		{"1000000000000167", 3},
		{"1000000000000168", 4},
		{"1000000000000169", 5},
		{"1000000000000170", 6},
		{"1000000000000171", 7},
		{"1000000000000172", 8},
		{"1000000000000173", 9},
		{"1000000000000174", 10},
		{"1000000000000175", 11},
		{"1000000000000176", 12},
		{"1000000000000177", 13},
		{"1000000000000178", 14},
		{"1000000000000179", 15},
		{"1000000000000180", 16},
		{"1000000000000181", 17},
		{"1000000000000182", 18},
		{"1000000000000183", 19},
		{"1000000000000184", 20},
		{"1000000000000185", 21},
		{"1000000000000186", 22},
		{"1000000000000187", 23},
		{"1000000000000188", 24},
		{"1000000000000189", 25},
		{"1000000000000190", 26},
		{"1000000000000191", 27},
		{"1000000000000192", 28},
		{"1000000000000193", 29},
		{"1000000000000194", 30},
		{"1000000000000195", 31},
		{"1000000000000196", 0},
		{"1000000000000197", 1},
		{"1000000000000198", 2},
		{"1000000000000199", 3},
		{"1000000000000200", 0},
	}

	for _, test := range tests {
		actualSlice, _ := shard.FindForKey(test.Key)
		if actualSlice != test.ExpectSlice {
			t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
		}
	}
}

func Test_PartitionPaddingMod_LeftEnd_KeyIsNegative(t *testing.T) {
	shard := newMycatPartitionPaddingModShard(PaddingModLeftEnd, PaddingModDefaultPadLength, PaddingModDefaultModBegin, PaddingModDefaultModEnd, PaddingModDefaultMod)
	err := shard.Init()
	if err != nil {
		t.Errorf("init error: %s", err.Error())
	}
	_, err = shard.FindForKey("-1")
	if err == nil {
		t.Errorf("expect error but return nil")
	}
}

func Test_PartitionPaddingMod_GetNewInstance(t *testing.T) {
	tests := []struct {
		PadFromStr      string
		PadLengthStr    string
		ModBeginStr     string
		ModEndStr       string
		ExpectPadFrom   int
		ExpectPadLength int
		ExpectModBegin  int
		ExpectModEnd    int
		Error           error
	}{
		{"1", "18", "10", "16", 1, 18, 10, 16, nil},
	}

	for _, test := range tests {
		shard, err := GetMycatPartitionPaddingModShard(test.PadFromStr, test.PadLengthStr, test.ModBeginStr, test.ModEndStr, 2)
		if (err == nil && test.Error != nil) || (err != nil && test.Error == nil) {
			t.Fatalf("err is not equal, expect: %v, actual: %v", test.Error, err)
		}
		if test.Error != nil {
			if err.Error() != test.Error.Error() {
				t.Fatalf("err is not equal, expect: %v, actual: %v", test.Error, err)
			} else {
				continue
			}
		}
		if test.ExpectPadFrom != shard.padFrom {
			t.Fatalf("padFrom not equal, expect: %d, actual: %d", test.ExpectPadFrom, shard.padFrom)
		}
		if test.ExpectPadLength != shard.padLength {
			t.Fatalf("padLength not equal, expect: %d, actual: %d", test.ExpectPadLength, shard.padLength)
		}
		if test.ExpectModBegin != shard.modBegin {
			t.Fatalf("modBegin not equal, expect: %d, actual: %d", test.ExpectModBegin, shard.modBegin)
		}
		if test.ExpectModEnd != shard.modEnd {
			t.Fatalf("modEnd not equal, expect: %d, actual: %d", test.ExpectModEnd, shard.modEnd)
		}
	}
}

func Test_PartitionPaddingMod_Init(t *testing.T) {
	tests := []struct {
		PadFrom   int
		PadLength int
		ModBegin  int
		ModEnd    int
		Mod       int
		Error     error
	}{
		// padFrom
		{-1, 1, 0, 1, 2, errors.New("invalid padding mod mode: -1")},
		{2, 1, 0, 1, 2, errors.New("invalid padding mod mode: 2")},
		{0, 1, 0, 1, 2, nil},
		{1, 1, 0, 1, 2, nil},
		// padLength
		{0, 0, 0, 1, 2, errors.New("invalid padding mod padLength: 0")},
		{0, -1, 0, 1, 2, errors.New("invalid padding mod padLength: -1")},
		// modBegin and modEnd
		{0, 1, -1, 1, 2, errors.New("invalid padding modBegin or modEnd: -1, 1")},
		{0, 1, 0, 0, 2, errors.New("invalid padding modBegin or modEnd: 0, 0")},
		{0, 1, 1, 1, 2, errors.New("invalid padding modBegin or modEnd: 1, 1")},
		{0, 1, 2, 1, 2, errors.New("invalid padding modBegin or modEnd: 2, 1")},
		// padLength and modLength
		{0, 1, 0, 2, 2, errors.New("invalid padding mod, padLength is less than modBegin - modEnd: 1, 0, 2")},
		{0, 2, 0, 2, 2, nil},
		// mod
		{0, 2, 0, 2, 1, errors.New("invalid padding mod number: 1")},
		{0, 2, 0, 2, 0, errors.New("invalid padding mod number: 0")},
		{0, 2, 0, 2, 5, nil},
	}

	for _, test := range tests {
		shard := newMycatPartitionPaddingModShard(test.PadFrom, test.PadLength, test.ModBegin, test.ModEnd, test.Mod)
		err := shard.Init()
		if (err == nil && test.Error != nil) || (err != nil && test.Error == nil) {
			t.Fatalf("err is not equal, expect: %v, actual: %v", test.Error, err)
		}
		if test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("err is not equal, expect: %v, actual: %v", test.Error, err)
		}
	}
}

func TestParseHashSliceStartEnd(t *testing.T) {
	tests := []struct {
		hashSlice string
		start     int
		end       int
		hasErr    bool
	}{
		{hashSlice: "2", start: 0, end: 2},
		{hashSlice: "1:2", start: 1, end: 2},
		{hashSlice: "1:", start: 1, end: 0},
		{hashSlice: "-1:", start: -1, end: 0},
		{hashSlice: ":-1", start: 0, end: -1},
		{hashSlice: ":", start: 0, end: 0},
		{hashSlice: "", hasErr: true},
		{hashSlice: "0:1:2", hasErr: true},
		{hashSlice: "a:1", hasErr: true},
		{hashSlice: "a", hasErr: true},
		{hashSlice: "a:", hasErr: true},
		{hashSlice: ":a", hasErr: true},
		{hashSlice: "1:a", hasErr: true},
	}
	for _, test := range tests {
		t.Run(test.hashSlice, func(t *testing.T) {
			start, end, err := parseHashSliceStartEnd(test.hashSlice)
			if err != nil {
				if test.hasErr {
					t.Logf("expected error: %v", err)
					return
				} else {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if test.start != start {
				t.Errorf("start not equal, expect: %d, actual: %d", test.start, start)
			}
			if test.end != end {
				t.Errorf("end not equal, expect: %d, actual: %d", test.end, end)
			}
		})
	}
}

func Test_MycatPartitionStringShard_32(t *testing.T) {
	shard := NewMycatPartitionStringShard(64, "64", "16", "32")
	err := shard.Init()
	if err != nil {
		t.Fatalf("init error: %s", err.Error())
	}
	tests := []struct {
		Key         interface{}
		ExpectSlice int
	}{
		{"", 0},
		{"hello, world", 24},
		{"你好, 中国", 40},
		{"?!)_FFSD", 58},
		{"ddda;kjelwr", 63},
		{"-120123012", 4},
		{"1029421093", 17},
		{"-50", 56},
		{"-49", 55},
		{"-48", 55},
		{"-47", 55},
		{"-46", 54},
		{"-45", 54},
		{"-44", 54},
		{"-43", 54},
		{"-42", 54},
		{"-41", 54},
		{"-40", 54},
		{"-39", 53},
		{"-38", 53},
		{"-37", 53},
		{"-36", 53},
		{"-35", 52},
		{"-34", 52},
		{"-33", 52},
		{"-32", 52},
		{"-31", 52},
		{"-30", 52},
		{"-29", 51},
		{"-28", 51},
		{"-27", 51},
		{"-26", 51},
		{"-25", 51},
		{"-24", 50},
		{"-23", 50},
		{"-22", 50},
		{"-21", 50},
		{"-20", 50},
		{"-19", 49},
		{"-18", 49},
		{"-17", 49},
		{"-16", 49},
		{"-15", 49},
		{"-14", 49},
		{"-13", 48},
		{"-12", 48},
		{"-11", 48},
		{"-10", 48},
		{"-9", 26},
		{"-8", 26},
		{"-7", 26},
		{"-6", 26},
		{"-5", 26},
		{"-4", 26},
		{"-3", 26},
		{"-2", 26},
		{"-1", 26},
		{"0", 3},
		{"1", 3},
		{"2", 3},
		{"3", 3},
		{"4", 3},
		{"5", 3},
		{"6", 3},
		{"7", 3},
		{"8", 3},
		{"9", 3},
		{"10", 33},
		{"11", 34},
		{"12", 34},
		{"13", 34},
		{"14", 34},
		{"15", 34},
		{"16", 34},
		{"17", 34},
		{"18", 34},
		{"19", 34},
		{"20", 35},
		{"21", 35},
		{"22", 36},
		{"23", 36},
		{"24", 36},
		{"25", 36},
		{"26", 36},
		{"27", 36},
		{"28", 36},
		{"29", 36},
		{"30", 37},
		{"31", 37},
		{"32", 37},
		{"33", 38},
		{"34", 38},
		{"35", 38},
		{"36", 38},
		{"37", 38},
		{"38", 38},
		{"39", 38},
		{"40", 39},
		{"41", 39},
		{"42", 39},
		{"43", 39},
		{"44", 40},
		{"45", 40},
		{"46", 40},
		{"47", 40},
		{"48", 40},
		{"49", 40},
		{"50", 41},
	}
	for _, test := range tests {
		testKey := test.Key.(string)
		t.Run(testKey, func(t *testing.T) {
			actualSlice, _ := shard.FindForKey(test.Key)
			if actualSlice != test.ExpectSlice {
				t.Errorf("not equal, case: %+v, expect: %d, actual: %d", test, actualSlice, test.ExpectSlice)
			}
		})
	}
}
