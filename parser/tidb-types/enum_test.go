// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/util/testleak"
)

func TestEnum(t *testing.T) {
	defer testleak.AfterTestT(t)()
	tbl := []struct {
		Elems    []string
		Name     string
		Expected int
	}{
		{[]string{"a", "b"}, "a", 1},
		{[]string{"a"}, "b", 0},
		{[]string{"a"}, "1", 1},
	}

	for _, tt := range tbl {
		e, err := ParseEnumName(tt.Elems, tt.Name)
		if tt.Expected == 0 {
			require.Error(t, err)
			require.Equal(t, float64(0), e.ToNumber())
			require.Equal(t, "", e.String())
			continue
		}

		require.NoError(t, err)
		require.Equal(t, tt.Elems[tt.Expected-1], e.String())
		require.Equal(t, float64(tt.Expected), e.ToNumber())
	}

	tblNumber := []struct {
		Elems    []string
		Number   uint64
		Expected int
	}{
		{[]string{"a"}, 1, 1},
		{[]string{"a"}, 0, 0},
	}

	for _, tt := range tblNumber {
		e, err := ParseEnumValue(tt.Elems, tt.Number)
		if tt.Expected == 0 {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, float64(tt.Expected), e.ToNumber())
	}
}
