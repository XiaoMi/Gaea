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

func TestSet(t *testing.T) {
	defer testleak.AfterTestT(t)()
	elems := []string{"a", "b", "c", "d"}
	tbl := []struct {
		Name          string
		ExpectedValue uint64
		ExpectedName  string
	}{
		{"a", 1, "a"},
		{"a,b,a", 3, "a,b"},
		{"b,a", 3, "a,b"},
		{"a,b,c,d", 15, "a,b,c,d"},
		{"d", 8, "d"},
		{"", 0, ""},
		{"0", 0, ""},
	}

	for _, tt := range tbl {
		e, err := ParseSetName(elems, tt.Name)
		require.NoError(t, err)
		require.Equal(t, float64(tt.ExpectedValue), e.ToNumber())
		require.Equal(t, tt.ExpectedName, e.String())
	}

	tblNumber := []struct {
		Number       uint64
		ExpectedName string
	}{
		{0, ""},
		{1, "a"},
		{3, "a,b"},
		{9, "a,d"},
	}

	for _, tt := range tblNumber {
		e, err := ParseSetValue(elems, tt.Number)
		require.NoError(t, err)
		require.Equal(t, tt.ExpectedName, e.String())
		require.Equal(t, float64(tt.Number), e.ToNumber())
	}

	tblErr := []string{
		"a.e",
		"e.f",
	}
	for _, tt := range tblErr {
		_, err := ParseSetName(elems, tt)
		require.Error(t, err)
	}

	tblNumberErr := []uint64{
		100, 16, 64,
	}
	for _, tt := range tblNumberErr {
		_, err := ParseSetValue(elems, tt)
		require.Error(t, err)
	}
}
