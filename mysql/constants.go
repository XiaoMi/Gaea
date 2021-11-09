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
	"strings"

	"github.com/pingcap/errors"

	"github.com/XiaoMi/Gaea/parser/format"
)

const (
	// MinProtocolVersion min protocol version
	MinProtocolVersion byte = 10
	// MaxPayloadLen max payload length
	MaxPayloadLen int = 1<<24 - 1
	// TimeFormat time format
	TimeFormat string = "2006-01-02 15:04:05"
	// ServerVersion server version
	ServerVersion string = "5.6.20-gaea"
	// MysqlNativePassword uses a salt and transmits a hash on the wire.
	MysqlNativePassword = "mysql_native_password"
	CachingSHA2Password = "caching_sha2_password"
	// ProtocolVersion is the current version of the protocol.
	// Always 10.
	ProtocolVersion = 10
)

const (
	// CursorTypeReadOnly readonly cursor
	CursorTypeReadOnly = 0x01
)

// Header information.
const (
	OKHeader          byte = 0x00
	ErrHeader         byte = 0xff
	EOFHeader         byte = 0xfe
	LocalInFileHeader byte = 0xfb
	AuthSwitchHeader  byte = 0xfe
)

// Server information.
const (
	ServerStatusInTrans            uint16 = 0x0001
	ServerStatusAutocommit         uint16 = 0x0002
	ServerMoreResultsExists        uint16 = 0x0008
	ServerStatusNoGoodIndexUsed    uint16 = 0x0010
	ServerStatusNoIndexUsed        uint16 = 0x0020
	ServerStatusCursorExists       uint16 = 0x0040
	ServerStatusLastRowSend        uint16 = 0x0080
	ServerStatusDBDropped          uint16 = 0x0100
	ServerStatusNoBackslashEscaped uint16 = 0x0200
	ServerStatusMetadataChanged    uint16 = 0x0400
	ServerStatusWasSlow            uint16 = 0x0800
	ServerPSOutParams              uint16 = 0x1000
)

// ErrTextLength error text length limit.
const ErrTextLength = 80

// Command information.
const (
	ComSleep byte = iota
	ComQuit
	ComInitDB
	ComQuery
	ComFieldList
	ComCreateDB
	ComDropDB
	ComRefresh
	ComShutdown
	ComStatistics
	ComProcessInfo
	ComConnect
	ComProcessKill
	ComDebug
	ComPing
	ComTime
	ComDelayedInsert
	ComChangeUser
	ComBinlogDump
	ComTableDump
	ComConnectOut
	ComRegisterSlave
	ComStmtPrepare
	ComStmtExecute
	ComStmtSendLongData
	ComStmtClose
	ComStmtReset
	ComSetOption
	ComStmtFetch
	ComDaemon
	ComBinlogDumpGtid
	ComResetConnection
	ComEnd
)

// Client information.
const (
	ClientLongPassword uint32 = 1 << iota
	ClientFoundRows
	ClientLongFlag
	ClientConnectWithDB
	ClientNoSchema
	ClientCompress
	ClientODBC
	ClientLocalFiles
	ClientIgnoreSpace
	ClientProtocol41
	ClientInteractive
	ClientSSL
	ClientIgnoreSigpipe
	ClientTransactions
	ClientReserved
	ClientSecureConnection
	ClientMultiStatements
	ClientMultiResults
	ClientPSMultiResults
	ClientPluginAuth
	ClientConnectAtts
	ClientPluginAuthLenencClientData
)

// PrivilegeType  privilege
type PrivilegeType uint32

