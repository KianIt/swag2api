package models

import (
	"go/ast"
)

type (
	Field    *ast.Field
	Fields   []Field
	BasicLit *ast.BasicLit
	Expr     ast.Expr
	Exprs    []Expr
	Stmt     ast.Stmt
	Stmts    []Stmt
	Decl     ast.Decl
)

// Ast returns Fields as []*ast.Field.
func (fs Fields) Ast() []*ast.Field {
	astFs := make([]*ast.Field, 0, len(fs))

	for _, f := range fs {
		astFs = append(astFs, f)
	}

	return astFs
}

// Ast returns Exprs as []ast.Expr.
func (es Exprs) Ast() []ast.Expr {
	astEs := make([]ast.Expr, 0, len(es))

	for _, e := range es {
		astEs = append(astEs, e)
	}

	return astEs
}

// Ast returns Stmts as []ast.Stmt.
func (ss Stmts) Ast() []ast.Stmt {
	astSs := make([]ast.Stmt, 0, len(ss))

	for _, e := range ss {
		astSs = append(astSs, e)
	}

	return astSs
}
