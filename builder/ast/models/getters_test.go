package models

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetField(t *testing.T) {
	field := GetField(
		"name",
		"typeExpr",
		"tag",
	)

	assert.Equal(t, "name", field.Names[0].Name)
	assert.Equal(t, "typeExpr", field.Type.(*ast.Ident).Name)
	assert.Equal(t, "tag", field.Tag.Value)
}

func TestGetStringBasicLit(t *testing.T) {
	lit := GetStringBasicLit("value")

	assert.Equal(t, "value", lit.Value)
}

func TestGetNeqExpr(t *testing.T) {
	expr := GetNeqExpr(
		ast.NewIdent("left"),
		ast.NewIdent("right"),
	)
	bin := expr.(*ast.BinaryExpr)

	assert.Equal(t, token.NEQ, bin.Op)
	assert.Equal(t, "left", bin.X.(*ast.Ident).Name)
	assert.Equal(t, "right", bin.Y.(*ast.Ident).Name)
}

func TestGetCallExpr(t *testing.T) {
	expr := GetCallExpr(
		ast.NewIdent("caller"),
		Exprs{ast.NewIdent("params")},
	)
	call := expr.(*ast.CallExpr)

	assert.Equal(t, "caller", call.Fun.(*ast.Ident).Name)
	assert.Equal(t, "params", call.Args[0].(*ast.Ident).Name)
}

func TestGetSelectorExpr(t *testing.T) {
	expr := GetSelectorExpr(
		ast.NewIdent("source"),
		"item",
	)
	selector := expr.(*ast.SelectorExpr)

	assert.Equal(t, "source", selector.X.(*ast.Ident).Name)
	assert.Equal(t, "item", selector.Sel.Name)
}

func TestGetStructLitExpr(t *testing.T) {
	expr := GetStructLitExpr(
		"name",
		Exprs{ast.NewIdent("element")},
	)
	lit := expr.(*ast.CompositeLit)

	assert.Equal(t, "name", lit.Type.(*ast.Ident).Name)
	assert.Equal(t, "element", lit.Elts[0].(*ast.Ident).Name)
}

func TestGetKeyValueExpr(t *testing.T) {
	expr := GetKeyValueExpr(
		ast.NewIdent("key"),
		ast.NewIdent("value"),
	)
	keyValue := expr.(*ast.KeyValueExpr)

	assert.Equal(t, "key", keyValue.Key.(*ast.Ident).Name)
	assert.Equal(t, "value", keyValue.Value.(*ast.Ident).Name)
}

func TestGetNameExpr(t *testing.T) {
	expr := GetNameExpr("name")
	ident := expr.(*ast.Ident)

	assert.Equal(t, "name", ident.Name)
}

func TestGetAssignStmt(t *testing.T) {
	stmt := GetAssignStmt(
		Exprs{ast.NewIdent("left")},
		Exprs{ast.NewIdent("right")},
	)
	assign := stmt.(*ast.AssignStmt)

	assert.Equal(t, token.ASSIGN, assign.Tok)
	assert.Equal(t, "left", assign.Lhs[0].(*ast.Ident).Name)
	assert.Equal(t, "right", assign.Rhs[0].(*ast.Ident).Name)
}

func TestGetAssignDefineStmt(t *testing.T) {
	stmt := GetAssignDefineStmt(
		Exprs{ast.NewIdent("left")},
		Exprs{ast.NewIdent("right")},
	)
	assign := stmt.(*ast.AssignStmt)

	assert.Equal(t, token.DEFINE, assign.Tok)
	assert.Equal(t, "left", assign.Lhs[0].(*ast.Ident).Name)
	assert.Equal(t, "right", assign.Rhs[0].(*ast.Ident).Name)
}

func TestGetStructStmt(t *testing.T) {
	stmt := GetStructStmt("name", Fields{&ast.Field{Names: []*ast.Ident{ast.NewIdent("fieldName")}}})
	decl := stmt.(*ast.DeclStmt)

	assert.Equal(t, token.TYPE, decl.Decl.(*ast.GenDecl).Tok)
	assert.Equal(t, "name", decl.Decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name)
	assert.Equal(t, "fieldName", decl.Decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0].Names[0].Name)
}

func TestGetIfStmt(t *testing.T) {
	stmt := GetIfStmt(
		ast.NewIdent("cond"),
		Stmts{GetExprStmt(ast.NewIdent("bodyStmt"))},
	)
	ifstmt := stmt.(*ast.IfStmt)

	assert.Equal(t, "cond", ifstmt.Cond.(*ast.Ident).Name)
	assert.Equal(t, "bodyStmt", ifstmt.Body.List[0].(*ast.ExprStmt).X.(*ast.Ident).Name)
}

func TestGetExprStmt(t *testing.T) {
	stmt := GetExprStmt(ast.NewIdent("expr"))
	expr := stmt.(*ast.ExprStmt)

	assert.Equal(t, "expr", expr.X.(*ast.Ident).Name)
}

func TestGetVarDecl(t *testing.T) {
	decl := GetVarDecl("name", "typeExpr")
	varDecl := decl.(*ast.GenDecl)

	assert.Equal(t, token.VAR, varDecl.Tok)
	assert.Equal(t, "name", varDecl.Specs[0].(*ast.ValueSpec).Names[0].Name)
	assert.Equal(t, "typeExpr", varDecl.Specs[0].(*ast.ValueSpec).Type.(*ast.Ident).Name)
}

func TestGetFuncDecl(t *testing.T) {
	decl := GetFuncDecl(
		"name",
		Fields{{Names: []*ast.Ident{ast.NewIdent("fieldName")}}},
		Stmts{GetExprStmt(ast.NewIdent("bodyStmt"))},
	)
	funcDecl := decl.(*ast.FuncDecl)

	assert.Equal(t, "name", funcDecl.Name.Name)
	assert.Equal(t, "fieldName", funcDecl.Type.Params.List[0].Names[0].Name)
	assert.Equal(t, "bodyStmt", funcDecl.Body.List[0].(*ast.ExprStmt).X.(*ast.Ident).Name)
}
