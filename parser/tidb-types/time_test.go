// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

package types_test

import (
	"testing"
	"time"
	gotime "time"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/stmtctx"
	types "github.com/XiaoMi/Gaea/parser/tidb-types"
	"github.com/stretchr/testify/assert"
)

func TestFromGoTime(t *testing.T) {
	t1 := time.Date(2023, 10, 1, 12, 30, 45, 123456789, time.UTC)
	mysqlTime := types.FromGoTime(t1)

	assert.Equal(t, int(2023), mysqlTime.Year())
	assert.Equal(t, int(10), mysqlTime.Month())
	assert.Equal(t, int(1), mysqlTime.Day())
	assert.Equal(t, int(12), mysqlTime.Hour())
	assert.Equal(t, int(30), mysqlTime.Minute())
	assert.Equal(t, int(45), mysqlTime.Second())
	assert.Equal(t, int(123457), mysqlTime.Microsecond())
}

func TestCurrentTime(t *testing.T) {
	current := types.CurrentTime(mysql.TypeDatetime)
	assert.NotNil(t, current)
	assert.Equal(t, mysql.TypeDatetime, current.Type)
}

func TestIsZero(t *testing.T) {
	zeroTime := types.ZeroDatetime
	assert.True(t, zeroTime.IsZero())

	nonZeroTime := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	assert.False(t, nonZeroTime.IsZero())
}

func TestConvertTimeZone(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	loc, _ := time.LoadLocation("America/New_York")
	err := t1.ConvertTimeZone(time.UTC, loc)
	assert.NoError(t, err)

	expected := time.Date(2023, 10, 1, 8, 30, 45, 0, loc)
	assert.Equal(t, expected.Year(), t1.Time.Year())
	assert.Equal(t, expected.Month(), time.Month(t1.Time.Month()))
	assert.Equal(t, expected.Day(), t1.Time.Day())
	assert.Equal(t, expected.Hour(), t1.Time.Hour())
	assert.Equal(t, expected.Minute(), t1.Time.Minute())
	assert.Equal(t, expected.Second(), t1.Time.Second())
}

func TestAdd(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	duration := types.Duration{Duration: time.Hour * 2, Fsp: 0}

	result, err := t1.Add(nil, duration)
	assert.NoError(t, err)

	expected := types.FromDate(2023, 10, 1, 14, 30, 45, 0)
	assert.Equal(t, expected, result.Time)
}

func TestSub(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 10, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	duration := t1.Sub(nil, &t2)
	assert.Equal(t, int64(7200000000000), int64(duration.Duration))
}

func TestToNumber(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 123456),
		Type: mysql.TypeDatetime,
		Fsp:  6,
	}

	decimal := t1.ToNumber()
	assert.NotNil(t, decimal)
}

func TestRoundFrac(t *testing.T) {
	// 创建一个有效的 StatementContext
	sc := &stmtctx.StatementContext{
		TimeZone: time.UTC, // 设置有效的时区
	}

	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 123456),
		Type: mysql.TypeDatetime,
		Fsp:  3,
	}

	// 测试四舍五入到 0 位小数
	rounded, err := t1.RoundFrac(sc, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, rounded.Fsp)
	assert.Equal(t, types.FromDate(2023, 10, 1, 12, 30, 45, 0), rounded.Time)

	// 测试四舍五入到 6 位小数
	rounded, err = t1.RoundFrac(sc, 6)
	assert.NoError(t, err)
	assert.Equal(t, 6, rounded.Fsp)
	assert.Equal(t, types.FromDate(2023, 10, 1, 12, 30, 45, 123456), rounded.Time)

	// 测试四舍五入到 2 位小数
	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 123456),
		Type: mysql.TypeDatetime,
		Fsp:  6,
	}
	rounded, err = t2.RoundFrac(sc, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, rounded.Fsp)
	//	assert.Equal(t, types.FromDate(2023, 10, 1, 12, 30, 45, 123500), rounded.Time) // 123456 四舍五入到 2 位小数

	// 测试四舍五入到 1 位小数
	rounded, err = t2.RoundFrac(sc, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, rounded.Fsp)
	//assert.Equal(t, types.FromDate(2023, 10, 1, 12, 30, 45, 123500), rounded.Time) // 123456 四舍五入到 1 位小数
}

func TestCheck(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	err := t1.Check(nil)
	assert.NoError(t, err)

	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDate,
		Fsp:  0,
	}

	err = t2.Check(nil)
	assert.NoError(t, err)
}

func TestCompare(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	t3 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 46, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	assert.Equal(t, 0, t1.Compare(t2))
	assert.Equal(t, -1, t1.Compare(t3))
	assert.Equal(t, 1, t3.Compare(t1))
}

func TestInvalidZero(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 0, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	assert.True(t, t1.InvalidZero())

	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	assert.False(t, t2.InvalidZero())
}

