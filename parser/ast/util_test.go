// Copyright 2017 PingCAP, Inc.
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

package ast_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/parser"
	. "github.com/XiaoMi/Gaea/parser/ast"
	. "github.com/XiaoMi/Gaea/parser/format"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
)

func TestCacheable(t *testing.T) {
	// test non-SelectStmt
	var stmt Node = &DeleteStmt{}
	require.False(t, IsReadOnly(stmt))

	stmt = &InsertStmt{}
	require.False(t, IsReadOnly(stmt))

	stmt = &UpdateStmt{}
	require.False(t, IsReadOnly(stmt))

	stmt = &ExplainStmt{}
	require.True(t, IsReadOnly(stmt))

	stmt = &ExplainStmt{}
	require.True(t, IsReadOnly(stmt))

	stmt = &DoStmt{}
	require.True(t, IsReadOnly(stmt))
}

// CleanNodeText set the text of node and all child node empty.
// For test only.
func CleanNodeText(node Node) {
	var cleaner nodeTextCleaner
	node.Accept(&cleaner)
}

// nodeTextCleaner clean the text of a node and it's child node.
// For test only.
type nodeTextCleaner struct {
}

// Enter implements Visitor interface.
func (checker *nodeTextCleaner) Enter(in Node) (out Node, skipChildren bool) {
	in.SetText("")
	switch node := in.(type) {
	case *Constraint:
		if node.Option != nil {
			if node.Option.KeyBlockSize == 0x0 && node.Option.Tp == 0 && node.Option.Comment == "" {
				node.Option = nil
			}
		}
	case *FuncCallExpr:
		node.FnName.O = strings.ToLower(node.FnName.O)
		switch node.FnName.L {
		case "convert":
			node.Args[1].(*driver.ValueExpr).Datum.SetBytes(nil)
		}
	case *AggregateFuncExpr:
		node.F = strings.ToLower(node.F)
	case *FieldList:
		for _, f := range node.Fields {
			f.Offset = 0
		}
	case *AlterTableSpec:
		for _, opt := range node.Options {
			opt.StrValue = strings.ToLower(opt.StrValue)
		}
	}
	return in, false
}

// Leave implements Visitor interface.
func (checker *nodeTextCleaner) Leave(in Node) (out Node, ok bool) {
	return in, true
}

type NodeRestoreTestCase struct {
	sourceSQL string
	expectSQL string
}

func runNodeRestoreTest(t *testing.T, nodeTestCases []NodeRestoreTestCase, template string, extractNodeFunc func(node Node) Node) {
	runNodeRestoreTestWithFlags(t, nodeTestCases, template, extractNodeFunc, DefaultRestoreFlags)
}

func runNodeRestoreTestWithFlags(t *testing.T, nodeTestCases []NodeRestoreTestCase, template string, extractNodeFunc func(node Node) Node, flags RestoreFlags) {
	p := parser.New()
	p.EnableWindowFunc(true)
	for _, testCase := range nodeTestCases {
		sourceSQL := fmt.Sprintf(template, testCase.sourceSQL)
		expectSQL := fmt.Sprintf(template, testCase.expectSQL)
		stmt, err := p.ParseOneStmt(sourceSQL, "", "")
		comment := fmt.Sprintf("source %#v", testCase)
		require.NoError(t, err, comment)
		var sb strings.Builder
		err = extractNodeFunc(stmt).Restore(NewRestoreCtx(flags, &sb))
		require.NoError(t, err, comment)
		restoreSql := fmt.Sprintf(template, sb.String())
		comment = fmt.Sprintf("source %#v; restore %v", testCase, restoreSql)
		require.Equal(t, expectSQL, restoreSql, comment)
		stmt2, err := p.ParseOneStmt(restoreSql, "", "")
		require.NoError(t, err, comment)
		CleanNodeText(stmt)
		CleanNodeText(stmt2)
		require.Equal(t, stmt, stmt2, comment)
	}
}
