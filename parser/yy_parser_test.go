package parser

import (
	"fmt"
	"testing"

	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
)

func TestNodeToString(t *testing.T) {
	tableName := "tb1"
	d := &driver.ValueExpr{}
	d.SetValue(tableName)
	s, err := NodeToStringWithoutQuote(d)
	if err != nil {
		t.Fatal(err)
	}
	if s != tableName {
		t.Errorf("table name not equal, expect: %s, actual: %s", tableName, s)
	}
}

type NodePrintVisitor struct {
}

func (v *NodePrintVisitor) Enter(n ast.Node) (ast.Node, bool) {
	fmt.Printf("enter: %T\n", n)
	return n, false
}

func (v *NodePrintVisitor) Leave(n ast.Node) (ast.Node, bool) {
	fmt.Printf("leave: %T\n", n)
	return n, true
}

func TestASTNode(t *testing.T) {
	sql := `desc xm_order`
	n, err := ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}
	v := &NodePrintVisitor{}
	n.Accept(v)
}
