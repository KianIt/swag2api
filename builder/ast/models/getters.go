package models

import (
	"go/ast"
	"go/token"
)

// GetField returns as AST field.
func GetField(name, typeExpr, tag string) Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  ast.NewIdent(typeExpr),
		Tag:   GetStringBasicLit(tag),
	}
}

// GetStringBasicLit returns an AST string literal.
func GetStringBasicLit(value string) BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

// GetNeqExpr returns as AST Not Equal expression.
func GetNeqExpr(left, right Expr) Expr {
	return &ast.BinaryExpr{
		Op: token.NEQ,
		X:  left,
		Y:  right,
	}
}

// GetCallExpr returns as AST Calling expression.
func GetCallExpr(caller Expr, params Exprs) Expr {
	return &ast.CallExpr{
		Fun:  caller,
		Args: params.Ast(),
	}
}

// GetSelectorExpr returns an AST Selector expression.
func GetSelectorExpr(source Expr, item string) Expr {
	return &ast.SelectorExpr{
		X:   source,
		Sel: ast.NewIdent(item),
	}
}

// GetStructLitExpr returns an AST struct literal.
func GetStructLitExpr(name string, elements Exprs) Expr {
	return &ast.CompositeLit{
		Type: ast.NewIdent(name),
		Elts: elements.Ast(),
	}
}

// GetKeyValueExpr returns an AST Key-Value expression.
func GetKeyValueExpr(key, value Expr) Expr {
	return &ast.KeyValueExpr{
		Key:   key,
		Value: value,
	}
}

// GetNameExpr returns an AST Name expression.
func GetNameExpr(name string) Expr {
	return ast.NewIdent(name)
}

// GetAssignStmt returns an AST Assign statement.
func GetAssignStmt(left, right Exprs) Stmt {
	return getAssignStmt(left, right, token.ASSIGN)
}

// GetAssignDefineStmt returns an AST Assign (Define) statement.
func GetAssignDefineStmt(left, right Exprs) Stmt {
	return getAssignStmt(left, right, token.DEFINE)
}

// getAssignStmt returns an AST assign statement.
func getAssignStmt(left, right Exprs, tok token.Token) Stmt {
	return &ast.AssignStmt{
		Lhs: left.Ast(),
		Tok: tok,
		Rhs: right.Ast(),
	}
}

// GetStructStmt returns an AST Struct statement.
func GetStructStmt(name string, fields Fields) Stmt {
	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(name),
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: fields.Ast(),
						},
					},
				},
			},
		},
	}
}

// GetIfStmt returns an AST If statement.
func GetIfStmt(cond Expr, bodyStmts Stmts) Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{
			List: bodyStmts.Ast(),
		},
	}
}

// GetExprStmt returns an AST expression as statement.
func GetExprStmt(expr Expr) Stmt {
	return &ast.ExprStmt{
		X: expr,
	}
}

// GetVarDecl returns an AST Variable declaration.
func GetVarDecl(name, typeExpr string) Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent(name)},
				Type:  ast.NewIdent(typeExpr),
			},
		},
	}
}

// GetFuncDecl returns an AST Function declaration.
func GetFuncDecl(name string, params Fields, bodyStmts Stmts) Decl {
	return &ast.FuncDecl{
		Type: &ast.FuncType{
			Func: token.Pos(0),
			Params: &ast.FieldList{
				List: params.Ast(),
			},
		},
		Name: ast.NewIdent(name),
		Body: &ast.BlockStmt{
			List: bodyStmts.Ast(),
		},
	}
}
