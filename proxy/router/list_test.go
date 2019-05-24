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

package router

import (
	"testing"

	"github.com/XiaoMi/Gaea/core/errors"
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

func TestParseYearRange(t *testing.T) {
	dateRange := "2014-2017"
	years, err := ParseYearRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}
	testCheckList(t, years, 2014, 2015, 2016, 2017)

	dateRange = "2017-2013"
	years, err = ParseYearRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}
	testCheckList(t, years, 2013, 2014, 2015, 2016, 2017)

	dateRange = "20120"
	years, err = ParseYearRange(dateRange)
	if err != errors.ErrDateRangeIllegal || years != nil {
		t.Fatal(err)
	}

	dateRange = "2o12"
	years, err = ParseYearRange(dateRange)
	if err == nil || years != nil {
		t.Failed()
	}

	dateRange = "2012"
	years, err = ParseYearRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}
	testCheckList(t, years, 2012)
}

func TestParseMonthRange(t *testing.T) {
	dateRange := "201602-201610"
	months, err := ParseMonthRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}

	testCheckList(t, months,
		201602,
		201603,
		201604,
		201605,
		201606,
		201607,
		201608,
		201609,
		201610,
	)

	dateRange = "201603-201511"
	months, err = ParseMonthRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}

	testCheckList(t, months,
		201511,
		201512,
		201601,
		201602,
		201603,
	)

	dateRange = "20120"
	months, err = ParseMonthRange(dateRange)
	if err != errors.ErrDateRangeIllegal || months != nil {
		t.Fatal(err)
	}

	dateRange = "2012o1"
	months, err = ParseMonthRange(dateRange)
	if err == nil || months != nil {
		t.Failed()
	}

	dateRange = "201201"
	months, err = ParseMonthRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}
	testCheckList(t, months, 201201)
}

func TestParseDayRange(t *testing.T) {
	dateRange := "20160227-20160304"
	days, err := ParseDayRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}

	testCheckList(t, days,
		20160227,
		20160228,
		20160229,
		20160301,
		20160302,
		20160303,
		20160304,
	)

	dateRange = "20160304-20160301"
	days, err = ParseDayRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}

	testCheckList(t, days,
		20160301,
		20160302,
		20160303,
		20160304,
	)

	dateRange = "2016034"
	days, err = ParseDayRange(dateRange)
	if err != errors.ErrDateRangeIllegal || days != nil {
		t.Fatal(err)
	}

	dateRange = "201603o4"
	days, err = ParseDayRange(dateRange)
	if err == nil {
		t.Failed()
	}

	dateRange = "20160304"
	days, err = ParseDayRange(dateRange)
	if err != nil {
		t.Fatal(err)
	}
	testCheckList(t, days, 20160304)
}