const (
	_ PrivilegeType = 1 << iota
	// CreatePriv is the privilege to create schema/table.
	CreatePriv
	// SelectPriv is the privilege to read from table.
	SelectPriv
	// InsertPriv is the privilege to insert data into table.
	InsertPriv
	// UpdatePriv is the privilege to update data in table.
	UpdatePriv
	// DeletePriv is the privilege to delete data from table.
	DeletePriv
	// ShowDBPriv is the privilege to run show databases statement.
	ShowDBPriv
	// SuperPriv enables many operations and server behaviors.
	SuperPriv
	// CreateUserPriv is the privilege to create user.
	CreateUserPriv
	// TriggerPriv is not checked yet.
	TriggerPriv
	// DropPriv is the privilege to drop schema/table.
	DropPriv
	// ProcessPriv pertains to display of information about the threads executing within the server.
	ProcessPriv
	// GrantPriv is the privilege to grant privilege to user.
	GrantPriv
	// ReferencesPriv is not checked yet.
	ReferencesPriv
	// AlterPriv is the privilege to run alter statement.
	AlterPriv
	// ExecutePriv is the privilege to run execute statement.
	ExecutePriv
	// IndexPriv is the privilege to create/drop index.
	IndexPriv
	// CreateViewPriv is the privilege to create view.
	CreateViewPriv
	// ShowViewPriv is the privilege to show create view.
	ShowViewPriv
	// CreateRolePriv the privilege to create a role.
	CreateRolePriv
	// DropRolePriv is the privilege to drop a role.
	DropRolePriv
	// AllPriv is the privilege for all actions.
	AllPriv
)

// AllPrivMask is the mask for PrivilegeType with all bits set to 1.
// If it's passed to RequestVerification, it means any privilege would be OK.
const AllPrivMask = AllPriv - 1

// MySQL type maximum length.
const (
	// For arguments that have no fixed number of decimals, the decimals value is set to 31,
	// which is 1 more than the maximum number of decimals permitted for the DECIMAL, FLOAT, and DOUBLE data types.
	NotFixedDec = 31

	MaxIntWidth             = 20
	MaxRealWidth            = 23
	MaxFloatingTypeScale    = 30
	MaxFloatingTypeWidth    = 255
	MaxDecimalScale         = 30
	MaxDecimalWidth         = 65
	MaxDateWidth            = 10 // YYYY-MM-DD.
	MaxDatetimeWidthNoFsp   = 19 // YYYY-MM-DD HH:MM:SS
	MaxDatetimeWidthWithFsp = 26 // YYYY-MM-DD HH:MM:SS[.fraction]
	MaxDatetimeFullWidth    = 29 // YYYY-MM-DD HH:MM:SS.###### AM
	MaxDurationWidthNoFsp   = 10 // HH:MM:SS
	MaxDurationWidthWithFsp = 15 // HH:MM:SS[.fraction]
	MaxBlobWidth            = 16777216
)

// SQLMode is the type for MySQL sql_mode.
// See https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html
type SQLMode int

// HasNoZeroDateMode detects if 'NO_ZERO_DATE' mode is set in SQLMode
func (m SQLMode) HasNoZeroDateMode() bool {
	return m&ModeNoZeroDate == ModeNoZeroDate
}

// HasNoZeroInDateMode detects if 'NO_ZERO_IN_DATE' mode is set in SQLMode
func (m SQLMode) HasNoZeroInDateMode() bool {
	return m&ModeNoZeroInDate == ModeNoZeroInDate
}

// HasErrorForDivisionByZeroMode detects if 'ERROR_FOR_DIVISION_BY_ZERO' mode is set in SQLMode
func (m SQLMode) HasErrorForDivisionByZeroMode() bool {
	return m&ModeErrorForDivisionByZero == ModeErrorForDivisionByZero
}

// HasOnlyFullGroupBy detects if 'ONLY_FULL_GROUP_BY' mode is set in SQLMode
func (m SQLMode) HasOnlyFullGroupBy() bool {
	return m&ModeOnlyFullGroupBy == ModeOnlyFullGroupBy
}

// HasStrictMode detects if 'STRICT_TRANS_TABLES' or 'STRICT_ALL_TABLES' mode is set in SQLMode
func (m SQLMode) HasStrictMode() bool {
	return m&ModeStrictTransTables == ModeStrictTransTables || m&ModeStrictAllTables == ModeStrictAllTables
}

// HasPipesAsConcatMode detects if 'PIPES_AS_CONCAT' mode is set in SQLMode
func (m SQLMode) HasPipesAsConcatMode() bool {
	return m&ModePipesAsConcat == ModePipesAsConcat
}

