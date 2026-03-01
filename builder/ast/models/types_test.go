package models

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields_AST(t *testing.T) {
	fields := Fields{
		&ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent("fieldName"),
			},
		},
	}
	astFields := fields.Ast()

	for i := range astFields {
		assert.Equal(t, (*ast.Field)(fields[i]), astFields[i])
	}
}

func TestExprs_AST(t *testing.T) {
	exprs := Exprs{
		ast.NewIdent("expr"),
	}
	astExprs := exprs.Ast()

	for i := range astExprs {
		assert.Equal(t, exprs[i], astExprs[i])
	}
}

func TestStmts_AST(t *testing.T) {
	stmts := Stmts{
		GetExprStmt(ast.NewIdent("expr")),
	}
	astStmts := stmts.Ast()

	for i := range astStmts {
		assert.Equal(t, stmts[i], astStmts[i])
	}
}
