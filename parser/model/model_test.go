// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/types"
)

func TestT(t *testing.T) {
	abc := NewCIStr("aBC")
	require.Equal(t, "aBC", abc.O)
	require.Equal(t, "abc", abc.L)
	require.Equal(t, "aBC", abc.String())

}

func TestModelBasic(t *testing.T) {
	column := &ColumnInfo{
		ID:           1,
		Name:         NewCIStr("c"),
		Offset:       0,
		DefaultValue: 0,
		FieldType:    *types.NewFieldType(0),
	}
	column.Flag |= mysql.PriKeyFlag

	index := &IndexInfo{
		Name:  NewCIStr("key"),
		Table: NewCIStr("t"),
		Columns: []*IndexColumn{
			{
				Name:   NewCIStr("c"),
				Offset: 0,
				Length: 10,
			}},
		Unique:  true,
		Primary: true,
	}

	fk := &FKInfo{
		RefCols: []CIStr{NewCIStr("a")},
		Cols:    []CIStr{NewCIStr("a")},
	}

	table := &TableInfo{
		ID:          1,
		Name:        NewCIStr("t"),
		Charset:     "utf8",
		Collate:     "utf8_bin",
		Columns:     []*ColumnInfo{column},
		Indices:     []*IndexInfo{index},
		ForeignKeys: []*FKInfo{fk},
		PKIsHandle:  true,
	}

	dbInfo := &DBInfo{
		ID:      1,
		Name:    NewCIStr("test"),
		Charset: "utf8",
		Collate: "utf8_bin",
		Tables:  []*TableInfo{table},
	}

	n := dbInfo.Clone()
	require.Equal(t, dbInfo, n)

	pkName := table.GetPkName()
	require.Equal(t, NewCIStr("c"), pkName)
	newColumn := table.GetPkColInfo()
	require.Equal(t, column, newColumn)
	inIdx := table.ColumnIsInIndex(column)
	require.True(t, inIdx)
	tp := IndexTypeBtree
	require.Equal(t, "BTREE", tp.String())
	tp = IndexTypeHash
	require.Equal(t, "HASH", tp.String())
	tp = 1e5
	require.Equal(t, "", tp.String())
	has := index.HasPrefixIndex()
	require.True(t, has)
	tbl := table.GetUpdateTime()
	require.Equal(t, TSConvert2Time(table.UpdateTS), tbl)

	// Corner cases
	column.Flag ^= mysql.PriKeyFlag
	pkName = table.GetPkName()
	require.Equal(t, NewCIStr(""), pkName)
	newColumn = table.GetPkColInfo()
	require.Nil(t, newColumn)
	anCol := &ColumnInfo{
		Name: NewCIStr("d"),
	}
	exIdx := table.ColumnIsInIndex(anCol)
	require.False(t, exIdx)
	anIndex := &IndexInfo{
		Columns: []*IndexColumn{},
	}
	no := anIndex.HasPrefixIndex()
	require.False(t, no)
}

func TestJobStartTime(t *testing.T) {
	job := &Job{
		ID:         123,
		BinlogInfo: &HistoryInfo{},
	}
	time := time.Unix(0, 0)
	require.Equal(t, TSConvert2Time(job.StartTS), time)
	ret := fmt.Sprintf("%s", job)
	require.Equal(t, ret, job.String())
}

