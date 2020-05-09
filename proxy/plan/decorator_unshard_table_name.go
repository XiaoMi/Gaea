package plan

import (
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
)

// TableNameDecorator decorate TableName
type UnshardTableNameDecorator struct {
	origin *ast.TableName
	phyDB  string
}

// CreateTableNameDecorator create TableNameDecorator
// the table has been checked before
func CreateUnshardTableNameDecorator(n *ast.TableName, phyDb string) *UnshardTableNameDecorator {
	ret := &UnshardTableNameDecorator{
		origin: n,
		phyDB:  phyDb,
	}
	return ret
}

// Restore implement ast.Node
func (t *UnshardTableNameDecorator) Restore(ctx *format.RestoreCtx) error {
	if t.origin.Schema.String() != "" {
		ctx.WriteName(t.phyDB)
		ctx.WritePlain(".")
	}
	ctx.WriteName(t.origin.Name.String())
	return nil
}

// Accept implement ast.Node
// do nothing and return current decorator
func (t *UnshardTableNameDecorator) Accept(v ast.Visitor) (ast.Node, bool) {
	return t, true
}

// Text implement ast.Node
func (t *UnshardTableNameDecorator) Text() string {
	return t.origin.Text()
}

// SetText implement ast.Node
func (t *UnshardTableNameDecorator) SetText(text string) {
	t.origin.SetText(text)
}
