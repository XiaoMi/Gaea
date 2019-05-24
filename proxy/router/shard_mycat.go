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
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/hack"
)

// MycatPartitionModShard mycat route algorithm: PartitionByMod
// take care: in Mycat, the key is parsed to a BigInteger, not int64.
type MycatPartitionModShard struct {
	ShardNum int
}

// NewMycatPartitionModShard constructor of MycatPartitionModShard
func NewMycatPartitionModShard(shardNum int) *MycatPartitionModShard {
	return &MycatPartitionModShard{
		ShardNum: shardNum,
	}
}

// FindForKey return result of calculated key
func (m *MycatPartitionModShard) FindForKey(key interface{}) (int, error) {
	h := hack.Abs(NumValue(key))
	return int(h % int64(m.ShardNum)), nil
}

const (
	// PartitionLength length of partition
	PartitionLength = 1024
	andValue        = PartitionLength - 1
)

// MycatPartitionLongShard mycat route algorithm: PartitionByLong
type MycatPartitionLongShard struct {
	shardNum  int
	countStr  string
	lengthStr string
	count     []int
	length    []int
	segment   []int
}

// NewMycatPartitionLongShard constructor of MycatPartitionLongShard
func NewMycatPartitionLongShard(shardNum int, partitionCount, partitionLength string) *MycatPartitionLongShard {
	return &MycatPartitionLongShard{
		shardNum:  shardNum,
		countStr:  partitionCount,
		lengthStr: partitionLength,
		count:     make([]int, 0, 10),
		length:    make([]int, 0, 10),
		segment:   make([]int, 0, PartitionLength),
	}
}

