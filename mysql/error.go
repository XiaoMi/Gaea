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

// Package mysql error.go contains error code, SQLSTATE value, and message string.
// The struct SqlError contains above informations and used in programs.

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
	"errors"
	"fmt"
)

var (
	// ErrBadConn bad connection error
	ErrBadConn = errors.New("connection was bad")
	// ErrMalformPacket packet error
	ErrMalformPacket = errors.New("Malform packet error")
	// ErrTxDone transaction done error
	ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")
)

// SQLError contains error code、SQLSTATE and message string
// https://dev.mysql.com/doc/refman/5.7/en/server-error-reference.html
type SQLError struct {
	Code    uint16
	Message string
	State   string
}

func (se *SQLError) Error() string {
	return fmt.Sprintf("ERROR %d (%s): %s", se.Code, se.State, se.Message)
}

// SQLCode returns the internal MySQL error code.
func (se *SQLError) SQLCode() uint16 {
	return se.Code
}

// SQLState returns the SQLSTATE value.
func (se *SQLError) SQLState() string {
	return se.State
}

// NewDefaultError default mysql error, must adapt errname message format
func NewDefaultError(errCode uint16, args ...interface{}) *SQLError {
	e := new(SQLError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DefaultMySQLState
	}

	if format, ok := MySQLErrName[errCode]; ok {
		e.Message = fmt.Sprintf(format, args...)
	} else {
		e.Message = fmt.Sprint(args...)
	}

	return e
}

// NewError create new error with specified code and message
func NewError(errCode uint16, message string) *SQLError {
	e := new(SQLError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DefaultMySQLState
	}

	e.Message = message

	return e
}

// NewErrf creates a SQL error, with an error code and a format specifier.
func NewErrf(errCode uint16, format string, args ...interface{}) *SQLError {
	e := &SQLError{Code: errCode}

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DefaultMySQLState
	}

	e.Message = fmt.Sprintf(format, args...)

	return e
}
