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

package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
)

var p = &mysql.Field{Name: []byte("?")}
var c = &mysql.Field{}

func calcParams(sql string) (paramCount int, offsets []int, err error) {
	count := 0
	quoteChar := ""
	paramCount = 0
	offsets = make([]int, 0)

	for i, elem := range []byte(sql) {
		if elem == '\\' {
			continue
		} else if elem == '"' || elem == '\'' {
			if quoteChar == "" {
				quoteChar = string(elem)
			} else if quoteChar == string(elem) {
				quoteChar = ""
			}
		} else if quoteChar == "" && elem == '?' {
			count++
			offsets = append(offsets, i)
		}

	}
	if quoteChar != "" {
		err = fmt.Errorf("fatal situation")
		return
	}

	paramCount = count

	return
}

func escapeSQL(sql string) string {
	t := make([]byte, 0, len(sql))
	for _, elem := range []byte(sql) {
		if elem == '\\' || elem == '\'' {
			t = append(t, '\\')
		}
		t = append(t, elem)
	}
	return string(t)
}

// Stmt prepare statement struct
type Stmt struct {
	id          uint32
	sql         string
	args        []interface{}
	columnCount int
	paramCount  int
	paramTypes  []byte
	offsets     []int
}

// ResetParams reset args
func (s *Stmt) ResetParams() {
	s.args = make([]interface{}, s.paramCount)
}

func (s *Stmt) SetParamTypes(paramTypes []byte) {
	s.paramTypes = paramTypes
}

func (s *Stmt) GetParamTypes() []byte {
	return s.paramTypes
}

// GetRewriteSQL get rewrite sql
func (s *Stmt) GetRewriteSQL() (string, error) {
	sql := s.sql
	var tmp = ""
	var pos = 0
	var offset = 0
	var quote = false
	for i := 0; i < s.paramCount; i++ {
		quote, tmp = util.ItoString(s.args[i])
		tmp = escapeSQL(tmp)
		pos = s.offsets[i]
		if quote {
			sql = util.Concat(sql[:pos+offset], "'", tmp, "'", sql[pos+offset+1:])
			offset = offset + len(tmp) - 1 + 2
		} else {
			sql = util.Concat(sql[:pos+offset], tmp, sql[pos+offset+1:])
			offset = offset + len(tmp) - 1
		}
	}
	return sql, nil
}

func (se *SessionExecutor) handleStmtExecute(data []byte) (*mysql.Result, error) {
	if len(data) < 9 {
		return nil, mysql.ErrMalformPacket
	}

	pos := 0
	id := binary.LittleEndian.Uint32(data[0:4])
	pos += 4

	s, ok := se.stmts[id]
	if !ok {
		return nil, mysql.NewDefaultError(mysql.ErrUnknownStmtHandler,
			strconv.FormatUint(uint64(id), 10), "stmt_execute")
	}

	flag := data[pos] & mysql.CursorTypeReadOnly
	pos++
	//now we only support CURSOR_TYPE_NO_CURSOR flag
	if flag != 0 {
		return nil, mysql.NewError(mysql.ErrUnknown, fmt.Sprintf("unsupported flag %d", flag))
	}

	//skip iteration-count, always 1
	pos += 4

	var nullBitmaps []byte
	var paramTypes []byte
	var paramValues []byte

	paramNum := s.paramCount

	var executeSQL string
	var err error
	if paramNum > 0 {
		nullBitmapLen := (s.paramCount + 7) >> 3
		if len(data) < (pos + nullBitmapLen + 1) {
			return nil, mysql.ErrMalformPacket
		}
		nullBitmaps = data[pos : pos+nullBitmapLen]
		pos += nullBitmapLen

		//new param bound flag
		if data[pos] == 1 {
			pos++
			if len(data) < (pos + (paramNum << 1)) {
				return nil, mysql.ErrMalformPacket
			}

			paramTypes = data[pos : pos+(paramNum<<1)]
			pos += (paramNum << 1)

			paramValues = data[pos:]
			s.SetParamTypes(paramTypes)
		} else {
			paramValues = data[pos+1:]
		}

		if err := se.bindStmtArgs(s, nullBitmaps, s.GetParamTypes(), paramValues); err != nil {
			return nil, err
		}

		executeSQL, err = s.GetRewriteSQL()
		if err != nil {
			return nil, err
		}
	} else {
		executeSQL = s.sql
	}

	defer s.ResetParams()

	// execute sql using ComQuery
	r, err := se.handleQuery(executeSQL)
	if err != nil {
		return nil, err
	}

	// build binary result set
	if r != nil && r.Resultset != nil {
		resultSet, err := mysql.BuildBinaryResultset(r.Fields, r.Values)
		if err != nil {
			return nil, err
		}
		r.Resultset = resultSet
	}

	return r, nil
}

