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

package mysql

import (
	"fmt"
	"testing"
)

var errTimeTpl = "invalid TypeDuration %s"

func Test_StringToMysqlTime(t *testing.T) {
	tests := []struct {
		timestr string
		value   TimeValue
		err     error
	}{
		{"00:00", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00")},
		{"00:00:00", TimeValue{}, nil},
		{"00:00:00.", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00:00.")},
		{"00:00:00.00", TimeValue{}, nil},
		{"00:00:00.000001", TimeValue{Microsecond: 1}, nil},
		{"00:00:00.001", TimeValue{Microsecond: 1000}, nil},
		{"00:00:00.999999", TimeValue{Microsecond: 999999}, nil},
		{"00:00:00.9999991", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00:00.9999991")},
		{"00:00:00.-01", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00:00.-01")},
		{"00:00:01", TimeValue{Second: 1}, nil},
		{"00:00:59", TimeValue{Second: 59}, nil},
		{"00:00:60", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00:60")},
		{"00:00:-01", TimeValue{}, fmt.Errorf(errTimeTpl, "00:00:-01")},
		{"00:01:00", TimeValue{Minute: 1}, nil},
		{"00:59:00", TimeValue{Minute: 59}, nil},
		{"00:60:00", TimeValue{}, fmt.Errorf(errTimeTpl, "00:60:00")},
		{"00:-01:00", TimeValue{}, fmt.Errorf(errTimeTpl, "00:-01:00")},
		{"00:1.5:00", TimeValue{}, fmt.Errorf(errTimeTpl, "00:1.5:00")},
		{"01:00:00", TimeValue{Hour: 1}, nil},
		{"23:00:00", TimeValue{Hour: 23}, nil},
		{"24:00:00", TimeValue{Day: 1, Hour: 0}, nil},
		{"25:00:00", TimeValue{Day: 1, Hour: 1}, nil},
		{"48:00:00", TimeValue{Day: 2, Hour: 0}, nil},
		{"1.1:00:00", TimeValue{}, fmt.Errorf(errTimeTpl, "1.1:00:00")},

		// negative time
		{"-00:00:00.00", TimeValue{IsNegative: true}, nil}, // mysql doesn't return this value
		{"-00:00:00.000001", TimeValue{IsNegative: true, Microsecond: 1}, nil},
		{"-00:00:00.001", TimeValue{IsNegative: true, Microsecond: 1000}, nil},
		{"-00:00:00.999999", TimeValue{IsNegative: true, Microsecond: 999999}, nil},
		{"-00:00:00.9999991", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:00:00.9999991")},
		{"-00:00:00.-01", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:00:00.-01")},
		{"-00:00:01", TimeValue{IsNegative: true, Second: 1}, nil},
		{"-00:00:59", TimeValue{IsNegative: true, Second: 59}, nil},
		{"-00:00:60", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:00:60")},
		{"-00:00:-01", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:00:-01")},
		{"-00:01:00", TimeValue{IsNegative: true, Minute: 1}, nil},
		{"-00:59:00", TimeValue{IsNegative: true, Minute: 59}, nil},
		{"-00:60:00", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:60:00")},
		{"-00:-01:00", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:-01:00")},
		{"-00:1.5:00", TimeValue{}, fmt.Errorf(errTimeTpl, "-00:1.5:00")},
		{"-01:00:00", TimeValue{IsNegative: true, Hour: 1}, nil},
		{"-23:00:00", TimeValue{IsNegative: true, Hour: 23}, nil},
		{"-24:00:00", TimeValue{IsNegative: true, Day: 1, Hour: 0}, nil},
		{"-25:00:00", TimeValue{IsNegative: true, Day: 1, Hour: 1}, nil},
		{"-48:00:00", TimeValue{IsNegative: true, Day: 2, Hour: 0}, nil},
		{"-1.1:00:00", TimeValue{}, fmt.Errorf(errTimeTpl, "-1.1:00:00")},
	}
	for _, test := range tests {
		t.Run(test.timestr, func(t *testing.T) {
			v, err := stringToMysqlTime(test.timestr)
			if err != nil {
				if test.err == nil {
					t.Errorf("expect no error, actual err: %v", err)
					t.FailNow()
				}
				if err.Error() != test.err.Error() {
					t.Errorf("error not equal, expect: %v, actual: %v", test.err, err)
					t.FailNow()
				}
			} else {
				if test.err != nil {
					t.Errorf("expect error: %v, actual no error, value: %v", test.err, v)
					t.FailNow()
				}

				if v != test.value {
					t.Errorf("result not equal, str: %s, expect: %v, actual: %v", test.timestr, test.value, v)
					t.FailNow()
				}
			}
		})
	}
}