// HasNoUnsignedSubtractionMode detects if 'NO_UNSIGNED_SUBTRACTION' mode is set in SQLMode
func (m SQLMode) HasNoUnsignedSubtractionMode() bool {
	return m&ModeNoUnsignedSubtraction == ModeNoUnsignedSubtraction
}

// HasHighNotPrecedenceMode detects if 'HIGH_NOT_PRECEDENCE' mode is set in SQLMode
func (m SQLMode) HasHighNotPrecedenceMode() bool {
	return m&ModeHighNotPrecedence == ModeHighNotPrecedence
}

// HasANSIQuotesMode detects if 'ANSI_QUOTES' mode is set in SQLMode
func (m SQLMode) HasANSIQuotesMode() bool {
	return m&ModeANSIQuotes == ModeANSIQuotes
}

// HasRealAsFloatMode detects if 'REAL_AS_FLOAT' mode is set in SQLMode
func (m SQLMode) HasRealAsFloatMode() bool {
	return m&ModeRealAsFloat == ModeRealAsFloat
}

// HasPadCharToFullLengthMode detects if 'PAD_CHAR_TO_FULL_LENGTH' mode is set in SQLMode
func (m SQLMode) HasPadCharToFullLengthMode() bool {
	return m&ModePadCharToFullLength == ModePadCharToFullLength
}

// HasNoBackslashEscapesMode detects if 'NO_BACKSLASH_ESCAPES' mode is set in SQLMode
func (m SQLMode) HasNoBackslashEscapesMode() bool {
	return m&ModeNoBackslashEscapes == ModeNoBackslashEscapes
}

// HasIgnoreSpaceMode detects if 'IGNORE_SPACE' mode is set in SQLMode
func (m SQLMode) HasIgnoreSpaceMode() bool {
	return m&ModeIgnoreSpace == ModeIgnoreSpace
}

// HasNoAutoCreateUserMode detects if 'NO_AUTO_CREATE_USER' mode is set in SQLMode
func (m SQLMode) HasNoAutoCreateUserMode() bool {
	return m&ModeNoAutoCreateUser == ModeNoAutoCreateUser
}

// HasAllowInvalidDatesMode detects if 'ALLOW_INVALID_DATES' mode is set in SQLMode
func (m SQLMode) HasAllowInvalidDatesMode() bool {
	return m&ModeAllowInvalidDates == ModeAllowInvalidDates
}

// consts for sql modes.
const (
	ModeNone        SQLMode = 0
	ModeRealAsFloat SQLMode = 1 << iota
	ModePipesAsConcat
	ModeANSIQuotes
	ModeIgnoreSpace
	ModeNotUsed
	ModeOnlyFullGroupBy
	ModeNoUnsignedSubtraction
	ModeNoDirInCreate
	ModePostgreSQL
	ModeOracle
	ModeMsSQL
	ModeDb2
	ModeMaxdb
	ModeNoKeyOptions
	ModeNoTableOptions
	ModeNoFieldOptions
	ModeMySQL323
	ModeMySQL40
	ModeANSI
	ModeNoAutoValueOnZero
	ModeNoBackslashEscapes
	ModeStrictTransTables
	ModeStrictAllTables
	ModeNoZeroInDate
	ModeNoZeroDate
	ModeInvalidDates
	ModeErrorForDivisionByZero
	ModeTraditional
	ModeNoAutoCreateUser
	ModeHighNotPrecedence
	ModeNoEngineSubstitution
	ModePadCharToFullLength
	ModeAllowInvalidDates
)

// GetSQLMode gets the sql mode for string literal. SQL_mode is a list of different modes separated by commas.
// The input string must be formatted by 'FormatSQLModeStr'
func GetSQLMode(s string) (SQLMode, error) {
	strs := strings.Split(s, ",")
	var sqlMode SQLMode
	for i, length := 0, len(strs); i < length; i++ {
		mode, ok := Str2SQLMode[strs[i]]
		if !ok && strs[i] != "" {
			return sqlMode, newInvalidModeErr(strs[i])
		}
		sqlMode = sqlMode | mode
	}
	return sqlMode, nil
}

