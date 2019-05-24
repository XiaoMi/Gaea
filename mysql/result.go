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
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/util/hack"
)

// RowData row in []byte format
type RowData []byte

// Parse parse data to field
func (p RowData) Parse(f []*Field, binary bool) ([]interface{}, error) {
	if binary {
		return p.ParseBinary(f)
	}
	return p.ParseText(f)
}

// ParseText parse text format data
func (p RowData) ParseText(f []*Field) ([]interface{}, error) {
	data := make([]interface{}, len(f))

	var err error
	var v []byte
	var isNull, isUnsigned bool
	var pos = 0
	var ok = false

	for i := range f {
		v, pos, isNull, ok = ReadLenEncStringAsBytes(p, pos)
		if !ok {
			return nil, fmt.Errorf("ReadLenEncStringAsBytes in ParseText failed")
		}

		if isNull {
			data[i] = nil
		} else {
			isUnsigned = (f[i].Flag&uint16(UnsignedFlag) > 0)
			switch f[i].Type {
			case TypeTiny, TypeShort, TypeLong, TypeInt24,
				TypeLonglong, TypeYear:
				if isUnsigned {
					data[i], err = strconv.ParseUint(string(v), 10, 64)
				} else {
					data[i], err = strconv.ParseInt(string(v), 10, 64)
				}
			case TypeFloat, TypeDouble, TypeNewDecimal:
				data[i], err = strconv.ParseFloat(string(v), 64)
			case TypeVarchar, TypeVarString,
				TypeString, TypeDatetime,
				TypeDate, TypeDuration, TypeTimestamp:
				data[i] = string(v)
			default:
				data[i] = v
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return data, nil
}

// ParseBinary parse binary format data
func (p RowData) ParseBinary(f []*Field) ([]interface{}, error) {
	data := make([]interface{}, len(f))

	if p[0] != OKHeader {
		return nil, ErrMalformPacket
	}

	pos := 1 + ((len(f) + 7 + 2) >> 3)

	nullBitmap := p[1:pos]

	var isUnsigned bool
	var isNull bool
	var err error
	var v []byte
	for i := range data {
		if nullBitmap[(i+2)/8]&(1<<(uint(i+2)%8)) > 0 {
			data[i] = nil
			continue
		}

		isUnsigned = f[i].Flag&uint16(UnsignedFlag) > 0

		switch f[i].Type {
		case TypeNull:
			data[i] = nil
			continue

		case TypeTiny:
			if isUnsigned {
				data[i] = uint64(p[pos])
			} else {
				data[i] = int64(p[pos])
			}
			pos++
			continue

		case TypeShort, TypeYear:
			if isUnsigned {
				data[i] = uint64(binary.LittleEndian.Uint16(p[pos : pos+2]))
			} else {
				var n int16
				err = binary.Read(bytes.NewBuffer(p[pos:pos+2]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = int64(n)
			}
			pos += 2
			continue

		case TypeInt24, TypeLong:
			if isUnsigned {
				data[i] = uint64(binary.LittleEndian.Uint32(p[pos : pos+4]))
			} else {
				var n int32
				err = binary.Read(bytes.NewBuffer(p[pos:pos+4]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = int64(n)
			}
			pos += 4
			continue

		case TypeLonglong:
			if isUnsigned {
				data[i] = binary.LittleEndian.Uint64(p[pos : pos+8])
			} else {
				var n int64
				err = binary.Read(bytes.NewBuffer(p[pos:pos+8]), binary.LittleEndian, &n)
				if err != nil {
					return nil, err
				}
				data[i] = int64(n)
			}
			pos += 8
			continue

		case TypeFloat:
			//data[i] = float64(math.Float32frombits(binary.LittleEndian.Uint32(p[pos : pos+4])))
			var n float32
			err = binary.Read(bytes.NewBuffer(p[pos:pos+4]), binary.LittleEndian, &n)
			if err != nil {
				return nil, err
			}
			data[i] = float64(n)
			pos += 4
			continue

		case TypeDouble:
			var n float64
			err = binary.Read(bytes.NewBuffer(p[pos:pos+8]), binary.LittleEndian, &n)
			if err != nil {
				return nil, err
			}
			data[i] = n
			pos += 8
			continue

		case TypeDecimal, TypeNewDecimal, TypeVarchar,
			TypeBit, TypeEnum, TypeSet, TypeTinyBlob,
			TypeMediumBlob, TypeLongBlob, TypeBlob,
			TypeVarString, TypeString, TypeGeometry:
			var ok = false
			v, pos, isNull, ok = ReadLenEncStringAsBytes(p, pos)
			if !ok {
				return nil, fmt.Errorf("ReadLenEncStringAsBytes in ParseBinary failed")
			}

			if !isNull {
				data[i] = v
				continue
			} else {
				data[i] = nil
				continue
			}
		case TypeDate, TypeNewDate:
			var num uint64
			num, pos, isNull, _ = ReadLenEncInt(p, pos)

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryDate(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		case TypeTimestamp, TypeDatetime:
			var num uint64
			num, pos, isNull, _ = ReadLenEncInt(p, pos)

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryDateTime(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		case TypeDuration:
			var num uint64
			num, pos, isNull, _ = ReadLenEncInt(p, pos)

			if isNull {
				data[i] = nil
				continue
			}

			data[i], err = FormatBinaryTime(int(num), p[pos:])
			pos += int(num)

			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("Stmt Unknown FieldType %d %s", f[i].Type, f[i].Name)
		}
	}

	return data, nil
}

// Result means mysql status、results after sql execution
type Result struct {
	Status uint16

	InsertID     uint64
	AffectedRows uint64

	*Resultset
}

// Resultset means mysql results of sql execution, included split table sql
type Resultset struct {
	Fields     []*Field        // columns information
	FieldNames map[string]int  // column information, key: column name value: index in Fields
	Values     [][]interface{} // values after sql handled

	RowDatas []RowData // data will returned
}

// RowNumber return row number of results
func (r *Resultset) RowNumber() int {
	return len(r.Values)
}

// ColumnNumber return column number of results
func (r *Resultset) ColumnNumber() int {
	return len(r.Fields)
}

// GetValue return value in special row and column
func (r *Resultset) GetValue(row, column int) (interface{}, error) {
	if row >= len(r.Values) || row < 0 {
		return nil, fmt.Errorf("invalid row index %d", row)
	}

	if column >= len(r.Fields) || column < 0 {
		return nil, fmt.Errorf("invalid column index %d", column)
	}

	return r.Values[row][column], nil
}

// NameIndex return column index in Fields
func (r *Resultset) NameIndex(name string) (int, error) {
	column, ok := r.FieldNames[name]
	if ok {
		return column, nil
	}
	return 0, fmt.Errorf("invalid field name %s", name)
}

// GetValueByName return value in special row and column
func (r *Resultset) GetValueByName(row int, name string) (interface{}, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return nil, err
	}
	return r.GetValue(row, column)
}

// IsNull check if value in special row and column is Null
func (r *Resultset) IsNull(row, column int) (bool, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return false, err
	}

	return d == nil, nil
}

// IsNullByName check if value in special row and column is Null, but the entry param is column name
func (r *Resultset) IsNullByName(row int, name string) (bool, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return false, err
	}
	return r.IsNull(row, column)
}

// GetUint return value in special row and column in uint64 type
func (r *Resultset) GetUint(row, column int) (uint64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case uint64:
		return v, nil
	case int64:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case []byte:
		return strconv.ParseUint(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// GetUintByName return value in special row and column in uint64 type, but the entry param is name
func (r *Resultset) GetUintByName(row int, name string) (uint64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}
	return r.GetUint(row, column)
}

// GetIntByName return value in special row and column in int64 type, but the entry param is name
func (r *Resultset) GetIntByName(row int, name string) (int64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}
	return r.GetInt(row, column)
}

// GetInt return value in special row and column in int64 type
func (r *Resultset) GetInt(row, column int) (int64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case uint64:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// GetFloat return value in special row and column in float64 type
func (r *Resultset) GetFloat(row, column int) (float64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case float64:
		return v, nil
	case uint64:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// GetFloatByName return value in special row and column in float64 type, but the entry param is name
func (r *Resultset) GetFloatByName(row int, name string) (float64, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return 0, err
	}
	return r.GetFloat(row, column)
}

// GetString return value in special row and column in string type
func (r *Resultset) GetString(row, column int) (string, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return "", err
	}

	switch v := d.(type) {
	case string:
		return v, nil
	case []byte:
		return hack.String(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("data type is %T", v)
	}
}

// GetStringByName return value in special row and column in string type, but the entry param is name
func (r *Resultset) GetStringByName(row int, name string) (string, error) {
	column, err := r.NameIndex(name)
	if err != nil {
		return "", err
	}
	return r.GetString(row, column)
}

// BuildResultset build resultset
func BuildResultset(fields []*Field, names []string, values [][]interface{}) (*Resultset, error) {
	var ExistFields bool
	r := new(Resultset)

	r.Fields = make([]*Field, len(names))
	r.FieldNames = make(map[string]int, len(names))

	//use the field def that get from true database
	if len(fields) != 0 {
		if len(r.Fields) == len(fields) {
			ExistFields = true
		} else {
			return nil, errors.ErrInvalidArgument
		}
	}

	var b []byte
	var err error

	for i, vs := range values {
		if len(vs) != len(r.Fields) {
			return nil, fmt.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		for j, value := range vs {
			// build fields
			if i == 0 {
				if ExistFields {
					r.Fields[j] = fields[j]
					r.FieldNames[string(r.Fields[j].Name)] = j
				} else {
					field := &Field{}
					r.Fields[j] = field
					r.FieldNames[string(r.Fields[j].Name)] = j
					field.Name = hack.Slice(names[j])
					if err = formatField(field, value); err != nil {
						return nil, err
					}
				}

			}
			// build row values
			b, err = formatValue(value)
			if err != nil {
				return nil, err
			}

			row = AppendLenEncStringBytes(row, b)
		}

		r.RowDatas = append(r.RowDatas, row)
	}
	//assign the values to the result
	r.Values = values

	return r, nil
}

// BuildBinaryResultset build binary resultset
// https://dev.mysql.com/doc/internals/en/binary-protocol-resultset.html
func BuildBinaryResultset(fields []*Field, values [][]interface{}) (*Resultset, error) {
	r := new(Resultset)
	r.Fields = make([]*Field, len(fields))
	for i := range fields {
		r.Fields[i] = fields[i]
	}

	bitmapLen := ((len(fields) + 7 + 2) >> 3)
	for i, v := range values {
		if len(v) != len(r.Fields) {
			return nil, fmt.Errorf("row %d has %d columns not equal %d", i, len(v), len(r.Fields))
		}

		var row []byte
		nullBitMap := make([]byte, bitmapLen)
		row = append(row, 0)
		row = append(row, nullBitMap...)
		for j, rowVal := range v {
			if rowVal == nil {
				bytePos := (j + 2) / 8
				bitPos := byte((j + 2) % 8)
				nullBitMap[bytePos] |= 1 << bitPos
				continue
			}

			var err error
			row, err = AppendBinaryValue(row, r.Fields[j].Type, rowVal)
			if err != nil {
				return nil, err
			}
		}
		copy(row[1:], nullBitMap)
		r.RowDatas = append(r.RowDatas, row)
	}

	return r, nil
}

// formatField encode field according to type of value if necessary
func formatField(field *Field, value interface{}) error {
	switch value.(type) {
	case int8, int16, int32, int64, int:
		field.Charset = 63
		field.Type = TypeLonglong
		field.Flag = uint16(BinaryFlag | NotNullFlag)
	case uint8, uint16, uint32, uint64, uint:
		field.Charset = 63
		field.Type = TypeLonglong
		field.Flag = uint16(BinaryFlag | NotNullFlag | UnsignedFlag)
	case float32, float64:
		field.Charset = 63
		field.Type = TypeDouble
		field.Flag = uint16(BinaryFlag | NotNullFlag)
	case string, []byte:
		field.Charset = 33
		field.Type = TypeVarString
	default:
		return fmt.Errorf("unsupport type %T for resultset", value)
	}
	return nil
}

// formatValue encode value into a string format
func formatValue(value interface{}) ([]byte, error) {
	if value == nil {
		return hack.Slice("NULL"), nil
	}
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return hack.Slice(v), nil
	default:
		return nil, fmt.Errorf("invalid type %T", value)
	}
}
