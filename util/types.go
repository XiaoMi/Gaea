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

package util

import (
	"strings"

	"github.com/XiaoMi/Gaea/parser/format"
	types "github.com/XiaoMi/Gaea/parser/tidb-types"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
)

// GetValueExprResult copy from ValueExpr.Restore()
// TODO: 分表列是否需要支持等值比较NULL
func GetValueExprResult(n *driver.ValueExpr) (interface{}, error) {
	switch n.Kind() {
	case types.KindNull:
		return nil, nil // TODO: null str or nil?
	case types.KindInt64:
		return n.GetInt64(), nil
	case types.KindUint64:
		return n.GetUint64(), nil
	case types.KindFloat32:
		return n.GetFloat32(), nil
	case types.KindFloat64:
		return n.GetFloat64(), nil
	case types.KindString, types.KindBytes:
		return n.GetString(), nil
	default:
		s := &strings.Builder{}
		ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
		err := n.Restore(ctx)
		if err != nil {
			return nil, err
		}
		return s.String(), nil
	}
}