func TestJobCodec(t *testing.T) {
	type A struct {
		Name string
	}
	job := &Job{
		ID:         1,
		TableID:    2,
		SchemaID:   1,
		BinlogInfo: &HistoryInfo{},
		Args:       []interface{}{NewCIStr("a"), A{Name: "abc"}},
	}
	job.BinlogInfo.AddDBInfo(123, &DBInfo{ID: 1, Name: NewCIStr("test_history_db")})
	job.BinlogInfo.AddTableInfo(123, &TableInfo{ID: 1, Name: NewCIStr("test_history_tbl")})

	// Test IsDependentOn.
	// job: table ID is 2
	// job1: table ID is 2
	var err error
	job1 := &Job{
		ID:         2,
		TableID:    2,
		SchemaID:   1,
		Type:       ActionRenameTable,
		BinlogInfo: &HistoryInfo{},
		Args:       []interface{}{int64(3), NewCIStr("new_table_name")},
	}
	job1.RawArgs, err = json.Marshal(job1.Args)
	require.NoError(t, err)
	isDependent, err := job.IsDependentOn(job1)
	require.NoError(t, err)
	require.True(t, isDependent)
	// job1: rename table, old schema ID is 3
	// job2: create schema, schema ID is 3
	job2 := &Job{
		ID:         3,
		TableID:    3,
		SchemaID:   3,
		Type:       ActionCreateSchema,
		BinlogInfo: &HistoryInfo{},
	}
	isDependent, err = job2.IsDependentOn(job1)
	require.NoError(t, err)
	require.True(t, isDependent)
	require.False(t, job.IsCancelled())
	b, err := job.Encode(false)
	require.NoError(t, err)
	newJob := &Job{}
	err = newJob.Decode(b)
	require.NoError(t, err)
	require.Equal(t, job.BinlogInfo, newJob.BinlogInfo)
	name := CIStr{}
	a := A{}
	err = newJob.DecodeArgs(&name, &a)
	require.NoError(t, err)
	require.Equal(t, NewCIStr(""), name)
	require.Equal(t, A{Name: ""}, a)
	require.Greater(t, len(newJob.String()), 0)

	job.BinlogInfo.Clean()
	b1, err := job.Encode(true)
	require.NoError(t, err)
	newJob = &Job{}
	err = newJob.Decode(b1)
	require.NoError(t, err)
	require.Equal(t, &HistoryInfo{}, newJob.BinlogInfo)
	name = CIStr{}
	a = A{}
	err = newJob.DecodeArgs(&name, &a)
	require.NoError(t, err)
	require.Equal(t, NewCIStr("a"), name)
	require.Equal(t, A{Name: "abc"}, a)
	require.Greater(t, len(newJob.String()), 0)

	b2, err := job.Encode(true)
	require.NoError(t, err)
	newJob = &Job{}
	err = newJob.Decode(b2)
	require.NoError(t, err)
	name = CIStr{}
	// Don't decode to a here.
	err = newJob.DecodeArgs(&name)
	require.NoError(t, err)
	require.Equal(t, NewCIStr("a"), name)
	require.Greater(t, len(newJob.String()), 0)

	job.State = JobStateDone
	require.True(t, job.IsDone())
	require.True(t, job.IsFinished())
	require.False(t, job.IsRunning())
	require.False(t, job.IsSynced())
	require.False(t, job.IsRollbackDone())

	job.SetRowCount(3)
	require.Equal(t, int64(3), job.GetRowCount())

}

func TestState(t *testing.T) {
	schemaTbl := []SchemaState{
		StateDeleteOnly,
		StateWriteOnly,
		StateWriteReorganization,
		StateDeleteReorganization,
		StatePublic,
	}

	for _, state := range schemaTbl {
		require.Greater(t, len(state.String()), 0)
	}

	jobTbl := []JobState{
		JobStateRunning,
		JobStateDone,
		JobStateCancelled,
		JobStateRollingback,
		JobStateRollbackDone,
		JobStateSynced,
	}

	for _, state := range jobTbl {
		require.Greater(t, len(state.String()), 0)
	}
}

func TestString(t *testing.T) {
	acts := []struct {
		act    ActionType
		result string
	}{
		{ActionNone, "none"},
		{ActionAddForeignKey, "add foreign key"},
		{ActionDropForeignKey, "drop foreign key"},
		{ActionTruncateTable, "truncate table"},
		{ActionModifyColumn, "modify column"},
		{ActionRenameTable, "rename table"},
		{ActionSetDefaultValue, "set default value"},
		{ActionCreateSchema, "create schema"},
		{ActionDropSchema, "drop schema"},
		{ActionCreateTable, "create table"},
		{ActionDropTable, "drop table"},
		{ActionAddIndex, "add index"},
		{ActionDropIndex, "drop index"},
		{ActionAddColumn, "add column"},
		{ActionDropColumn, "drop column"},
	}

	for _, v := range acts {
		str := v.act.String()
		require.Equal(t, v.result, str)
	}
}

func TestUnmarshalCIStr(t *testing.T) {
	var ci CIStr

	// Test unmarshal CIStr from a single string.
	str := "aaBB"
	buf, err := json.Marshal(str)
	require.NoError(t, err)
	require.NoError(t, ci.UnmarshalJSON(buf))
	require.Equal(t, str, ci.O)
	require.Equal(t, "aabb", ci.L)

	buf, err = json.Marshal(ci)
	require.NoError(t, err)
	require.Equal(t, `{"O":"aaBB","L":"aabb"}`, string(buf))
	require.NoError(t, ci.UnmarshalJSON(buf))
	require.Equal(t, str, ci.O)
	require.Equal(t, "aabb", ci.L)
}
