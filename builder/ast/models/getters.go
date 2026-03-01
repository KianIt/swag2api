package models

import (
	"go/ast"
	"go/token"
)

func GetField(name, typeExpr, tag string) Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  ast.NewIdent(typeExpr),
		Tag:   GetStringBasicLit(tag),
	}
}

func GetNeqExpr(left, right Expr) Expr {
	return &ast.BinaryExpr{
		Op: token.NEQ,
		X:  left,
		Y:  right,
	}
}

func GetCallExpr(caller Expr, params Exprs) Expr {
	return &ast.CallExpr{
		Fun:  caller,
		Args: params.Ast(),
	}
}

func GetSelectorExpr(source Expr, item string) Expr {
	return &ast.SelectorExpr{
		X:   source,
		Sel: ast.NewIdent(item),
	}
}

func GetNameExpr(name string) Expr {
	return ast.NewIdent(name)
}

func GetAssignStmt(left, right Exprs) Stmt {
	return getAssignStmt(left, right, token.ASSIGN)
}

func GetAssignDefineStmt(left, right Exprs) Stmt {
	return getAssignStmt(left, right, token.DEFINE)
}

func getAssignStmt(left, right Exprs, tok token.Token) Stmt {
	return &ast.AssignStmt{
		Lhs: left.Ast(),
		Tok: tok,
		Rhs: right.Ast(),
	}
}

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

func GetIfStmt(cond Expr, bodyStmts Stmts) Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{
			List: bodyStmts.Ast(),
		},
	}
}

func GetExprStmt(expr Expr) Stmt {
	return &ast.ExprStmt{
		X: expr,
	}
}

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

func GetStringBasicLit(value string) BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

func GetStructLitExpr(name string, elements Exprs) Expr {
	return &ast.CompositeLit{
		Type: ast.NewIdent(name),
		Elts: elements.Ast(),
	}
}

func GetKeyValueExpr(key, value Expr) Expr {
	return &ast.KeyValueExpr{
		Key:   key,
		Value: value,
	}
}