// Str2SQLMode is the string represent of sql_mode to sql_mode map.
var Str2SQLMode = map[string]SQLMode{
	"REAL_AS_FLOAT":              ModeRealAsFloat,
	"PIPES_AS_CONCAT":            ModePipesAsConcat,
	"ANSI_QUOTES":                ModeANSIQuotes,
	"IGNORE_SPACE":               ModeIgnoreSpace,
	"NOT_USED":                   ModeNotUsed,
	"ONLY_FULL_GROUP_BY":         ModeOnlyFullGroupBy,
	"NO_UNSIGNED_SUBTRACTION":    ModeNoUnsignedSubtraction,
	"NO_DIR_IN_CREATE":           ModeNoDirInCreate,
	"POSTGRESQL":                 ModePostgreSQL,
	"ORACLE":                     ModeOracle,
	"MSSQL":                      ModeMsSQL,
	"DB2":                        ModeDb2,
	"MAXDB":                      ModeMaxdb,
	"NO_KEY_OPTIONS":             ModeNoKeyOptions,
	"NO_TABLE_OPTIONS":           ModeNoTableOptions,
	"NO_FIELD_OPTIONS":           ModeNoFieldOptions,
	"MYSQL323":                   ModeMySQL323,
	"MYSQL40":                    ModeMySQL40,
	"ANSI":                       ModeANSI,
	"NO_AUTO_VALUE_ON_ZERO":      ModeNoAutoValueOnZero,
	"NO_BACKSLASH_ESCAPES":       ModeNoBackslashEscapes,
	"STRICT_TRANS_TABLES":        ModeStrictTransTables,
	"STRICT_ALL_TABLES":          ModeStrictAllTables,
	"NO_ZERO_IN_DATE":            ModeNoZeroInDate,
	"NO_ZERO_DATE":               ModeNoZeroDate,
	"INVALID_DATES":              ModeInvalidDates,
	"ERROR_FOR_DIVISION_BY_ZERO": ModeErrorForDivisionByZero,
	"TRADITIONAL":                ModeTraditional,
	"NO_AUTO_CREATE_USER":        ModeNoAutoCreateUser,
	"HIGH_NOT_PRECEDENCE":        ModeHighNotPrecedence,
	"NO_ENGINE_SUBSTITUTION":     ModeNoEngineSubstitution,
	"PAD_CHAR_TO_FULL_LENGTH":    ModePadCharToFullLength,
	"ALLOW_INVALID_DATES":        ModeAllowInvalidDates,
}

func newInvalidModeErr(s string) error {
	return NewDefaultError(ErrWrongValueForVar, "sql_mode", s)
}

// PriorityEnum is defined for Priority const values.
type PriorityEnum int

// Priority const values.
// See https://dev.mysql.com/doc/refman/5.7/en/insert.html
const (
	NoPriority PriorityEnum = iota
	LowPriority
	HighPriority
	DelayedPriority
)

// Priority2Str is used to convert the statement priority to string.
var Priority2Str = map[PriorityEnum]string{
	NoPriority:      "NO_PRIORITY",
	LowPriority:     "LOW_PRIORITY",
	HighPriority:    "HIGH_PRIORITY",
	DelayedPriority: "DELAYED",
}

// Str2Priority is used to convert a string to a priority.
func Str2Priority(val string) PriorityEnum {
	val = strings.ToUpper(val)
	switch val {
	case "NO_PRIORITY":
		return NoPriority
	case "HIGH_PRIORITY":
		return HighPriority
	case "LOW_PRIORITY":
		return LowPriority
	case "DELAYED":
		return DelayedPriority
	default:
		return NoPriority
	}
}

// Restore implements Node interface.
func (n *PriorityEnum) Restore(ctx *format.RestoreCtx) error {
	switch *n {
	case NoPriority:
		return nil
	case LowPriority:
		ctx.WriteKeyWord("LOW_PRIORITY")
	case HighPriority:
		ctx.WriteKeyWord("HIGH_PRIORITY")
	case DelayedPriority:
		ctx.WriteKeyWord("DELAYED")
	default:
		return errors.Errorf("undefined PriorityEnum Type[%d]", *n)
	}
	return nil
}
