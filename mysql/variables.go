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
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/core/errors"
)

type verifyFunc func(interface{}) error

// allowed session variables
const (
	SQLModeStr     = "sql_mode"
	SQLSafeUpdates = "sql_safe_updates"
	TimeZone       = "time_zone"
)

// not allowed session variables
const (
	MaxAllowedPacket = "max_allowed_packet"
)

var variableVerifyFuncMap = map[string]verifyFunc{
	SQLModeStr:     verifySQLMode,
	SQLSafeUpdates: verifyOnOffInteger,
	TimeZone:       verifyTimeZone,
}

// SessionVariables variables in session
type SessionVariables struct {
	variables map[string]*Variable
	unused    map[string]*Variable
}

// NewSessionVariables constructor of SessionVariables
func NewSessionVariables() *SessionVariables {
	return &SessionVariables{
		variables: make(map[string]*Variable),
		unused:    make(map[string]*Variable),
	}
}

// Equals check if equal of SessionVariables
func (s *SessionVariables) Equals(dst *SessionVariables) bool {
	if len(s.variables) != len(dst.variables) {
		return false
	}

	for _, v := range s.variables {
		if dstV, ok := dst.variables[v.Name()]; !ok {
			return false
		} else if dstV != v {
			return false
		}
	}
	return true
}

