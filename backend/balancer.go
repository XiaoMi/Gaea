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

package backend

import (
	"fmt"
	"sync"
	"time"

	"math/rand"
)

type balancer struct {
	mu          sync.Mutex // 保护并发安全
	nextIndex   int        // 当前轮询指针，指向 roundRobinQ 中的下一个候选
	roundRobinQ []int      // 按照权重扩展后的候选队列（存储原始连接池的下标）
	poolIndices []int      // 连接池的下标
	poolWeights []int      // 连接池的下标对应的权重
}

// 计算 GCD（欧几里得算法）
func gcd(ary []int) int {
	if len(ary) == 0 {
		return 1
	}

	gcdHelper := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	g := ary[0]
	for i := 1; i < len(ary); i++ {
		g = gcdHelper(g, ary[i])
		if g == 1 {
			return 1 // 归一化到最小单位
		}
	}
	return g
}

// newBalancer 根据传入的节点下标和对应权重构造新的 balancer, 根据 indices 和对应的 weights（按 gcd 归一化），扩展 roundRobinQ 以进行权重轮询。
// 每次调用 next()，返回 roundRobinQ 中的下一个节点（基于权重比例）。
// getIndicesAndWeights中过滤掉了权重为0和为负数的节点
func newBalancer(indices []int, weights []int) (*balancer, error) {

	if len(indices) != len(weights) {
		return nil, fmt.Errorf("indices and weights length mismatch, indices: %d, weights: %d", len(indices), len(weights))
	}
	if len(indices) == 0 {
		return nil, nil
	}
	// 计算 GCD，确保权重按比例缩放
	gcdVal := gcd(weights)

	// 计算队列所需大小
	var sum int
	for _, weight := range weights {
		sum += weight / gcdVal
	}

	// 预分配 `queue`（避免动态扩容带来的额外开销）
	queue := make([]int, 0, len(indices))

	// 构造 roundRobinQ
	for i, idx := range indices {
		repeat := weights[i] / gcdVal
		for j := 0; j < repeat; j++ {
			queue = append(queue, idx)
		}
	}

	// 生成 balancer
	b := &balancer{
		mu:          sync.Mutex{},
		nextIndex:   0,
		roundRobinQ: queue,
		poolWeights: weights,
		poolIndices: indices,
	}

	// 使 `roundRobinQ` 随机化，防止固定分布
	if len(b.roundRobinQ) > 1 {
		rand.Seed(time.Now().UnixNano()) // 使用时间戳作为随机种子
		rand.Shuffle(len(b.roundRobinQ), func(i, j int) {
			b.roundRobinQ[i], b.roundRobinQ[j] = b.roundRobinQ[j], b.roundRobinQ[i]
		})
	}

	return b, nil
}

// next 返回 roundRobinQ 中下一个候选的节点下标，
func (b *balancer) next() (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.roundRobinQ) == 0 {
		return 0, fmt.Errorf("no candidate available in balancer")
	}
	if len(b.roundRobinQ) == 1 {
		return b.roundRobinQ[0], nil
	}
	// 直接取出 roundRobinQ 中的元素
	elem := b.roundRobinQ[b.nextIndex]

	// 轮询更新 nextIndex
	b.nextIndex = (b.nextIndex + 1) % len(b.roundRobinQ)

	return elem, nil
}