// long data and generic args are all in s.args
func (se *SessionExecutor) bindStmtArgs(s *Stmt, nullBitmap, paramTypes, paramValues []byte) error {
	args := s.args

	pos := 0

	var v []byte
	var isNull bool

	for i := 0; i < s.paramCount; i++ {
		if nullBitmap[i>>3]&(1<<(uint(i)%8)) > 0 {
			args[i] = nil
			continue
		}

		if (i<<1)+1 >= len(paramTypes) {
			return mysql.ErrMalformPacket
		}

		tp := paramTypes[i<<1]
		isUnsigned := (paramTypes[(i<<1)+1] & 0x80) > 0

		if s.args[i] != nil {
			continue
		}
		switch tp {
		case mysql.TypeNull:
			args[i] = nil
			continue

		case mysql.TypeTiny:
			if len(paramValues) < (pos + 1) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint8(paramValues[pos])
			} else {
				args[i] = int8(paramValues[pos])
			}

			pos++
			continue

		case mysql.TypeShort, mysql.TypeYear:
			if len(paramValues) < (pos + 2) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint16(binary.LittleEndian.Uint16(paramValues[pos : pos+2]))
			} else {
				args[i] = int16((binary.LittleEndian.Uint16(paramValues[pos : pos+2])))
			}
			pos += 2
			continue

		case mysql.TypeInt24, mysql.TypeLong:
			if len(paramValues) < (pos + 4) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = uint32(binary.LittleEndian.Uint32(paramValues[pos : pos+4]))
			} else {
				args[i] = int32(binary.LittleEndian.Uint32(paramValues[pos : pos+4]))
			}
			pos += 4
			continue

		case mysql.TypeLonglong:
			if len(paramValues) < (pos + 8) {
				return mysql.ErrMalformPacket
			}

			if isUnsigned {
				args[i] = binary.LittleEndian.Uint64(paramValues[pos : pos+8])
			} else {
				args[i] = int64(binary.LittleEndian.Uint64(paramValues[pos : pos+8]))
			}
			pos += 8
			continue

		case mysql.TypeFloat:
			if len(paramValues) < (pos + 4) {
				return mysql.ErrMalformPacket
			}

			args[i] = float32(math.Float32frombits(binary.LittleEndian.Uint32(paramValues[pos : pos+4])))
			pos += 4
			continue

		case mysql.TypeDouble:
			if len(paramValues) < (pos + 8) {
				return mysql.ErrMalformPacket
			}

			args[i] = math.Float64frombits(binary.LittleEndian.Uint64(paramValues[pos : pos+8]))
			pos += 8
			continue

		case mysql.TypeDecimal, mysql.TypeNewDecimal, mysql.TypeVarchar,
			mysql.TypeBit, mysql.TypeEnum, mysql.TypeSet, mysql.TypeTinyBlob,
			mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob,
			mysql.TypeVarString, mysql.TypeString, mysql.TypeGeometry,
			mysql.TypeDate, mysql.TypeNewDate,
			mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDuration, mysql.TypeJSON:
			if len(paramValues) < (pos + 1) {
				return mysql.ErrMalformPacket
			}

			var ok = false
			v, pos, isNull, ok = mysql.ReadLenEncStringAsBytes(paramValues, pos)
			if !ok {
				return errors.New("ReadLenEncStringAsBytes in bindStmtArgs failed")
			}

			if !isNull {
				args[i] = v
				continue
			} else {
				args[i] = nil
				continue
			}
		default:
			return fmt.Errorf("Stmt Unknown FieldType %d", tp)
		}
	}
	return nil
}

func (se *SessionExecutor) handleStmtSendLongData(data []byte) error {
	if len(data) < 6 {
		return mysql.ErrMalformPacket
	}

	id := binary.LittleEndian.Uint32(data[0:4])

	s, ok := se.stmts[id]
	if !ok {
		return mysql.NewDefaultError(mysql.ErrUnknownStmtHandler,
			strconv.FormatUint(uint64(id), 10), "stmt_send_longdata")
	}

	paramID := binary.LittleEndian.Uint16(data[4:6])
	if paramID >= uint16(s.paramCount) {
		return mysql.NewDefaultError(mysql.ErrWrongArguments, "stmt_send_longdata")
	}

	if s.args[paramID] == nil {
		tmpSlice := make([]byte, len(data)-6)
		copy(tmpSlice, data[6:])
		s.args[paramID] = tmpSlice
	} else {
		if b, ok := s.args[paramID].([]byte); ok {
			b = append(b, data[6:]...)
			s.args[paramID] = b
		} else {
			return fmt.Errorf("invalid param long data type %T", s.args[paramID])
		}
	}

	return nil
}

func (se *SessionExecutor) handleStmtReset(data []byte) error {
	if len(data) < 4 {
		return mysql.ErrMalformPacket
	}

	id := binary.LittleEndian.Uint32(data[0:4])

	s, ok := se.stmts[id]
	if !ok {
		return mysql.NewDefaultError(mysql.ErrUnknownStmtHandler,
			strconv.FormatUint(uint64(id), 10), "stmt_reset")
	}

	s.ResetParams()
	return nil
}
