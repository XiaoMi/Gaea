// Copyright 2019 PingCAP, Inc.
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
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckFsp(t *testing.T) {

	obtained, err := CheckFsp(UnspecifiedFsp)
	require.Equal(t, DefaultFsp, obtained)
	require.NoError(t, err)

	obtained, err = CheckFsp(-2019)
	require.Equal(t, DefaultFsp, obtained)
	require.EqualError(t, err, "Invalid fsp -2019")

	obtained, err = CheckFsp(MinFsp - 4294967296)
	require.Equal(t, DefaultFsp, obtained)
	require.EqualError(t, err, "Invalid fsp "+strconv.Itoa(MinFsp-4294967296))

	// UnspecifiedFsp
	obtained, err = CheckFsp(-1)
	require.Equal(t, DefaultFsp, obtained)
	require.NoError(t, err)

	obtained, err = CheckFsp(MaxFsp + 1)
	require.Equal(t, DefaultFsp, obtained)
	require.EqualError(t, err, "Invalid fsp "+strconv.Itoa(MaxFsp+1))

	obtained, err = CheckFsp(MaxFsp + 2019)
	require.Equal(t, DefaultFsp, obtained)
	require.EqualError(t, err, "Invalid fsp "+strconv.Itoa(MaxFsp+2019))

	obtained, err = CheckFsp(MaxFsp + 4294967296)
	require.Equal(t, DefaultFsp, obtained)
	require.EqualError(t, err, "Invalid fsp "+strconv.Itoa(MaxFsp+4294967296))

	obtained, err = CheckFsp((MaxFsp + MinFsp) / 2)
	require.Equal(t, (MaxFsp+MinFsp)/2, obtained)
	require.NoError(t, err)

	obtained, err = CheckFsp(5)
	require.Equal(t, 5, obtained)
	require.NoError(t, err)

}

func TestParseFrac(t *testing.T) {
	obtained, overflow, err := ParseFrac("", 5)
	require.Equal(t, 0, obtained)
	require.False(t, overflow)
	require.NoError(t, err)

	obtained, overflow, err = ParseFrac("999", 200)
	require.Equal(t, 0, obtained)
	require.False(t, overflow)
	require.Error(t, err)
	require.Regexp(t, "^Invalid fsp ", err.Error())

	obtained, overflow, err = ParseFrac("NotNum", MaxFsp)
	require.Equal(t, 0, obtained)
	require.False(t, overflow)
	require.Error(t, err)
	require.Regexp(t, "strconv.ParseInt:.*", err.Error())

	obtained, overflow, err = ParseFrac("1235", 6)
	require.Equal(t, 123500, obtained)
	require.False(t, overflow)
	require.NoError(t, err)

	obtained, overflow, err = ParseFrac("123456", 4)
	require.Equal(t, 123500, obtained)
	require.False(t, overflow)
	require.NoError(t, err)

	// 1236 round 3 -> 124 -> 124000
	obtained, overflow, err = ParseFrac("1236", 3)
	require.Equal(t, 124000, obtained)
	require.False(t, overflow)
	require.NoError(t, err)

	// 03123 round 2 -> 3 -> 30000
	obtained, overflow, err = ParseFrac("0312", 2)
	require.Equal(t, 30000, obtained)
	require.False(t, overflow)
	require.NoError(t, err)

	// 999 round 2 -> 100 -> overflow
	obtained, overflow, err = ParseFrac("999", 2)
	require.Equal(t, 0, obtained)
	require.True(t, overflow)
	require.NoError(t, err)
}

func TestAlignFrac(t *testing.T) {
	obtained := alignFrac("100", 6)
	require.Equal(t, "100000", obtained)
	obtained = alignFrac("10000000000", 6)
	require.Equal(t, "10000000000", obtained)
}
