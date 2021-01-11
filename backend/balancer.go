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
	"math/rand"
	"time"

	"github.com/XiaoMi/Gaea/core/errors"
)

type balancer struct {
	total       int
	lastIndex   int
	roundRobinQ []int
	nodeWeights []int
}

// calculate gcd ?
func gcd(ary []int) int {
	var i int
	min := ary[0]
	length := len(ary)
	for i = 0; i < length; i++ {
		if ary[i] < min {
			min = ary[i]
		}
	}

	for {
		isCommon := true
		for i = 0; i < length; i++ {
			if ary[i]%min != 0 {
				isCommon = false
				break
			}
		}
		if isCommon {
			break
		}
		min--
		if min < 1 {
			break
		}
	}
	return min
}

func newBalancer(nodeWeights []int, total int) *balancer {
	var sum int
	var s balancer
	s.total = total
	s.lastIndex = 0
	gcd := gcd(nodeWeights)

	for _, weight := range nodeWeights {
		sum += weight / gcd
	}

	s.roundRobinQ = make([]int, 0, sum)
	for index, weight := range nodeWeights {
		for j := 0; j < weight/gcd; j++ {
			s.roundRobinQ = append(s.roundRobinQ, index)
		}
	}

	//random order
	if 1 < len(s.nodeWeights) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < sum; i++ {
			x := r.Intn(sum)
			temp := s.roundRobinQ[x]
			other := sum % (x + 1)
			s.roundRobinQ[x] = s.roundRobinQ[other]
			s.roundRobinQ[other] = temp
		}
	}
	return &s
}

func (b *balancer) next() (int, error) {
	var index int
	queueLen := len(b.roundRobinQ)
	if queueLen == 0 {
		return 0, errors.ErrNoDatabase
	}
	if queueLen == 1 {
		index = b.roundRobinQ[0]
		return index, nil
	}

	b.lastIndex = b.lastIndex % queueLen
	index = b.roundRobinQ[b.lastIndex]
	if index >= b.total {
		return 0, errors.ErrNoDatabase
	}
	b.lastIndex++
	b.lastIndex = b.lastIndex % queueLen
	return index, nil
}