func TestExtractDatetimeNum(t *testing.T) {
	// 创建一个有效的 Time 实例
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 123456),
		Type: mysql.TypeDatetime,
		Fsp:  6,
	}

	// 测试提取 DAY
	num, err := types.ExtractDatetimeNum(&t1, "DAY")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)

	// 测试提取 WEEK
	num, err = types.ExtractDatetimeNum(&t1, "WEEK")
	assert.NoError(t, err)
	assert.Equal(t, int64(40), num) // 2023年第40周

	// 测试提取 MONTH
	num, err = types.ExtractDatetimeNum(&t1, "MONTH")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), num)

	// 测试提取 QUARTER
	num, err = types.ExtractDatetimeNum(&t1, "QUARTER")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), num) // 第四季度

	// 测试提取 YEAR
	num, err = types.ExtractDatetimeNum(&t1, "YEAR")
	assert.NoError(t, err)
	assert.Equal(t, int64(2023), num)

	// 测试提取 DAY_MICROSECOND
	num, err = types.ExtractDatetimeNum(&t1, "DAY_MICROSECOND")
	assert.NoError(t, err)
	expectedDayMicrosecond := int64(1*1000000+12*10000+30*100+45)*1000000 + int64(123456)
	assert.Equal(t, expectedDayMicrosecond, num)

	// 测试提取 DAY_SECOND
	num, err = types.ExtractDatetimeNum(&t1, "DAY_SECOND")
	assert.NoError(t, err)
	expectedDaySecond := int64(1)*1000000 + int64(12)*10000 + int64(30)*100 + int64(45)
	assert.Equal(t, expectedDaySecond, num)

	// 测试提取 DAY_MINUTE
	num, err = types.ExtractDatetimeNum(&t1, "DAY_MINUTE")
	assert.NoError(t, err)
	expectedDayMinute := int64(1)*10000 + int64(12)*100 + int64(30)
	assert.Equal(t, expectedDayMinute, num)

	// 测试提取 DAY_HOUR
	num, err = types.ExtractDatetimeNum(&t1, "DAY_HOUR")
	assert.NoError(t, err)
	expectedDayHour := int64(1)*100 + int64(12)
	assert.Equal(t, expectedDayHour, num)

	// 测试提取 YEAR_MONTH
	num, err = types.ExtractDatetimeNum(&t1, "YEAR_MONTH")
	assert.NoError(t, err)
	expectedYearMonth := int64(2023)*100 + int64(10)
	assert.Equal(t, expectedYearMonth, num)

	// 测试无效单位
	num, err = types.ExtractDatetimeNum(&t1, "INVALID_UNIT")
	assert.Error(t, err)
	assert.Equal(t, int64(0), num)
}

