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

package mysql

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/util/hack"
)

// FieldData means filed data, is []byte
type FieldData []byte

// Field to represent column field
type Field struct {
	Data         FieldData
	Schema       []byte
	Table        []byte
	OrgTable     []byte
	Name         []byte
	OrgName      []byte
	Charset      uint16
	ColumnLength uint32
	Type         uint8
	Flag         uint16
	Decimal      uint8

	DefaultValueLength uint64
	DefaultValue       []byte
}

// TimeValue mysql time value
type TimeValue struct {
	IsNegative  bool
	Day         int
	Hour        int
	Minute      int
	Second      int
	Microsecond int
}

// IsNull check TimeValue if null
func (m *TimeValue) IsNull() bool {
	return m.Day == 0 && m.Hour == 0 && m.Minute == 0 && m.Second == 0 && m.Microsecond == 0
}

// Parse parse []byte to Field
func (p FieldData) Parse() (f *Field, err error) {
	f = new(Field)

	data := make([]byte, len(p))
	copy(data, p)
	f.Data = data

	pos := 0
	ok := false
	//skip catelog, always def
	pos, ok = skipLenEncString(p, pos)
	if !ok {
		return f, errors.New("skipLenEncString in Parse failed")
	}

	//schema
	f.Schema, pos, _, ok = ReadLenEncStringAsBytes(p, pos)
	if !ok {
		return f, errors.New("read Schema failed")
	}

	//table
	f.Table, pos, _, ok = ReadLenEncStringAsBytes(p, pos)
	if !ok {
		return f, errors.New("read Table failed")
	}

	//org_table
	f.OrgTable, pos, _, ok = ReadLenEncStringAsBytes(p, pos)
	if !ok {
		return f, errors.New("read OrgTable failed")
	}

	//name
	f.Name, pos, _, ok = ReadLenEncStringAsBytes(p, pos)
	if !ok {
		return f, errors.New("read Name failed")
	}

	//org_name
	f.OrgName, pos, _, ok = ReadLenEncStringAsBytes(p, pos)
	if !ok {
		return f, errors.New("read OrgName failed")
	}

	//skip oc
	pos++

	//charset
	f.Charset = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	//column length
	f.ColumnLength = binary.LittleEndian.Uint32(p[pos:])
	pos += 4

	//type
	f.Type = p[pos]
	pos++

	//flag
	f.Flag = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	//decimals 1
	f.Decimal = p[pos]
	pos++

	//filter [0x00][0x00]
	pos += 2

	f.DefaultValue = nil
	//if more data, command was field list
	if len(p) > pos {
		//length of default value lenenc-int
		f.DefaultValueLength, pos, _, _ = ReadLenEncInt(p, pos)

		if pos+int(f.DefaultValueLength) > len(p) {
			err = ErrMalformPacket
			return
		}

		//default value string[$len]
		f.DefaultValue = p[pos:(pos + int(f.DefaultValueLength))]
	}

	return
}

// Dump dume field into binary []byte
func (f *Field) Dump() []byte {
	if f.Data != nil {
		return []byte(f.Data)
	}

	l := len(f.Schema) + len(f.Table) + len(f.OrgTable) + len(f.Name) + len(f.OrgName) + len(f.DefaultValue) + 48

	data := make([]byte, 0, l)

	data = AppendLenEncStringBytes(data, []byte("def"))
	data = AppendLenEncStringBytes(data, f.Schema)
	data = AppendLenEncStringBytes(data, f.Table)
	data = AppendLenEncStringBytes(data, f.OrgTable)
	data = AppendLenEncStringBytes(data, f.Name)
	data = AppendLenEncStringBytes(data, f.OrgName)

	data = append(data, 0x0c)

	data = AppendUint16(data, f.Charset)
	data = AppendUint32(data, f.ColumnLength)
	data = append(data, f.Type)
	data = AppendUint16(data, f.Flag)
	data = append(data, f.Decimal)
	data = append(data, 0, 0)

	if f.DefaultValue != nil {
		data = AppendUint64(data, f.DefaultValueLength)
		data = append(data, f.DefaultValue...)
	}

	return data
}

// FieldType return type of field
func FieldType(value interface{}) (typ uint8, err error) {
	switch value.(type) {
	case int8, int16, int32, int64, int:
		typ = TypeLonglong
	case uint8, uint16, uint32, uint64, uint:
		typ = TypeLonglong
	case float32, float64:
		typ = TypeDouble
	case string, []byte:
		typ = TypeVarString
	case nil:
		typ = TypeNull
	default:
		err = fmt.Errorf("unsupport type %T for resultset", value)
	}
	return
}

func stringToMysqlTime(s string) (TimeValue, error) {
	var v TimeValue

	timeFields := strings.SplitN(s, ":", 2)
	if len(timeFields) != 2 {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	hour, err := strconv.ParseInt(timeFields[0], 10, 64)
	if err != nil {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	if strings.HasPrefix(timeFields[0], "-") {
		v.IsNegative = true
		hour = hack.Abs(hour)
	}

	day := int(hour / 24)
	hourRest := int(hour % 24)

	timeRest := strconv.Itoa(hourRest) + ":" + timeFields[1]
	ts, err := time.Parse("15:04:05", timeRest)
	if err != nil {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}
	if ts.Nanosecond()%1000 != 0 {
		return v, fmt.Errorf("invalid TypeDuration %s", s)
	}

	v.Day = day
	v.Hour = ts.Hour()
	v.Minute = ts.Minute()
	v.Second = ts.Second()
	v.Microsecond = ts.Nanosecond() / 1000
	return v, nil
}

func mysqlTimeToBinaryResult(v TimeValue) []byte {
	var t []byte
	var length uint8
	if v.IsNull() {
		length = 0
		t = append(t, length)
	} else {
		if v.Microsecond == 0 {
			length = 8
		} else {
			length = 12
		}
		t = append(t, length)
		if v.IsNegative {
			t = append(t, 1)
		} else {
			t = append(t, 0)
		}
		t = AppendUint32(t, uint32(v.Day))
		t = append(t, uint8(v.Hour))
		t = append(t, uint8(v.Minute))
		t = append(t, uint8(v.Second))
		if v.Microsecond != 0 {
			t = AppendUint32(t, uint32(v.Microsecond))
		}
	}
	return t
}