// SetEqualsWith set the SessionVariables equals with the dst, and variables not contained in dst are moved to unused.
func (s *SessionVariables) SetEqualsWith(dst *SessionVariables) ( /*changed*/ bool, error) {
	if len(s.variables) == 0 && len(dst.variables) != 0 {
		for _, v := range dst.variables {
			if err := s.Set(v.Name(), v.Get()); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	if len(s.variables) != 0 && len(dst.variables) == 0 {
		for _, v := range s.variables {
			s.unused[v.Name()] = v
			delete(s.variables, v.Name())
		}
		return true, nil
	}

	changed := false
	for variableName := range variableVerifyFuncMap {
		srcVar, srcOK := s.variables[variableName]
		dstVar, dstOK := dst.variables[variableName]
		if srcOK && dstOK {
			if srcVar.Get() != dstVar.Get() {
				changed = true
				srcVar.Set(dstVar.Get())
			}
		} else if srcOK && !dstOK {
			changed = true
			s.unused[variableName] = srcVar
			delete(s.variables, variableName)
		} else if !srcOK && dstOK {
			changed = true
			s.Set(variableName, dstVar.Get())
		}
	}

	return changed, nil
}

// Delete delete variables with specific key
func (s *SessionVariables) Delete(key string) {
	delete(s.variables, formatVariableName(key))
}

// Set store variable in session
func (s *SessionVariables) Set(key string, value interface{}) error {
	formatKey := formatVariableName(key)
	verifyFunc, ok := variableVerifyFuncMap[formatKey]
	if !ok {
		return fmt.Errorf("variable not support")
	}

	if variable, ok := s.variables[formatKey]; ok {
		return variable.Set(value)
	}

	variable, err := NewVariable(formatKey, value, verifyFunc)
	if err != nil {
		return err
	}
	s.variables[formatKey] = variable
	return nil
}

// Get return variable with specific key
func (s *SessionVariables) Get(key string) (interface{}, bool) {
	v, ok := s.variables[key]
	return v, ok
}

// GetAll return all variables in session
func (s *SessionVariables) GetAll() map[string]*Variable {
	return s.variables
}

// GetUnusedAndClear unused variables
func (s *SessionVariables) GetUnusedAndClear() map[string]*Variable {
	unused := s.unused
	s.unused = make(map[string]*Variable)
	return unused
}

func formatVariableName(name string) string {
	name = strings.Trim(name, "'`\"")
	name = strings.ToLower(name)
	return name
}

// Variable variable definition in session
type Variable struct {
	name   string
	value  interface{}
	verify verifyFunc
}

// NewVariable constructor of Variable
func NewVariable(name string, value interface{}, verify verifyFunc) (*Variable, error) {
	v := &Variable{
		name:   formatVariableName(name),
		value:  value,
		verify: verify,
	}
	if err := v.verify(value); err != nil {
		return nil, err
	}
	return v, nil
}

// Set store data
func (v *Variable) Set(value interface{}) error {
	if err := v.verify(value); err != nil {
		return err
	}
	v.value = value
	return nil
}

// Name name of variable
func (v *Variable) Name() string {
	return v.name
}

// Get return value in Variable
func (v *Variable) Get() interface{} {
	return v.value
}

func verifySQLMode(v interface{}) error {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid type of sql mode")
	}
	if value == "" {
		return nil
	}

	value = strings.Trim(value, "'`\"")
	value = strings.ToUpper(value)
	values := strings.Split(value, ",")
	for _, sqlMode := range values {
		if _, ok := SQLModeSet[sqlMode]; !ok {
			return errors.ErrInvalidSQLMode
		}
	}
	return nil
}

// SQLModeSet https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html
var SQLModeSet = map[string]bool{
	// Full List of SQL Modes
	"ALLOW_INVALID_DATES":        true,
	"ANSI_QUOTES":                true,
	"ERROR_FOR_DIVISION_BY_ZERO": true,
	"HIGH_NOT_PRECEDENCE":        true,
	"IGNORE_SPACE":               true,
	"NO_AUTO_CREATE_USER":        true,
	"NO_AUTO_VALUE_ON_ZERO":      true,
	"NO_BACKSLASH_ESCAPES":       true,
	"NO_DIR_IN_CREATE":           true,
	"NO_ENGINE_SUBSTITUTION":     true,
	"NO_FIELD_OPTIONS":           true,
	"NO_KEY_OPTIONS":             true,
	"NO_TABLE_OPTIONS":           true,
	"NO_UNSIGNED_SUBTRACTION":    true,
	"NO_ZERO_DATE":               true,
	"NO_ZERO_IN_DATE":            true,
	"ONLY_FULL_GROUP_BY":         true,
	"PAD_CHAR_TO_FULL_LENGTH":    true,
	"PIPES_AS_CONCAT":            true,
	"REAL_AS_FLOAT":              true,
	"STRICT_ALL_TABLES":          true,
	"STRICT_TRANS_TABLES":        true,

	// Combination SQL Modes
	"ANSI":        true,
	"DB2":         true,
	"MAXDB":       true,
	"MSSQL":       true,
	"MYSQL323":    true,
	"MYSQL40":     true,
	"ORACLE":      true,
	"POSTGRESQL":  true,
	"TRADITIONAL": true,
}

func verifyOnOffInteger(v interface{}) error {
	val, ok := v.(int64)
	if !ok {
		return fmt.Errorf("value is not int64")
	}
	if val != 0 && val != 1 {
		return fmt.Errorf("value is not 0 or 1")
	}
	return nil
}

func verifyTimeZone(v interface{}) error {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid type of time_zone")
	}
	values := strings.Split(value, ":")
	if len(values) != 2 {
		return fmt.Errorf("invalid format of time_zone")
	}
	if values[0][0] != '+' && values[0][0] != '-' {
		return fmt.Errorf("invalid format of time_zone")
	}
	hour, err := strconv.Atoi(values[0])
	if err != nil {
		return fmt.Errorf("invalid hour of time_zone")
	}
	minute, err := strconv.Atoi(values[1])
	if err != nil {
		return fmt.Errorf("invalid minute of time_zone")
	}
	var directMinute int
	if hour < 0 {
		directMinute = hour*60 - minute
	} else {
		directMinute = hour*60 + minute
	}
	if directMinute < -779 || directMinute > 780 {
		return fmt.Errorf("exceed limit of time_zone")
	}

	return nil
}