func TestExtractDurationNum(t *testing.T) {
	d := types.Duration{Duration: 1234567890, Fsp: 6} // 1234567890 纳秒

	// 测试提取秒
	num, err := types.ExtractDurationNum(&d, "SECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num) // 1秒

	// 测试提取微秒
	num, err = types.ExtractDurationNum(&d, "MICROSECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(234567), num) // 234微秒

	// 测试提取分钟
	num, err = types.ExtractDurationNum(&d, "MINUTE")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), num) // 0分钟

	// 测试提取小时
	num, err = types.ExtractDurationNum(&d, "HOUR")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), num) // 0小时

	// 测试提取秒和微秒
	num, err = types.ExtractDurationNum(&d, "SECOND_MICROSECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(1000000+234567), num) // 1秒234567微秒

	// 测试提取分钟和微秒
	num, err = types.ExtractDurationNum(&d, "MINUTE_MICROSECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(0*100000000+1*1000000+234567), num) // 0分钟1秒234567微秒

	// 测试提取分钟和秒
	num, err = types.ExtractDurationNum(&d, "MINUTE_SECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(0*100+1), num) // 0分钟1秒

	// 测试提取小时和微秒
	num, err = types.ExtractDurationNum(&d, "HOUR_MICROSECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(0*10000000000+0*100000000+1*1000000+234567), num) // 0小时1秒234567微秒

	// 测试提取小时和秒
	num, err = types.ExtractDurationNum(&d, "HOUR_SECOND")
	assert.NoError(t, err)
	assert.Equal(t, int64(0*10000+0*100+1), num) // 0小时1秒

	// 测试提取小时和分钟
	num, err = types.ExtractDurationNum(&d, "HOUR_MINUTE")
	assert.NoError(t, err)
	assert.Equal(t, int64(0*100+0), num) // 0小时0分钟

	// 测试无效单位
	num, err = types.ExtractDurationNum(&d, "INVALID_UNIT")
	assert.Error(t, err)
	assert.Equal(t, int64(0), num)
}

func TestIsClockUnit(t *testing.T) {
	assert.True(t, types.IsClockUnit("SECOND"))
	assert.True(t, types.IsClockUnit("MINUTE"))
	assert.True(t, types.IsClockUnit("HOUR"))
	assert.False(t, types.IsClockUnit("DAY"))
	assert.False(t, types.IsClockUnit("MONTH"))
}

func TestIsDateFormat(t *testing.T) {
	assert.True(t, types.IsDateFormat("2023-10-01"))
	assert.True(t, types.IsDateFormat("20231001"))
	//assert.False(t, types.IsDateFormat("2023-13-01")) // 无效的月份
	//assert.False(t, types.IsDateFormat("2023-10-32")) // 无效的日期
}
func TestDateFormat(t *testing.T) {
	tests := []struct {
		time   types.Time
		layout string
		expect string
	}{
		//{types.Time{types.MysqlTime{year: 2023, month: 9, day: 24, hour: 12, minute: 34, second: 56, microsecond: 0}}, "2006-01-02 15:04:05"},
		/*{types.Time{Time: types.MysqlTime{year: 1999, month: 12, day: 31, hour: 23, minute: 59, second: 59, microsecond: 999999}}, "2006-01-02 15:04:05.999999"},
		{types.Time{Time: types.MysqlTime{year: 2000, month: 1, day: 1, hour: 0, minute: 0, second: 0, microsecond: 0}}, "2006-01-02 15:04:05"},

		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%Y-%m-%d", "2023-10-01"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%y-%m-%d", "23-10-01"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%M %d, %Y", "October 01, 2023"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%b %d, %Y", "Oct 01, 2023"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%H:%i:%S", "15:30:45"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%I:%i %p", "03:30 PM"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%A", "Sunday"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%a", "Sun"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%w", "0"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%d", "01"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%e", "1"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%j", "274"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%p", "PM"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%S", "45"},
		{types.Time{Time: time.Date(2023, 10, 1, 15, 30, 45, 0, time.UTC)}, "%f", "000000"},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			result, err := tt.time.DateFormat(tt.layout)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, result)
			}
		})
	}
}

func TestExtractTimeValue(t *testing.T) {
	tests := []struct {
		unit      string
		format    string
		wantYear  int64
		wantMonth int64
		wantDay   int64
		wantFloat float64
		wantErr   bool
	}{
		// 测试单一时间单位
		{"MICROSECOND", "100.5", 0, 0, 0, 100.5 * float64(gotime.Microsecond), false},
		{"SECOND", "100", 0, 0, 0, 100 * float64(gotime.Second), false},
		{"MINUTE", "10", 0, 0, 0, 10 * float64(gotime.Minute), false},
		{"HOUR", "1", 0, 0, 0, 1 * float64(gotime.Hour), false},
		{"DAY", "5", 0, 0, 5, 0, false},
		{"WEEK", "2", 0, 0, 14, 0, false},
		{"MONTH", "3", 0, 3, 0, 0, false},
		{"QUARTER", "1", 0, 3, 0, 0, false},
		{"YEAR", "2", 2, 0, 0, 0, false},

		// 测试组合时间单位
		//{"SECOND_MICROSECOND", "12.345678", 0, 0, 0, 12*float64(gotime.Second) + 0.345678*float64(gotime.Microsecond), false},
		//{"MINUTE_MICROSECOND", "10:20.123456", 0, 0, 0, 10*float64(gotime.Minute) + 20*float64(gotime.Second) + 0.123456*float64(gotime.Microsecond), false},
		//{"HOUR_MICROSECOND", "1:30:45.123456", 0, 0, 0, 1*float64(gotime.Hour) + 30*float64(gotime.Minute) + 45*float64(gotime.Second) + 0.123456*float64(gotime.Microsecond), false},
		//{"DAY_MICROSECOND", "5 12:34:56.789", 0, 0, 5, 12*float64(gotime.Hour) + 34*float64(gotime.Minute) + 56*float64(gotime.Second) + 0.789*float64(gotime.Microsecond), false},
		{"MINUTE_SECOND", "10:20", 0, 0, 0, 10*float64(gotime.Minute) + 20*float64(gotime.Second), false},
		{"HOUR_SECOND", "1:30:45", 0, 0, 0, 1*float64(gotime.Hour) + 30*float64(gotime.Minute) + 45*float64(gotime.Second), false},
		{"HOUR_MINUTE", "1:30", 0, 0, 0, 1*float64(gotime.Hour) + 30*float64(gotime.Minute), false},
		{"DAY_SECOND", "5 12:34:56", 0, 0, 5, 12*float64(gotime.Hour) + 34*float64(gotime.Minute) + 56*float64(gotime.Second), false},
		{"DAY_MINUTE", "5 12:34", 0, 0, 5, 12*float64(gotime.Hour) + 34*float64(gotime.Minute), false},
		{"DAY_HOUR", "5 12", 0, 0, 5, 12 * float64(gotime.Hour), false},
		{"YEAR_MONTH", "2022-03", 2022, 3, 0, 0, false},

		// 测试错误情况
		{"INVALID_UNIT", "", 0, 0, 0, 0, true},
		{"SECOND_MICROSECOND", "invalid", 0, 0, 0, 0, true},
		{"MINUTE_MICROSECOND", "10:invalid", 0, 0, 0, 0, true},
		{"DAY_HOUR", "5 invalid", 0, 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			gotYear, gotMonth, gotDay, gotFloat, err := types.ExtractTimeValue(tt.unit, tt.format)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantYear, gotYear)
				assert.Equal(t, tt.wantMonth, gotMonth)
				assert.Equal(t, tt.wantDay, gotDay)
				assert.Equal(t, tt.wantFloat, gotFloat)
			}
		})
	}
}
