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

import "testing"

func TestGetString(t *testing.T) {
	tests := []struct {
		v  interface{}
		vs string
	}{
		{int(-1), "-1"},
		{int64(-1), "-1"},
		{"-1", "-1"},
		{[]byte("-1"), "-1"},
		{int(0), "0"},
		{int64(0), "0"},
		{uint(0), "0"},
		{uint64(0), "0"},
		{"0", "0"},
		{[]byte("0"), "0"},
		{int(1), "1"},
		{int64(1), "1"},
		{uint(1), "1"},
		{uint64(1), "1"},
		{"1", "1"},
		{[]byte("1"), "1"},
	}
	for _, test := range tests {
		t.Run(test.vs, func(t *testing.T) {
			actualVs := GetString(test.v)
			if actualVs != test.vs {
				t.Errorf("not equal, v: %v, expect: %s, actual: %s", test.v, test.vs, actualVs)
			}
		})
	}
}
