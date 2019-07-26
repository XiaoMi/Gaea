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

package plan

/*2,5 ==> [2,3,4]*/
func makeList(start, end int) []int {
	if start >= end {
		return []int{}
	}
	list := make([]int, end-start)
	for i := start; i < end; i++ {
		list[i-start] = i
	}
	return list
}

//if value is 2016, and indexs is [2015,2016,2017]
//the result is [2015,2016]
// the indexs must be sorted
func makeLeList(value int, indexs []int) []int {
	for k, v := range indexs {
		if v == value {
			return indexs[:k+1]
		}
	}
	return nil
}

//if value is 2016, and indexs is [2015,2016,2017,2018]
//the result is [2016,2017,2018]
// the indexs must be sorted
func makeGeList(value int, indexs []int) []int {
	for k, v := range indexs {
		if v == value {
			return indexs[k:]
		}
	}
	return nil
}

//if value is 2016, and indexs is [2015,2016,2017,2018]
//the result is [2015]
// the indexs must be sorted
func makeLtList(value int, indexs []int) []int {
	for k, v := range indexs {
		if v == value {
			return indexs[:k]
		}
	}
	return nil
}

//if value is 2016, and indexs is [2015,2016,2017,2018]
//the result is [2017,2018]
// the indexs must be sorted
func makeGtList(value int, indexs []int) []int {
	for k, v := range indexs {
		if v == value {
			return indexs[k+1:]
		}
	}
	return nil
}

//if start is 2016, end is 2017. indexs is [2015,2016,2017,2018]
//the result is [2016,2017]
// the indexs must be sorted
func makeBetweenList(start, end int, indexs []int) []int {
	var startIndex, endIndex int
	var SetStart bool
	if end < start {
		start, end = end, start
	}
	for k, v := range indexs {
		if v == start {
			startIndex = k
			SetStart = true
		}
		if v == end {
			endIndex = k
			if SetStart {
				return indexs[startIndex : endIndex+1]
			}
		}
	}
	return nil
}

// l1 & l2
func interList(l1 []int, l2 []int) []int {
	if len(l1) == 0 || len(l2) == 0 {
		return []int{}
	}

	l3 := make([]int, 0, len(l1)+len(l2))
	var i = 0
	var j = 0
	for i < len(l1) && j < len(l2) {
		if l1[i] == l2[j] {
			l3 = append(l3, l1[i])
			i++
			j++
		} else if l1[i] < l2[j] {
			i++
		} else {
			j++
		}
	}

	return l3
}

// l1 | l2
func unionList(l1 []int, l2 []int) []int {
	if len(l1) == 0 {
		return l2
	} else if len(l2) == 0 {
		return l1
	}

	l3 := make([]int, 0, len(l1)+len(l2))

	var i = 0
	var j = 0
	for i < len(l1) && j < len(l2) {
		if l1[i] < l2[j] {
			l3 = append(l3, l1[i])
			i++
		} else if l1[i] > l2[j] {
			l3 = append(l3, l2[j])
			j++
		} else {
			l3 = append(l3, l1[i])
			i++
			j++
		}
	}

	if i != len(l1) {
		l3 = append(l3, l1[i:]...)
	} else if j != len(l2) {
		l3 = append(l3, l2[j:]...)
	}

	return l3
}

// l1 - l2
func differentList(l1 []int, l2 []int) []int {
	if len(l1) == 0 {
		return []int{}
	} else if len(l2) == 0 {
		return l1
	}

	l3 := make([]int, 0, len(l1))

	var i = 0
	var j = 0
	for i < len(l1) && j < len(l2) {
		if l1[i] < l2[j] {
			l3 = append(l3, l1[i])
			i++
		} else if l1[i] > l2[j] {
			j++
		} else {
			i++
			j++
		}
	}

	if i != len(l1) {
		l3 = append(l3, l1[i:]...)
	}

	return l3
}

func cleanList(l []int) []int {
	s := make(map[int]struct{})
	listLen := len(l)
	l2 := make([]int, 0, listLen)

	for i := 0; i < listLen; i++ {
		k := l[i]
		s[k] = struct{}{}
	}
	for k := range s {
		l2 = append(l2, k)
	}
	return l2
}

// list l need to be sorted
func distinctList(l []int) []int {
	if len(l) < 2 {
		return l
	}
	var ret []int
	current := l[0]
	ret = append(ret, current)
	for i := 1; i < len(l); i++ {
		if l[i] != current {
			current = l[i]
			ret = append(ret, current)
		}
	}
	return ret
}