// Init init MycatPartitionLongShard
func (m *MycatPartitionLongShard) Init() error {
	countList, err := toIntArray(m.countStr)
	if err != nil {
		return err
	}

	lengthList, err := toIntArray(m.lengthStr)
	if err != nil {
		return err
	}

	countSize := len(countList)
	lengthSize := len(lengthList)
	if countSize != lengthSize {
		return errors.New("error, check your scope & scopeLength definition")
	}

	segmentLength := 0
	for i := 0; i < countSize; i++ {
		segmentLength += countList[i]
	}

	if segmentLength != m.shardNum {
		return errors.New("segmentLength is not equal to shardNum")
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
		return errors.New("error, check your partitionScope definition")
	}

	segment := make([]int, PartitionLength)
	for i := 1; i < len(ai); i++ {
		for j := ai[i-1]; j < ai[i]; j++ {
			segment[j] = i - 1
		}
	}

	m.count = countList
	m.length = lengthList
	m.segment = segment
	return nil
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

// FindForKey return MycatPartitionLongShard calculated result
func (m *MycatPartitionLongShard) FindForKey(key interface{}) (int, error) {
	h := NumValue(key)
	return m.segment[int(h)&andValue], nil
}

// MycatPartitionStringShard mycat route algorithm: PartitionByString
type MycatPartitionStringShard struct {
	*MycatPartitionLongShard
	hashSliceStr   string
	hashSliceStart int
	hashSliceEnd   int
}

// NewMycatPartitionStringShard constructor of MycatPartitionStringShard
func NewMycatPartitionStringShard(shardNum int, partitionCount, partitionLength string, hashSliceStr string) *MycatPartitionStringShard {
	longShard := NewMycatPartitionLongShard(shardNum, partitionCount, partitionLength)
	stringShard := &MycatPartitionStringShard{
		MycatPartitionLongShard: longShard,
		hashSliceStr:            hashSliceStr,
	}
	return stringShard
}

// Init init NewMycatPartitionStringShard
func (m *MycatPartitionStringShard) Init() error {
	if err := m.MycatPartitionLongShard.Init(); err != nil {
		return err
	}

	start, end, err := parseHashSliceStartEnd(m.hashSliceStr)
	if err != nil {
		return err
	}
	m.hashSliceStart = start
	m.hashSliceEnd = end
	return nil
}

/**
 * "2" -&gt; (0,2)<br/>
 * "1:2" -&gt; (1,2)<br/>
 * "1:" -&gt; (1,0)<br/>
 * "-1:" -&gt; (-1,0)<br/>
 * ":-1" -&gt; (0,-1)<br/>
 * ":" -&gt; (0,0)<br/>
 */
// copied from mycat
func parseHashSliceStartEnd(hashSliceStr string) (int, int, error) {
	hashSliceStr = strings.TrimSpace(hashSliceStr)
	strs := strings.Split(hashSliceStr, ":")

	if len(strs) == 1 {
		v, err := strconv.Atoi(strs[0])
		if err != nil {
			return -1, -1, err
		}
		if v >= 0 {
			return 0, v, nil
		}
		return v, 0, nil
	}

	if len(strs) == 2 {
		start, err := parseHashSliceValue(strs[0])
		if err != nil {
			return -1, -1, fmt.Errorf("parse hash slice start error: %v", err)
		}
		end, err := parseHashSliceValue(strs[1])
		if err != nil {
			return -1, -1, fmt.Errorf("parse hash slice end error: %v", err)
		}
		return start, end, nil
	}

	return -1, -1, fmt.Errorf("invalid hash slice str")
}

// 如果是空字符串, 则返回0
// 如果是非数字, 则返回err
// 如果是数字, 则正常返回
func parseHashSliceValue(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(str)
	return v, err
}

// FindForKey return MycatPartitionStringShard calculated result
func (m *MycatPartitionStringShard) FindForKey(key interface{}) (int, error) {
	keyStr := GetString(key)
	var start int
	if m.hashSliceStart >= 0 {
		start = m.hashSliceStart
	} else {
		start = len(keyStr) + m.hashSliceStart
	}

	var end int
	if m.hashSliceEnd > 0 {
		end = m.hashSliceEnd
	} else {
		end = len(keyStr) + m.hashSliceEnd
	}
	h := stringHash(keyStr, start, end)
	return m.segment[int(h)&andValue], nil
}

// copied from mycat
func stringHash(s string, start, end int) int64 {
	input := []rune(s)
	if start < 0 {
		start = 0
	}

	if end > len(input) {
		end = len(input)
	}

	var h int64
	for i := start; i < end; i++ {
		h = (h << 5) - h + int64(input[i])
	}

	return h
}

const (
	defaultWeight = 1
)

// MycatPartitionMurmurHashShard mycat route algorithm: PartitionByMurmurHash
type MycatPartitionMurmurHashShard struct {
	seed               int
	count              int
	virtualBucketTimes int
	weightMap          map[int]int
	bucketMap          *treemap.Map
	hashFunction       *util.MurmurHash
}

// NewMycatPartitionMurmurHashShard constructor of MycatPartitionMurmurHashShard
func NewMycatPartitionMurmurHashShard(seedStr, virtualBucketTimesStr string, count int) (*MycatPartitionMurmurHashShard, error) {
	seed, err := strconv.Atoi(seedStr)
	if err != nil {
		return nil, err
	}
	if virtualBucketTimesStr == "" {
		virtualBucketTimesStr = "160"
	}
	virtualBucketTimes, err := strconv.Atoi(virtualBucketTimesStr)
	if err != nil {
		return nil, err
	}

	return &MycatPartitionMurmurHashShard{
		seed:               seed,
		count:              count,
		virtualBucketTimes: virtualBucketTimes,
	}, nil
}

// Init init MycatPartitionMurmurHashShard
func (m *MycatPartitionMurmurHashShard) Init() error {
	m.bucketMap = treemap.NewWithIntComparator()

	m.hashFunction = util.NewMurmurHash(m.seed) //计算一致性哈希的对象
	if err := m.generateBucketMap(); err != nil {
		return fmt.Errorf("murmur hash error, %v", err)
	}

	m.weightMap = make(map[int]int)
	return nil
}

// FindForKey return MycatPartitionMurmurHashShard calculated result
func (m *MycatPartitionMurmurHashShard) FindForKey(key interface{}) (int, error) {
	keyStr := GetString(key)
	hashKey := m.hashFunction.HashUnencodedChars(keyStr)
	_, ret := m.bucketMap.Ceiling(int(hashKey))
	if ret != nil {
		return ret.(int), nil
	}
	k, v := m.bucketMap.Min()
	if k == nil {
		return 0, errors.New("bucket map is empty")
	}
	return v.(int), nil
}

// SetWeightMapFromFile init weight
func (m *MycatPartitionMurmurHashShard) SetWeightMapFromFile(weightMapPath string) error {
	weightMap, err := parseWeightMapFile(weightMapPath)
	if err != nil {
		return err
	}
	m.weightMap = weightMap
	return nil
}

/**
 * 节点的权重，没有指定权重的节点默认是1。以properties文件的格式填写，以从0开始到count-1的整数值也就是节点索引为key，以节点权重值为值。
 * 所有权重值必须是正整数，否则以1代替
 */
// TODO: not implemented
func parseWeightMapFile(fileName string) (map[int]int, error) {
	return nil, errors.New("not implemented")
}

/**
 * 得到桶的权重，桶就是实际存储数据的DB实例
 * 从0开始的桶编号为key，权重为值，权重默认为1。
 * 键值必须都是整数
 *
 * @param bucket
 * @return
 */
func (m *MycatPartitionMurmurHashShard) getWeight(bucket int) int {
	w, ok := m.weightMap[bucket]
	if !ok {
		return defaultWeight
	}
	return w
}

func (m *MycatPartitionMurmurHashShard) generateBucketMap() error {
	hash := m.hashFunction
	for i := 0; i < m.count; i++ { //构造一致性哈希环，用TreeMap表示
		buf := bytes.NewBufferString("SHARD-" + strconv.Itoa(i))
		shard := m.virtualBucketTimes * m.getWeight(i)
		for n := 0; n < shard; n++ {
			// the hashKey is same with mycat, but maybe ugly, like:
			// SHARD-0-NODE-0, SHARD-0-NODE-0-NODE-1, SHARD-0-NODE-0-NODE-1-NODE-2,
			// not SHARD-0-NODE-0, SHARD-0-NODE-1, SHARD-0-NODE-2,
			// take care!
			buf.WriteString("-NODE-" + strconv.Itoa(n))
			hashKey := hash.HashUnencodedChars(buf.String())
			m.bucketMap.Put(hashKey, i)
		}
	}
	return nil
}

// mod padding
const (
	PaddingModLeftEnd  = 0
	PaddingModRightEnd = 1

	PaddingModDefaultPadFrom   = PaddingModRightEnd
	PaddingModDefaultPadLength = 18
	PaddingModDefaultModBegin  = 10
	PaddingModDefaultModEnd    = 16
	PaddingModDefaultMod       = 2
)

// MycatPartitionPaddingModShard mcyat partition padding mod
type MycatPartitionPaddingModShard struct {
	padFrom   int
	padLength int
	modBegin  int
	modEnd    int
	mod       int
}

func newDefaultMycatPartitionPaddingModShard() *MycatPartitionPaddingModShard {
	return &MycatPartitionPaddingModShard{
		padFrom:   PaddingModDefaultPadFrom,
		padLength: PaddingModDefaultPadLength,
		modBegin:  PaddingModDefaultModBegin,
		modEnd:    PaddingModDefaultModEnd,
		mod:       PaddingModDefaultMod,
	}
}

// GetMycatPartitionPaddingModShard wrapper newDefaultMycatPartitionPaddingModShard
func GetMycatPartitionPaddingModShard(padFromStr, padLengthStr, modBeginStr, modEndStr string, mod int) (shard *MycatPartitionPaddingModShard, err error) {
	padFrom, err := strconv.Atoi(padFromStr)
	if err != nil {
		return nil, err
	}

	padLength, err := strconv.Atoi(padLengthStr)
	if err != nil {
		return nil, err
	}

	modBegin, err := strconv.Atoi(modBeginStr)
	if err != nil {
		return nil, err
	}

	modEnd, err := strconv.Atoi(modEndStr)
	if err != nil {
		return nil, err
	}

	return newMycatPartitionPaddingModShard(padFrom, padLength, modBegin, modEnd, mod), nil
}

func newMycatPartitionPaddingModShard(padFrom, padLength, modBegin, modEnd, mod int) *MycatPartitionPaddingModShard {
	return &MycatPartitionPaddingModShard{
		padFrom:   padFrom,
		padLength: padLength,
		modBegin:  modBegin,
		modEnd:    modEnd,
		mod:       mod,
	}
}

// Init init MycatPartitionPaddingModShard and check params
func (m *MycatPartitionPaddingModShard) Init() error {
	return m.checkParam()
}

func (m *MycatPartitionPaddingModShard) checkParam() error {
	if m.padFrom != PaddingModLeftEnd && m.padFrom != PaddingModRightEnd {
		return fmt.Errorf("invalid padding mod mode: %d", m.padFrom)
	}

	if m.mod < PaddingModDefaultMod {
		return fmt.Errorf("invalid padding mod number: %d", m.mod)
	}

	if m.modBegin < 0 || m.modBegin >= m.modEnd {
		return fmt.Errorf("invalid padding modBegin or modEnd: %d, %d", m.modBegin, m.modEnd)
	}

	if m.padLength <= 0 {
		return fmt.Errorf("invalid padding mod padLength: %d", m.padLength)
	}

	if m.padLength < (m.modEnd - m.modBegin) {
		return fmt.Errorf("invalid padding mod, padLength is less than modBegin - modEnd: %d, %d, %d", m.padLength, m.modBegin, m.modEnd)
	}

	return nil
}

// FindForKey return MycatPartitionPaddingModShard calculated result
func (m *MycatPartitionPaddingModShard) FindForKey(key interface{}) (int, error) {
	h := NumValue(key) // assert the key is number
	keyStr := strconv.FormatInt(h, 10)

	var paddingKey string
	if len(keyStr) > m.padLength {
		paddingKey = keyStr[0:m.padLength]
	} else {
		if m.padFrom == PaddingModLeftEnd {
			if h < 0 { // Compatible with xiaomi mycat version
				return -1, errors.New("padding left to a negative is not allowed")
			}
			paddingKey = util.Left(keyStr, m.padLength, "0")
		} else {
			paddingKey = util.Right(keyStr, m.padLength, "0")
		}
	}

	modSegment := paddingKey[m.modBegin:m.modEnd]
	bigNum, err := strconv.ParseInt(modSegment, 10, 64) // in mycat, this is a BigInteger, but here is int64
	if err != nil {
		return -1, err
	}
	bigNumAbs := hack.Abs(bigNum)
	return int(bigNumAbs % int64(m.mod)), nil
}
