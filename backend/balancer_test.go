// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func checkRoundRobin(rb []int, indices []int, weights []int) bool {
	countMap := make(map[int]int)
	for _, node := range rb {
		countMap[node]++
	}

	gcdVal := gcd(weights)

	// 使用 indices 确保检查的是 roundRobinQ 里的 index，而不是 weights 的索引
	for i, poolIdx := range indices {
		expected := weights[i] / gcdVal
		if countMap[poolIdx] != expected {
			fmt.Printf("Mismatch for index %d (node %d): got %d, expected %d\n", i, poolIdx, countMap[poolIdx], expected)
			return false
		}
	}
	return true
}

func TestBalancer(t *testing.T) {
	indices := []int{0, 1, 2}
	weights := []int{2, 1, 4}
	b, err := newBalancer(indices, weights)
	assert.NoError(t, err)
	assert.True(t, checkRoundRobin(b.roundRobinQ, indices, weights))

	indices = []int{2, 1, 0}
	weights = []int{2, 1, 4}
	b, err = newBalancer(indices, weights)
	assert.NoError(t, err)
	assert.True(t, checkRoundRobin(b.roundRobinQ, indices, weights))

}

func TestBalancerSequence(t *testing.T) {
	indices := []int{0, 1, 2, 3}
	weights := []int{2, 2, 2, 4}
	b, err := newBalancer(indices, weights)
	assert.NoError(t, err)
	assert.True(t, checkRoundRobin(b.roundRobinQ, indices, weights))
}

func TestBalancerNext(t *testing.T) {
	// 假设连接池有4个节点，下标为 [0,1,2,3] 权重分别为[1,2,3,4]
	indices := []int{0, 1}
	weights := []int{1, 2}
	b, err := newBalancer(indices, weights)
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.NoError(t, err)
		assert.Contains(t, indices, idx)
	}

	indices = []int{1, 0}
	weights = []int{2, 1}
	b, err = newBalancer(indices, weights)
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.NoError(t, err)
		assert.Contains(t, indices, idx)
	}

	indices = []int{2, 3}
	weights = []int{3, 4}
	b, err = newBalancer(indices, weights)
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.NoError(t, err)
		assert.Contains(t, indices, idx)
	}
}

func TestBalancerLocalDataCenter(t *testing.T) {
	// 假设连接池有4个节点，下标为 [0,1,2,3] 权重分别为[1,2,3,4] 匹配本地数据中心的 indices 为 [2,3]

	localIndices := []int{2, 3} // 本地数据中心的 indices
	localWeights := []int{3, 4} // 本地数据中心的权重
	b, err := newBalancer(localIndices, localWeights)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.Nil(t, err)
		assert.Contains(t, localIndices, idx) // 只允许返回 2 或 3
	}
}

// TestNewBalancer_ValidInputs 测试正常输入情况
func TestNewBalancer_ValidInputs(t *testing.T) {
	indices := []int{0, 1, 2}
	weights := []int{2, 3, 5}

	b, err := newBalancer(indices, weights)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, len(b.roundRobinQ), 10)

	countMap := make(map[int]int)
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.Nil(t, err)
		countMap[idx]++
	}

	assert.True(t, checkRoundRobin(b.roundRobinQ, indices, weights))
	assert.Equal(t, countMap[0], 2)
	assert.Equal(t, countMap[1], 3)
	assert.Equal(t, countMap[2], 5)
}

func TestNewBalancer_ValidInputsShuffleIndices(t *testing.T) {
	indices := []int{2, 1, 0}
	weights := []int{2, 3, 5}

	b, err := newBalancer(indices, weights)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, len(b.roundRobinQ), 10)

	countMap := make(map[int]int)
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.Nil(t, err)
		countMap[idx]++
	}

	assert.True(t, checkRoundRobin(b.roundRobinQ, indices, weights))
	assert.Equal(t, countMap[2], 2)
	assert.Equal(t, countMap[1], 3)
	assert.Equal(t, countMap[0], 5)
}

// TestNewBalancer_ZeroWeight 测试包含 `0` 权重的情况
func TestNewBalancer_ZeroWeight(t *testing.T) {
	indices := []int{0, 1, 2}
	weights := []int{2, 0, 5} // index `1` 权重为 `0`，应被忽略

	b, err := newBalancer(indices, weights)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.NotContains(t, b.roundRobinQ, 1) // 确保 `1` 没有出现在 roundRobinQ
	for i := 0; i < 10; i++ {
		idx, err := b.next()
		assert.Nil(t, err)
		assert.NotEqual(t, 1, idx) // `1` 绝不应该被选中
	}
}

// TestNewBalancer_EmptyIndices 测试空索引列表（不报错，balancer 为 nil）
func TestNewBalancer_EmptyIndices(t *testing.T) {
	b, err := newBalancer([]int{}, []int{})
	assert.Nil(t, b)
	assert.Nil(t, err)
}

// TestNewBalancer_IndexWeightMismatch 测试索引和权重数量不匹配（应 panic）
func TestNewBalancer_IndexWeightMismatch(t *testing.T) {
	b, err := newBalancer([]int{0, 1}, []int{2})
	assert.Nil(t, b)
	assert.NotNil(t, err)

}

