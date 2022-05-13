// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

package parser

// analyzer.go contains utility analysis functions.

import (
	"strings"
	"unicode"
)

// These constants are used to identify the SQL statement type.
const (
	StmtSelect = iota
	StmtStream
	StmtInsert
	StmtReplace
	StmtUpdate
	StmtDelete
	StmtDDL
	StmtBegin
	StmtCommit
	StmtRollback
	StmtSet
	StmtShow
	StmtUse
	StmtOther
	StmtUnknown
	StmtComment
	StmtSavepoint
)
const (
	eofChar = 0x100
)

// Preview analyzes the beginning of the query using a simpler and faster
// textual comparison to identify the statement type.
func Preview(sql string) int {
	trimmed := StripLeadingComments(sql)

	firstWord := trimmed
	if end := strings.IndexFunc(trimmed, unicode.IsSpace); end != -1 {
		firstWord = trimmed[:end]
	}
	firstWord = strings.TrimLeftFunc(firstWord, func(r rune) bool { return !unicode.IsLetter(r) })
	// Comparison is done in order of priority.
	loweredFirstWord := strings.ToLower(firstWord)
	switch loweredFirstWord {
	case "select":
		return StmtSelect
	case "stream":
		return StmtStream
	case "insert":
		return StmtInsert
	case "replace":
		return StmtReplace
	case "update":
		return StmtUpdate
	case "delete":
		return StmtDelete
	}
	// For the following statements it is not sufficient to rely
	// on loweredFirstWord. This is because they are not statements
	// in the grammar and we are relying on Preview to parse them.
	// For instance, we don't want: "BEGIN JUNK" to be parsed
	// as StmtBegin.
	trimmedNoComments, _ := SplitMarginComments(trimmed)
	switch strings.ToLower(trimmedNoComments) {
	case "begin", "start transaction":
		return StmtBegin
	}
	switch loweredFirstWord {
	case "create", "alter", "rename", "drop", "truncate", "flush":
		return StmtDDL
	case "set":
		return StmtSet
	case "show":
		return StmtShow
	case "use":
		return StmtUse
	case "analyze", "describe", "desc", "explain", "repair", "optimize":
		return StmtOther
	case "commit":
		return StmtCommit
	case "rollback":
		return StmtRollback
	case "savepoint":
		return StmtSavepoint
	}
	if strings.Index(trimmedNoComments, "release savepoint") == 0 {
		return StmtSavepoint
	}
	if strings.Index(trimmed, "/*!") == 0 {
		return StmtComment
	}
	return StmtUnknown
}

// StmtType returns the statement type as a string
func StmtType(stmtType int) string {
	switch stmtType {
	case StmtSelect:
		return "SELECT"
	case StmtStream:
		return "STREAM"
	case StmtInsert:
		return "INSERT"
	case StmtReplace:
		return "REPLACE"
	case StmtUpdate:
		return "UPDATE"
	case StmtDelete:
		return "DELETE"
	case StmtDDL:
		return "DDL"
	case StmtBegin:
		return "BEGIN"
	case StmtCommit:
		return "COMMIT"
	case StmtRollback:
		return "ROLLBACK"
	case StmtSavepoint:
		return "SAVEPOINT"
	case StmtSet:
		return "SET"
	case StmtShow:
		return "SHOW"
	case StmtUse:
		return "USE"
	case StmtOther:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

// SplitStatementToPieces split raw sql statement that may have multi sql pieces to sql pieces
// returns the sql pieces blob contains; or error if sql cannot be parsed
func SplitStatementToPieces(blob string) (pieces []string, err error) {
	// fast path: the vast majority of SQL statements do not have semicolons in them
	if blob == "" {
		return nil, nil
	}
	switch strings.IndexByte(blob, ';') {
	case -1: // if there is no semicolon, return blob as a whole
		return []string{blob}, nil
	case len(blob) - 1: // if there's a single semicolon and it's the last character, return blob without it
		return []string{blob[:len(blob)-1]}, nil
	}

	pieces = make([]string, 0, 16)
	tokenizer := NewScanner(blob)

	tkn := 0
	var pos Pos
	var stmt string
	stmtBegin := 0
	emptyStatement := true
loop:
	for {
		tkn, pos, _ = tokenizer.scan()
		switch tkn {
		case ';':
			stmt = blob[stmtBegin:pos.Offset]
			if !emptyStatement {
				pieces = append(pieces, stmt)
				emptyStatement = true
			}
			stmtBegin = pos.Offset + 1
		case 0, eofChar:
			blobTail := pos.Offset - 1
			if stmtBegin < blobTail {
				stmt = blob[stmtBegin : blobTail+1]
				if !emptyStatement {
					pieces = append(pieces, stmt)
				}
			}
			break loop
		default:
			emptyStatement = false
		}
	}

	if len(tokenizer.errs) > 0 {
		err = tokenizer.errs[0]
	}

	return
}
