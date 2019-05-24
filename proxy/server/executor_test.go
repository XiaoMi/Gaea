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
package server

import (
	"testing"

	"fmt"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
)

func TestGetVariableExprResult(t *testing.T) {
	tests := []struct {
		variable []string
		expect   string
	}{
		{[]string{"ON", "on", "'on'", "`on`"}, "on"},
		{[]string{"OFF", "off", "'off'", "`off`"}, "off"},
		{[]string{"1", "'1'", "`1`"}, "1"},
		{[]string{"0", "'0'", "`0`"}, "0"},
	}
	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			for _, v := range test.variable {
				sql := fmt.Sprintf("set autocommit = %s", v)
				s, err := parser.ParseSQL(sql)
				if err != nil {
					t.Fatal(err)
				}
				stmt := s.(*ast.SetStmt)
				for _, v := range stmt.Variables {
					actual := getVariableExprResult(v.Value)
					if actual != test.expect {
						t.Errorf("not equal, expect: %v, actual: %v", test.expect, actual)
					}
				}
			}
		})
	}
}
