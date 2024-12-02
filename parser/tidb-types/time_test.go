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
	"fmt"
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

func TestString(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	expected := "2023-10-01 12:30:45"
	assert.Equal(t, expected, t1.String())

	t2 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 123456),
		Type: mysql.TypeDatetime,
		Fsp:  6,
	}

	expectedWithFsp := "2023-10-01 12:30:45.123456"
	assert.Equal(t, expectedWithFsp, t2.String())
}

func TestConvert(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local
	converted, err := t1.Convert(sc, mysql.TypeTimestamp)
	assert.NoError(t, err)
	assert.Equal(t, mysql.TypeTimestamp, converted.Type)
	assert.Equal(t, t1.Time, converted.Time)
	assert.Equal(t, t1.Fsp, converted.Fsp)
}

func TestConvertToDuration(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	duration, err := t1.ConvertToDuration()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(12*3600+30*60+45)*time.Second, duration.Duration)
	assert.Equal(t, 0, duration.Fsp)
}

func TestCompareString(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	tm := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	// 比较成功
	result, err := tm.CompareString(sc, "2023-10-1 12:30:45")
	assert.NoError(t, err)
	assert.Equal(t, 0, result)

	// 比较失败
	_, err = tm.CompareString(sc, "invalid-time")
	assert.Error(t, err)
}

func TestGetFsp(t *testing.T) {
	// 没有小数部分
	assert.Equal(t, 0, types.GetFsp("2023-10-1 12:30:45"))

	// 有小数部分
	assert.Equal(t, 3, types.GetFsp("2023-10-1 12:30:45.123"))

	// 小数部分超过6位
	assert.Equal(t, 6, types.GetFsp("2023-10-1 12:30:45.123456789"))
}

func TestToPackedUint(t *testing.T) {
	tm := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 500000),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	_, err := tm.ToPackedUint()
	assert.NoError(t, err)
}

func TestFromPackedUint(t *testing.T) {
	var tm types.Time
	tm.FromPackedUint(0x07E70A010C1E2)
}

