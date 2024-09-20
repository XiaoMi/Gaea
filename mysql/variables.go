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
)

type verifyFunc func(interface{}) error

// allowed session variables
const (
	SQLModeStr             = "sql_mode"
	SQLSafeUpdates         = "sql_safe_updates"
	TimeZone               = "time_zone"
	SQLSelectLimit         = "sql_select_limit"
	TxReadOnly             = "tx_read_only"
	TransactionReadOnly    = "transaction_read_only"
	CharacterSetConnection = "character_set_connection"
	CharacterSetResults    = "character_set_results"
	CharacterSetClient     = "character_set_client"
	GroupConcatMaxLen      = "group_concat_max_len"
	MaxExecutionTime       = "max_execution_time"
	UniqueChecks           = "unique_checks"
	TransactionIsolation   = "transaction_isolation"
)

// not allowed session variables
const (
	MaxAllowedPacket = "max_allowed_packet"
)

var variableVerifyFuncMap = map[string]verifyFunc{
	SQLModeStr:             verifySQLMode,
	SQLSafeUpdates:         verifyOnOffInteger,
	TimeZone:               verifyTimeZone,
	SQLSelectLimit:         verifyInteger,
	TxReadOnly:             verifyOnOffInteger,
	TransactionReadOnly:    verifyOnOffInteger,
	CharacterSetConnection: verifyString,
	CharacterSetResults:    verifyString,
	CharacterSetClient:     verifyString,
	GroupConcatMaxLen:      verifyInteger,
	MaxExecutionTime:       verifyInteger,
	UniqueChecks:           verifyOnOffInteger,
	TransactionIsolation:   verifyString,
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

// SetEqualsWith sets the SessionVariables of the current instance (s) to be equal to those of the destination instance (dst).
// This method ensures that all session variables from dst are replicated in s, and any variables not present in dst are moved to the 'unused' map.
// It returns a boolean indicating if any variables were changed, and an error if any operations fail.
//
// The method handles several scenarios:
//
//  1. If s.variables is empty and dst.variables is not, all variables from dst are copied to s.
//     This scenario is for initializing s with the settings from dst when s has no prior variables.
//
//  2. If s.variables is not empty and dst.variables is empty, all variables from s are moved to s.unused.
//     This scenario applies when it's determined that no session variables should be actively used in s,
//     perhaps resetting s to a state without active session variables.
//
// 3. The general case when both s.variables and dst.variables have variables:
//   - It iterates through all variables in dst.variables:
//     a. If a variable also exists in s.variables and their values differ, the value from dst is set in s.
//     b. If a variable does not exist in s.variables, it is added to s.
//   - It then checks for variables that exist in s.variables but not in dst.variables.
//     a. Such variables are moved to s.unused to indicate they are no longer actively needed.
//
// SetEqualsWith ensures that s.variables reflect exactly what is in dst.variables post-execution, with any extraneous variables moved to unused.
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

	// 用于记录是否有变量更新
	changed := false

	for name, dstVar := range dst.variables {
		if srcVar, ok := s.variables[name]; ok {
			// 如果源存在相同名称的变量并且值不同，则更新
			if srcVar.Get() != dstVar.Get() {
				if err := srcVar.Set(dstVar.Get()); err != nil {
					return false, err
				}
				changed = true
			}
		} else {
			// 如果源不存在这个变量，则添加
			if err := s.Set(name, dstVar.Get()); err != nil {
				return false, err
			}
			changed = true
		}
	}

	// 检查源中有而目标中没有的变量，这些变量应被视为不再使用
	for name, srcVar := range s.variables {
		if _, ok := dst.variables[name]; !ok {
			s.unused[name] = srcVar
			delete(s.variables, name)
			changed = true
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
		// Allows the program to perform default validation processing on unsupported or undefined variable names
		// instead of returning an error directly
		verifyFunc = verifyDefault
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

// Reset removes any session variables that are not recognized according to the current verification rules.
func (s *SessionVariables) Reset(err error) {
	// Retrieve all current session variables.
	allVars := s.GetAll()
	// Iterate through all the variables.
	for key := range allVars {
		// Check if there is a verification function for the key in the map.
		if _, ok := variableVerifyFuncMap[key]; !ok {
			// If the key is not found in the verification function map, delete it from session variables.
			s.Delete(key)
		}
	}
	// Check if the error is related to an invalid sql_mode
	if IsWrongValueForSQLModeErr(err) {
		// Remove the invalid sql_mode from session variables
		s.RemoveInvalidSQLMode()
	}

}

func (s *SessionVariables) RemoveInvalidSQLMode() {
	// Check if 'sql_mode' exists in the session variables
	if _, ok := s.Get("sql_mode"); ok {
		// Assume the verification function is available in the variableVerifyFuncMap
		if _, exists := variableVerifyFuncMap["sql_mode"]; exists {
			// If verification fails, remove 'sql_mode'
			s.Delete("sql_mode")
		}
	}
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

func verifyInteger(v interface{}) error {
	_, ok := v.(int64)
	if !ok {
		return fmt.Errorf("value is not int64")
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

func verifyString(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("value is not string type")
	}
	return nil
}

func verifyDefault(v interface{}) error {
	return nil
}
