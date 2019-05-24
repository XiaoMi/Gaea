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
	"fmt"
	"testing"

	"github.com/XiaoMi/Gaea/mysql"
)

func TestLimitSelectResult(t *testing.T) {
	tests := []struct {
		offset       int64
		count        int64
		retLen       int64
		trimedRetLen int64
	}{
		{-1, -1, 5, 5},
		{0, -1, 5, 5},
		{0, 3, 5, 3},
		{0, 10, 5, 5},
		{3, 10, 5, 2},
		{3, 1, 5, 1},
		{4, 1, 5, 1},
		{5, 1, 5, 0},
		{5, 10, 5, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d:%d", test.offset, test.count), func(t *testing.T) {
			info := &SelectPlan{
				offset: test.offset,
				count:  test.count,
			}

			ret := &mysql.Result{
				Resultset: &mysql.Resultset{
					Values:   make([][]interface{}, test.retLen),
					RowDatas: make([]mysql.RowData, test.retLen),
				},
			}

			if err := limitSelectResult(info, ret); err != nil {
				t.Fatalf("limitSelectResult error: %v", err)
			}

			if int64(len(ret.Values)) != test.trimedRetLen {
				t.Errorf("len Values not equal, expect: %d, actual: %d", test.trimedRetLen, len(ret.Values))
			}
		})
	}
}