// TestBalancer_Next_EmptyQueue 测试 next() 在 roundRobinQ 为空的情况下
func TestBalancer_Next_EmptyQueue(t *testing.T) {
	b := &balancer{} // 创建一个空的 balancer
	_, err := b.next()
	assert.NotNil(t, err)
	assert.Equal(t, "no candidate available in balancer", err.Error())
}

// TestBalancer_WeightDistribution 确保 next() 按权重比例返回索引
func TestBalancer_WeightDistribution(t *testing.T) {
	indices := []int{0, 1, 2}
	weights := []int{2, 3, 5} // 预期比例为 2:3:5

	b, err := newBalancer(indices, weights)
	assert.Nil(t, err)
	countMap := make(map[int]int)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		idx, err := b.next()
		assert.Nil(t, err)
		countMap[idx]++
	}

	// 计算实际比例
	total := countMap[0] + countMap[1] + countMap[2]
	expectedRatio := []float64{0.2, 0.3, 0.5} // 2:3:5 的比例
	actualRatio := []float64{
		float64(countMap[0]) / float64(total),
		float64(countMap[1]) / float64(total),
		float64(countMap[2]) / float64(total),
	}

	for i, ratio := range expectedRatio {
		assert.Equal(t, ratio, actualRatio[i])
	}
}

func TestBalancerConcurrentDistribution(t *testing.T) {
	nodeWeights := []int{2, 3, 5}
	totalNodes := len(nodeWeights)
	indices := make([]int, totalNodes)
	for i := range indices {
		indices[i] = i
	}

	b, err := newBalancer(indices, nodeWeights)
	assert.NoError(t, err)
	var counts = make([]int64, totalNodes)

	var wg sync.WaitGroup
	concurrency := 100
	callsPerGoroutine := 1000
	totalCalls := concurrency * callsPerGoroutine

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				idx, err := b.next()
				assert.NoError(t, err)
				assert.Less(t, idx, totalNodes)
				atomic.AddInt64(&counts[idx], 1)
			}
		}()
	}
	wg.Wait()

	expected := make([]int64, totalNodes)
	totalWeight := 0
	for _, w := range nodeWeights {
		totalWeight += w
	}

	for i, w := range nodeWeights {
		expected[i] = int64(totalCalls * w / totalWeight)
	}

	for i := 0; i < totalNodes; i++ {
		assert.InDelta(t, expected[i], counts[i], float64(expected[i])*0.05) // 允许5%的误差
	}
}

func TestBalancerSequenceNonFixedSeed(t *testing.T) {
	indices := []int{0, 1, 2}
	weights := []int{2, 3, 5}

	// 1. 调用 newBalancer，根据权重生成并随机打乱 roundRobinQ
	b, err := newBalancer(indices, weights)
	assert.NoError(t, err)
	queueLen := len(b.roundRobinQ)
	if queueLen == 0 {
		t.Fatal("roundRobinQ is empty, can't test sequence")
	}

	// 2.统计在多次调用 next() 中每个 index 出现的次数, 调用次数足够多
	frequency := make(map[int]int)
	iterations := 3000
	expectedResult := map[int]int{
		0: 600, 1: 900, 2: 1500,
	}

	// 3. 连续调用 next()，检查返回结果是否与 expectedQueue 循环一致
	for i := 0; i < iterations; i++ {
		idx, err := b.next()
		if err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
		frequency[idx]++
	}

	// 4. 检查每个 index 出现的次数是否是期望值
	for i := 0; i < len(indices); i++ {
		got := frequency[i]
		want := expectedResult[i]
		if got != want {
			t.Errorf("index %d: got frequency=%d, want=%d", i, got, want)
		}
	}
}

func TestBalancerSequenceFixedSeed(t *testing.T) {
	mockey.PatchConvey("test LocalSlaveReadPreferred local one down", t, func() {
		// 固定随机时间种子
		mockey.Mock(time.Now).Return(time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)).Build()
		indices := []int{0, 1, 2}
		weights := []int{2, 3, 5}

		// 1. 调用 newBalancer，根据权重生成并随机打乱 roundRobinQ
		b, err := newBalancer(indices, weights)
		assert.NoError(t, err)
		queueLen := len(b.roundRobinQ)
		if queueLen == 0 {
			t.Fatal("roundRobinQ is empty, can't test sequence")
		}

		// 2.统计在多次调用 next() 中每个 index 出现的次数, 调用次数足够多
		frequency := make(map[int]int)
		iterations := 3000
		expectedResult := map[int]int{
			0: 600, 1: 900, 2: 1500,
		}

		// 3. 连续调用 next()，检查返回结果是否与 expectedQueue 循环一致
		for i := 0; i < iterations; i++ {
			idx, err := b.next()
			if err != nil {
				t.Fatalf("unexpected error on call %d: %v", i, err)
			}
			frequency[idx]++
		}

		// 4. 检查每个 index 出现的次数是否是期望值
		for i := 0; i < len(indices); i++ {
			got := frequency[i]
			want := expectedResult[i]
			if got != want {
				t.Errorf("index %d: got frequency=%d, want=%d", i, got, want)
			}
		}
	})
}

func TestBalancerNextIndexMax(t *testing.T) {
	var u uint32 = math.MaxUint32
	res := atomic.AddUint32(&u, 1)
	if int(res) < 0 {
		t.Errorf("index %d: overflow", u)
	}
}
