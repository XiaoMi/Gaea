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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitMarginComments(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantQuery    string
		wantComments MarginComments
	}{
		{
			name:      "test no comment",
			sql:       "select * from t",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: "",
			},
		},
		{
			name:      "test leading comment",
			sql:       "/*master*/ select * from t",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "/*master*/ ",
				Trailing: "",
			},
		},
		{
			name:      "test internal comment",
			sql:       "select /*master*/ * from t",
			wantQuery: "select /*master*/ * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: "",
			},
		},
		{
			name:      "test mycat hint comment",
			sql:       "select * from t /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: " /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
			},
		},
		{
			name:      "test trailing comment",
			sql:       "select * from t /*master*/",
			wantQuery: "select * from t",
			wantComments: MarginComments{
				Leading:  "",
				Trailing: " /*master*/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotComments := SplitMarginComments(tt.sql)
			assert.Equalf(t, tt.wantQuery, gotQuery, "SplitMarginComments(%v)", tt.sql)
			assert.Equalf(t, tt.wantComments, gotComments, "SplitMarginComments(%v)", tt.sql)
		})
	}
}
