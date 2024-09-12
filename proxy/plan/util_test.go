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

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func testCheckList(t *testing.T, l []int, checkList ...int) {
	if len(l) != len(checkList) {
		t.Fatal("invalid list len", len(l), len(checkList))
	}

	for i := 0; i < len(l); i++ {
		if l[i] != checkList[i] {
			t.Fatal("invalid list item", l[i], i)
		}
	}
}

func TestListSet(t *testing.T) {
	var l1 []int
	var l2 []int
	var l3 []int

	l1 = []int{1, 2, 3}
	l2 = []int{2}

	l3 = interList(l1, l2)
	testCheckList(t, l3, 2)

	l1 = []int{1, 2, 3}
	l2 = []int{2, 3}

	l3 = interList(l1, l2)
	testCheckList(t, l3, 2, 3)

	l1 = []int{1, 2, 4}
	l2 = []int{2, 3}

	l3 = interList(l1, l2)
	testCheckList(t, l3, 2)

	l1 = []int{1, 2, 4}
	l2 = []int{}

	l3 = interList(l1, l2)
	testCheckList(t, l3)

	l1 = []int{1, 2, 3}
	l2 = []int{2}

	l3 = unionList(l1, l2)
	testCheckList(t, l3, 1, 2, 3)

	l1 = []int{1, 2, 4}
	l2 = []int{3}

	l3 = unionList(l1, l2)
	testCheckList(t, l3, 1, 2, 3, 4)

	l1 = []int{1, 2, 3}
	l2 = []int{2, 3, 4}

	l3 = unionList(l1, l2)
	testCheckList(t, l3, 1, 2, 3, 4)

	l1 = []int{1, 2, 3}
	l2 = []int{}

	l3 = unionList(l1, l2)
	testCheckList(t, l3, 1, 2, 3)

	l1 = []int{1, 2, 3, 4}
	l2 = []int{2}

	l3 = differentList(l1, l2)
	testCheckList(t, l3, 1, 3, 4)

	l1 = []int{1, 2, 3, 4}
	l2 = []int{}

	l3 = differentList(l1, l2)
	testCheckList(t, l3, 1, 2, 3, 4)

	l1 = []int{1, 2, 3, 4}
	l2 = []int{1, 3, 5}

	l3 = differentList(l1, l2)
	testCheckList(t, l3, 2, 4)

	l1 = []int{1, 2, 3}
	l2 = []int{1, 3, 5, 6}

	l3 = differentList(l1, l2)
	testCheckList(t, l3, 2)

	l1 = []int{1, 2, 3, 4}
	l2 = []int{2, 3}

	l3 = differentList(l1, l2)
	testCheckList(t, l3, 1, 4)

	l1 = []int{1, 2, 2, 1, 5, 3, 5, 2}
	l2 = cleanList(l1)
	sort.Sort(sort.IntSlice(l2))
	testCheckList(t, l2, 1, 2, 3, 5)
}

func TestMakeLeList(t *testing.T) {
	l1 := []int{20150802, 20150812, 20150822, 20150823, 20150825, 20150828}
	l2 := makeLeList(20150822, l1)
	testCheckList(t, l2, 20150802, 20150812, 20150822)
	l3 := makeLeList(20150824, l1)
	testCheckList(t, l3, []int{}...)
}

func TestMakeLtList(t *testing.T) {
	l1 := []int{20150802, 20150812, 20150822, 20150823, 20150825, 20150828}
	l2 := makeLtList(20150822, l1)
	testCheckList(t, l2, 20150802, 20150812)
	l3 := makeLtList(20150824, l1)
	testCheckList(t, l3, []int{}...)
	l4 := makeLtList(20150802, l1)
	testCheckList(t, l4, []int{}...)
}

func TestMakeGeList(t *testing.T) {
	l1 := []int{20150802, 20150812, 20150822, 20150823, 20150825, 20150828}
	l2 := makeGeList(20150822, l1)
	testCheckList(t, l2, 20150822, 20150823, 20150825, 20150828)
	l3 := makeGeList(20150828, l1)
	testCheckList(t, l3, 20150828)
}

func TestMakeGtList(t *testing.T) {
	l1 := []int{20150802, 20150812, 20150822, 20150823, 20150825, 20150828}
	l2 := makeGtList(20150822, l1)
	testCheckList(t, l2, 20150823, 20150825, 20150828)
	l3 := makeGtList(20150824, l1)
	testCheckList(t, l3, []int{}...)
	l4 := makeGtList(20150828, l1)
	testCheckList(t, l4, []int{}...)
}

func TestMakeBetweenList(t *testing.T) {
	// 测试正常情况
	t.Run("Normal Case", func(t *testing.T) {
		indexs := []int{10, 20, 30, 40, 50}
		expected := []int{20, 30, 40}
		result := makeBetweenList(20, 40, indexs)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	// 测试反转的情况
	t.Run("Reversed Range", func(t *testing.T) {
		indexs := []int{10, 20, 30, 40, 50}
		expected := []int{20, 30, 40}
		result := makeBetweenList(40, 20, indexs)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	// 测试没有找到起始或结束索引的情况
	t.Run("Not Found Start or End Index", func(t *testing.T) {
		indexs := []int{10, 20, 30, 40, 50}
		expected := ([]int)(nil)
		result := makeBetweenList(25, 35, indexs)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	// 测试空列表的情况
	t.Run("Empty List", func(t *testing.T) {
		indexs := []int{}
		expected := ([]int)(nil)
		result := makeBetweenList(10, 20, indexs)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	// 测试包含重复元素的情况
	t.Run("Duplicate Elements", func(t *testing.T) {
		indexs := []int{10, 20, 30, 20, 50}
		expected := []int{20, 30}
		result := makeBetweenList(20, 30, indexs)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})
}

func TestRemoveDuplicatesString(t *testing.T) {
	testCases := []struct {
		arr1     []string
		arr2     []string
		expected []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{"a"}, []string{"b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"d", "e", "f"}, []string{"a", "b", "c", "d", "e", "f"}},
		{[]string{"a", "a", "b", "c"}, []string{"b", "d", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
		{[]string{"a", "b", "b", "c", "c", "c"}, []string{"a", "a", "d", "d", "d", "d"}, []string{"a", "b", "c", "d"}},
		{[]string{"a", "a", "a", "a", "a"}, []string{"b", "b", "b", "b", "b"}, []string{"a", "b"}},
	}

	for _, tc := range testCases {
		output := removeDuplicatesString(tc.arr1, tc.arr2)
		assert.Equal(t, len(output), len(tc.expected))
	}
}