func TestTimestampDiff(t *testing.T) {
	t1 := types.Time{
		Time: types.FromDate(2023, 10, 1, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}
	t2 := types.Time{
		Time: types.FromDate(2023, 10, 2, 12, 30, 45, 0),
		Type: mysql.TypeDatetime,
		Fsp:  0,
	}

	// 测试 "DAY" 单位
	assert.Equal(t, int64(1), types.TimestampDiff("DAY", t1, t2))

	// 测试 "HOUR" 单位
	assert.Equal(t, int64(24), types.TimestampDiff("HOUR", t1, t2))

	// 测试 "MINUTE" 单位
	assert.Equal(t, int64(1440), types.TimestampDiff("MINUTE", t1, t2))

	// 测试 "SECOND" 单位
	assert.Equal(t, int64(86400), types.TimestampDiff("SECOND", t1, t2))
}

func TestParseDateFormat(t *testing.T) {
	// 测试有效的日期格式
	assert.Equal(t, []string{"2023", "10", "01"}, types.ParseDateFormat("2023-10-01"))
}

func TestParseDatetime(t *testing.T) {
	sc := &stmtctx.StatementContext{TimeZone: time.UTC}

	// 测试有效的日期时间格式
	types.ParseDatetime(sc, "2023-10-01 12:30:45.123456")

	// 测试无效的日期时间格式
	types.ParseDatetime(sc, "2023-10-01 12:30:45.1234567")
}

func TestParseYear(t *testing.T) {
	// 测试有效的年份格式
	_, err := types.ParseYear("2023")
	assert.NoError(t, err)

	// 测试无效的年份格式
	_, err = types.ParseYear("2023a")
}

func TestAdjustYear(t *testing.T) {
	// 测试调整年份
	types.AdjustYear(0, true)
	types.AdjustYear(69, false)
	types.AdjustYear(69, false)
	types.AdjustYear(70, false)
	types.AdjustYear(99, false)
}

func TestDurationAdd(t *testing.T) {
	d1 := types.Duration{Duration: time.Hour, Fsp: 0}
	d2 := types.Duration{Duration: time.Hour * 2, Fsp: 0}

	result, err := d1.Add(d2)
	assert.NoError(t, err)
	assert.Equal(t, time.Hour*3, result.Duration)
	assert.Equal(t, 0, result.Fsp)
}

func TestDurationSub(t *testing.T) {
	d1 := types.Duration{Duration: time.Hour * 3, Fsp: 0}
	d2 := types.Duration{Duration: time.Hour * 2, Fsp: 0}

	result, err := d1.Sub(d2)
	assert.NoError(t, err)
	assert.Equal(t, time.Hour, result.Duration)
	assert.Equal(t, 0, result.Fsp)
}

func TestDurationString(t *testing.T) {
	d := types.Duration{Duration: time.Hour + time.Minute + time.Second, Fsp: 0}
	assert.Equal(t, "01:01:01", d.String())
}

func TestDurationToNumber(t *testing.T) {
	d := types.Duration{Duration: time.Hour + time.Minute + time.Second, Fsp: 0}
	d.ToNumber()
}

func TestDurationConvertToTime(t *testing.T) {
	sc := &stmtctx.StatementContext{TimeZone: time.UTC}
	d := types.Duration{Duration: time.Hour + time.Minute + time.Second, Fsp: 0}

	_, err := d.ConvertToTime(sc, mysql.TypeDatetime)
	assert.NoError(t, err)
}

func TestDurationRoundFrac(t *testing.T) {
	d := types.Duration{Duration: time.Hour + time.Minute + time.Second + 500*time.Millisecond, Fsp: 0}

	d.RoundFrac(1)
}

func TestDurationCompare(t *testing.T) {
	d1 := types.Duration{Duration: time.Hour, Fsp: 0}
	d2 := types.Duration{Duration: time.Hour * 2, Fsp: 0}

	assert.Equal(t, -1, d1.Compare(d2))
	assert.Equal(t, 0, d1.Compare(d1))
	assert.Equal(t, 1, d2.Compare(d1))
}

func TestDurationCompareString(t *testing.T) {
	sc := &stmtctx.StatementContext{TimeZone: time.UTC}
	d1 := types.Duration{Duration: time.Hour, Fsp: 0}

	result, err := d1.CompareString(sc, "02:00:00")
	assert.NoError(t, err)
	assert.Equal(t, -1, result)
}

func TestDurationHour(t *testing.T) {
	d := types.Duration{Duration: time.Hour + time.Minute + time.Second, Fsp: 0}
	assert.Equal(t, 1, d.Hour())
}

// 测试 Minute 方法
func TestDurationMinute(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     int
	}{
		{"整点分钟", 11 * time.Minute, 11},
		{"非整点分钟", 11*time.Minute + 30*time.Second, 11},
		{"零分钟", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := types.Duration{Duration: tt.duration}
			if got := d.Minute(); got != tt.want {
				t.Errorf("Duration.Minute() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 Second 方法
func TestDurationSecond(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     int
	}{
		{"整点秒", 11 * time.Second, 11},
		{"非整点秒", 11*time.Second + 500*time.Millisecond, 11},
		{"零秒", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := types.Duration{Duration: tt.duration}
			if got := d.Second(); got != tt.want {
				t.Errorf("Duration.Second() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 MicroSecond 方法
func TestDurationMicroSecond(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     int
	}{
		{"整点微秒", 11 * time.Microsecond, 11},
		{"非整点微秒", 11*time.Microsecond + 500*time.Nanosecond, 11},
		{"零微秒", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := types.Duration{Duration: tt.duration}
			if got := d.MicroSecond(); got != tt.want {
				t.Errorf("Duration.MicroSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	// 创建一个 StatementContext 用于测试
	sc := &stmtctx.StatementContext{}

	// 测试用例1：正常情况，解析一个标准的持续时间字符串
	t.Run("Normal Case", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "12:34:56.789", 3)
		assert.NoError(t, err)
		assert.Equal(t, types.Duration{Duration: 45296789000000, Fsp: 3}, duration)
	})

	// 测试用例2：包含天数的持续时间字符串
	t.Run("With Day", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "1 12:34:56.789", 3)
		assert.NoError(t, err)
		assert.Equal(t, types.Duration{Duration: 131696789000000, Fsp: 3}, duration)
	})

	// 测试用例3：负的持续时间字符串
	t.Run("Negative Duration", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "-12:34:56.789", 3)
		assert.NoError(t, err)
		assert.Equal(t, types.Duration{Duration: -45296789000000, Fsp: 3}, duration)
	})

	// 测试用例4：空字符串
	t.Run("Empty String", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "", 3)
		assert.NoError(t, err)
		assert.Equal(t, types.ZeroDuration, duration)
	})

	// 测试用例5：无效的持续时间字符串
	t.Run("Invalid Duration", func(t *testing.T) {
		_, err := types.ParseDuration(sc, "invalid", 3)
		assert.Error(t, err)
	})

	// 测试用例6：超出范围的持续时间字符串
	t.Run("Overflow Duration", func(t *testing.T) {
		_, err := types.ParseDuration(sc, "999999:59:59.999999", 6)
		assert.Error(t, err)
	})

	// 测试用例7：检查 Fsp 参数的有效性
	t.Run("Invalid Fsp", func(t *testing.T) {
		_, err := types.ParseDuration(sc, "12:34:56.789", 7)
		assert.Error(t, err)
	})

	// 测试用例8：没有小数部分的持续时间字符串
	t.Run("No Fractional Part", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "12:34:56", 0)
		assert.NoError(t, err)
		assert.Equal(t, types.Duration{Duration: 45296000000000, Fsp: 0}, duration)
	})

	// 测试用例9：只有小时和分钟的持续时间字符串
	t.Run("Only Hours and Minutes", func(t *testing.T) {
		duration, err := types.ParseDuration(sc, "12:34", 0)
		assert.NoError(t, err)
		assert.Equal(t, types.Duration{Duration: 45240000000000, Fsp: 0}, duration)
	})

	// 测试用例10：只有分钟的持续时间字符串
	t.Run("Only Minutes", func(t *testing.T) {
		_, err := types.ParseDuration(sc, "34", 0)
		assert.NoError(t, err)
	})
}

func TestParseTime(t *testing.T) {
	sc := &stmtctx.StatementContext{}

	// 测试正常情况
	timeStr := "2023-10-01 12:00:00"
	tp := mysql.TypeDatetime
	fsp := 0
	result, err := types.ParseTime(sc, timeStr, tp, fsp)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	timeStr = "invalid time"
	_, err = types.ParseTime(sc, timeStr, tp, fsp)
	assert.Error(t, err)
}

func TestParseTimeFromFloatString(t *testing.T) {
	sc := &stmtctx.StatementContext{}

	// 测试正常情况
	timeStr := "20231001120000"
	tp := mysql.TypeDatetime
	fsp := 0
	result, err := types.ParseTimeFromFloatString(sc, timeStr, tp, fsp)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	timeStr = "invalid time"
	_, err = types.ParseTimeFromFloatString(sc, timeStr, tp, fsp)
	assert.Error(t, err)
}

func TestParseTimestamp(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local

	// 测试正常情况
	timeStr := "2023-10-01 12:00:00"
	result, err := types.ParseTimestamp(sc, timeStr)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	timeStr = "invalid time"
	_, err = types.ParseTimestamp(sc, timeStr)
	assert.Error(t, err)
}

func TestParseDate(t *testing.T) {
	sc := &stmtctx.StatementContext{}

	// 测试正常情况
	timeStr := "2023-10-01"
	result, err := types.ParseDate(sc, timeStr)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01", result.String())

	// 测试错误情况
	timeStr = "invalid time"
	_, err = types.ParseDate(sc, timeStr)
	assert.Error(t, err)
}

func TestParseTimeFromNum(t *testing.T) {
	sc := &stmtctx.StatementContext{}

	// 测试正常情况
	num := int64(20231001120000)
	tp := mysql.TypeDatetime
	fsp := 0
	result, err := types.ParseTimeFromNum(sc, num, tp, fsp)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	num = int64(-1)
	_, err = types.ParseTimeFromNum(sc, num, tp, fsp)
	assert.Error(t, err)
}

func TestParseDatetimeFromNum(t *testing.T) {
	sc := &stmtctx.StatementContext{}

	// 测试正常情况
	num := int64(20231001120000)
	result, err := types.ParseDatetimeFromNum(sc, num)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	num = int64(-1)
	_, err = types.ParseDatetimeFromNum(sc, num)
	assert.Error(t, err)
}

func TestParseTimestampFromNum(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local

	// 测试正常情况
	num := int64(20231001120000)
	result, err := types.ParseTimestampFromNum(sc, num)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01 12:00:00", result.String())

	// 测试错误情况
	num = int64(-1)
	_, err = types.ParseTimestampFromNum(sc, num)
	assert.Error(t, err)
}

func TestParseDateFromNum(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local

	// 测试正常情况
	num := int64(20231001)
	result, err := types.ParseDateFromNum(sc, num)
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01", result.String())

	// 测试错误情况
	num = int64(-1)
	_, err = types.ParseDateFromNum(sc, num)
	assert.Error(t, err)
}

func TestTimeFromDays(t *testing.T) {
	// 测试正常情况
	num := int64(738427) // 2023-10-01
	result := types.TimeFromDays(num)

	// 测试错误情况
	num = int64(-1)
	result = types.TimeFromDays(num)
	assert.Equal(t, "0000-00-00", result.String())
}

func TestParseTimeFromInt64(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local
	tests := []struct {
		num      int64
		expected types.Time
	}{
		{20231001, types.Time{Time: types.FromDate(2023, 10, 1, 0, 0, 0, 0)}},
		{20231001150405, types.Time{Time: types.FromDate(2023, 10, 1, 15, 4, 5, 0)}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.num), func(t *testing.T) {
			result, err := types.ParseTimeFromInt64(sc, test.num)
			assert.NoError(t, err)
			assert.Equal(t, test.expected.Time, result.Time)
		})
	}
}

func TestStrToDate(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	sc.TimeZone = time.Local
	tests := []struct {
		date     string
		format   string
		expected bool
	}{
		{"2023-10-01", "%Y-%m-%d", true},
		{"2023-10-01 15:04:05", "%Y-%m-%d %H:%i:%s", true},
		{"2023-10-01", "%Y-%m-%d %H:%i:%s", false},
		{"2023-10-01 15:04:05", "%Y-%m-%d", false},
	}

	for _, test := range tests {
		t.Run(test.date, func(t *testing.T) {
			var time types.Time
			time.StrToDate(sc, test.date, test.format)
		})
	}
}
