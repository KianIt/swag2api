package builder

import (
	"fmt"
	"go/ast"
	"go/token"
	"net/http"
	"os"
	"testing"

	astm "github.com/KianIt/swag2api/builder/ast/models"
	templateModels "github.com/KianIt/swag2api/builder/ast/template/models"
	"github.com/KianIt/swag2api/builder/models"
	s2aModels "github.com/KianIt/swag2api/models"
	parserModels "github.com/KianIt/swag2api/parser/models"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
)

func TestGetHandleResultStmt(t *testing.T) {
	stmt := getHandleResultStmt("errName")
	assert.Equal(
		t,
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.Ident{Name: templateModels.HandleResult},
				Args: []ast.Expr{
					&ast.Ident{Name: models.WToken.String()},
					&ast.Ident{Name: "errName"},
					&ast.Ident{Name: models.ResultValueToken.String()},
				},
			},
		},
		stmt,
	)
}

func TestGetHandleBadRequestStmt(t *testing.T) {
	stmt := getHandleBadRequestStmt("errName")
	assert.Equal(
		t,
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X:  &ast.Ident{Name: "errName"},
				Y:  &ast.Ident{Name: models.NilToken.String()},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
							Args: []ast.Expr{
								&ast.Ident{Name: models.WToken.String()},
								&ast.Ident{Name: "errName"},
							},
						},
					},
					&ast.ExprStmt{
						X: &ast.Ident{Name: models.ReturnToken.String()},
					},
				},
			},
		},
		stmt,
	)
}

func TestGetFuncHandlerResultStmts(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		stmts := getFuncHandlerResultStmts(
			s2aModels.Function{
				Results: s2aModels.Results{
					{
						Field: s2aModels.Field{
							Name:     "resultString",
							TypeExpr: "string",
						},
					},
					{
						Field: s2aModels.Field{
							Name:     "resultInt",
							TypeExpr: "int",
						},
					},
				},
			},
		)
		assert.Equal(
			t,
			astm.Stmts{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: &ast.Ident{Name: models.ResultTypeToken.String()},
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{{Name: "ResultString"}},
												Type:  &ast.Ident{Name: "string"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"resultString\"`",
												},
											},
											{
												Names: []*ast.Ident{{Name: "ResultInt"}},
												Type:  &ast.Ident{Name: "int"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"resultInt\"`",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.ResultValueToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{Name: models.ResultTypeToken.String()},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "ResultString"},
									Value: &ast.Ident{Name: "resultString"},
								},
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "ResultInt"},
									Value: &ast.Ident{Name: "resultInt"},
								},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: templateModels.HandleResult},
						Args: []ast.Expr{
							&ast.Ident{Name: models.WToken.String()},
							&ast.Ident{Name: models.NilToken.String()},
							&ast.Ident{Name: models.ResultValueToken.String()},
						},
					},
				},
			},
			stmts,
		)
	})

	t.Run("has error", func(t *testing.T) {
		stmts := getFuncHandlerResultStmts(
			s2aModels.Function{
				Results: s2aModels.Results{
					{
						Field: s2aModels.Field{
							Name:     "resultString",
							TypeExpr: "string",
						},
					},
					{
						Field: s2aModels.Field{
							Name:     "resultError",
							TypeExpr: "error",
						},
					},
				},
			},
		)
		assert.Equal(
			t,
			astm.Stmts{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: &ast.Ident{Name: models.ResultTypeToken.String()},
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{{Name: "ResultString"}},
												Type:  &ast.Ident{Name: "string"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"resultString\"`",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.ResultValueToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{Name: models.ResultTypeToken.String()},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "ResultString"},
									Value: &ast.Ident{Name: "resultString"},
								},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: templateModels.HandleResult},
						Args: []ast.Expr{
							&ast.Ident{Name: models.WToken.String()},
							&ast.Ident{Name: "resultError"},
							&ast.Ident{Name: models.ResultValueToken.String()},
						},
					},
				},
			},
			stmts,
		)
	})
}

func TestGetFuncHandlerFuncCallStmts(t *testing.T) {
	f := s2aModels.Function{
		Name: "function",
		Params: s2aModels.Params{
			{
				Field: s2aModels.Field{
					Name:     "paramStringPath",
					TypeExpr: "string",
				},
				Type:   s2aModels.String,
				Origin: s2aModels.Path,
			},
			{
				Field: s2aModels.Field{
					Name:     "paramIntQuery",
					TypeExpr: "int",
				},
				Type:   s2aModels.Int,
				Origin: s2aModels.Query,
			},
			{
				Field: s2aModels.Field{
					Name:     "paramBoolBody",
					TypeExpr: "bool",
				},
				Type:   s2aModels.Bool,
				Origin: s2aModels.Body,
			},
		},
		Results: s2aModels.Results{
			{
				Field: s2aModels.Field{
					Name:     "resultString",
					TypeExpr: "string",
				},
			},
			{
				Field: s2aModels.Field{
					Name:     "resultError",
					TypeExpr: "error",
				},
			},
		},
	}

	t.Run("no body struct", func(t *testing.T) {
		stmts := getFuncHandlerFuncCallStmts(f, false)

		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "resultString"},
						&ast.Ident{Name: "resultError"},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "function"},
							Args: []ast.Expr{
								&ast.Ident{Name: "paramStringPath"},
								&ast.Ident{Name: "paramIntQuery"},
								&ast.Ident{Name: "paramBoolBody"},
							},
						},
					},
				},
			},
			stmts,
		)
	})

	t.Run("has body struct", func(t *testing.T) {
		stmts := getFuncHandlerFuncCallStmts(f, true)

		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "resultString"},
						&ast.Ident{Name: "resultError"},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "function"},
							Args: []ast.Expr{
								&ast.Ident{Name: "paramStringPath"},
								&ast.Ident{Name: "paramIntQuery"},
								&ast.SelectorExpr{
									X:   &ast.Ident{Name: models.BodyValueToken.String()},
									Sel: &ast.Ident{Name: "ParamBoolBody"},
								},
							},
						},
					},
				},
			},
			stmts,
		)
	})
}

func TestGetFuncHandlerParamUnmarshalStmt(t *testing.T) {
	paramNameOrigin := "paramPath"
	paramName := "param"
	typeExpr := "string"

	t.Run("not body", func(t *testing.T) {
		stmt, errName := getFuncHandlerParamUnmarshalStmt(paramNameOrigin, paramName, typeExpr, false)

		assert.Equal(
			t,
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					&ast.Ident{Name: paramName},
					&ast.Ident{Name: paramName + models.UnmarshalErrSuffixToken.String()},
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{Name: templateModels.UnmarshalString + fmt.Sprintf("[%s]", typeExpr)},
						Args: []ast.Expr{
							&ast.Ident{Name: paramNameOrigin},
						},
					},
				},
			},
			stmt,
		)
		assert.Equal(t, paramName+models.UnmarshalErrSuffixToken.String(), errName)
	})

	t.Run("body", func(t *testing.T) {
		stmt, errName := getFuncHandlerParamUnmarshalStmt(paramNameOrigin, paramName, typeExpr, true)

		assert.Equal(
			t,
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					&ast.Ident{Name: paramName},
					&ast.Ident{Name: paramName + models.UnmarshalErrSuffixToken.String()},
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{Name: templateModels.UnmarshalBytes + fmt.Sprintf("[%s]", typeExpr)},
						Args: []ast.Expr{
							&ast.Ident{Name: paramNameOrigin},
						},
					},
				},
			},
			stmt,
		)
		assert.Equal(t, paramName+models.UnmarshalErrSuffixToken.String(), errName)
	})
}

func TestGetFuncHandlerBodyParamStmts(t *testing.T) {
	t.Run("one param", func(t *testing.T) {
		stmts, isStruct := getFuncHandlerBodyParamStmts(s2aModels.Params{
			{
				Field: s2aModels.Field{
					Name:     "param",
					TypeExpr: "string",
				},
				Type:   s2aModels.String,
				Origin: s2aModels.Body,
			},
		})

		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyToken.String()},
						&ast.Ident{Name: models.BodyReadErrToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.Ident{Name: models.BodyGetToken.String()},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyReadErrToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyReadErrToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "param"},
						&ast.Ident{Name: "param" + models.UnmarshalErrSuffixToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.UnmarshalBytes + "[string]"},
							Args: []ast.Expr{
								&ast.Ident{Name: models.BodyToken.String()},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: "param" + models.UnmarshalErrSuffixToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: "param" + models.UnmarshalErrSuffixToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
			},
			stmts,
		)
		assert.False(t, isStruct)
	})

	t.Run("two params", func(t *testing.T) {
		stmts, isStruct := getFuncHandlerBodyParamStmts(s2aModels.Params{
			{
				Field: s2aModels.Field{
					Name:     "param1",
					TypeExpr: "string",
				},
				Type:   s2aModels.String,
				Origin: s2aModels.Body,
			},
			{
				Field: s2aModels.Field{
					Name:     "param2",
					TypeExpr: "int",
				},
				Type:   s2aModels.Int,
				Origin: s2aModels.Body,
			},
		})

		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyToken.String()},
						&ast.Ident{Name: models.BodyReadErrToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.Ident{Name: models.BodyGetToken.String()},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyReadErrToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyReadErrToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: &ast.Ident{Name: models.BodyTypeToken.String()},
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{{Name: "Param1"}},
												Type:  &ast.Ident{Name: "string"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"param1\"`",
												},
											},
											{
												Names: []*ast.Ident{{Name: "Param2"}},
												Type:  &ast.Ident{Name: "int"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"param2\"`",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyValueToken.String()},
						&ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.UnmarshalBytes + fmt.Sprintf("[%s]", models.BodyTypeToken.String())},
							Args: []ast.Expr{
								&ast.Ident{Name: models.BodyToken.String()},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
			},
			stmts,
		)
		assert.True(t, isStruct)
	})
}

func TestBuilder_addImport(t *testing.T) {
	builder := NewBuilder(parserModels.ParsingInfo{})

	builder.addImport("encoding/json")
	builder.addImport("net/http")

	file := builder.fc.Compile()
	assert.Equal(
		t,
		[]ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: &ast.Ident{},
						Path: &ast.BasicLit{
							Value: "\"encoding/json\"",
							Kind:  token.STRING,
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{},
						Path: &ast.BasicLit{
							Value: "\"net/http\"",
							Kind:  token.STRING,
						},
					},
				},
				Lparen: 0,
				Rparen: 10,
			},
		},
		file.Decls,
	)
}

func TestBuilder_addAliasedImport(t *testing.T) {
	builder := NewBuilder(parserModels.ParsingInfo{})

	builder.addAliasedImport("encoding/json", "json")
	builder.addAliasedImport("net/http", "http")

	file := builder.fc.Compile()
	assert.Equal(
		t,
		[]ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: &ast.Ident{Name: "json"},
						Path: &ast.BasicLit{
							Value: "\"encoding/json\"",
							Kind:  token.STRING,
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{Name: "http"},
						Path: &ast.BasicLit{
							Value: "\"net/http\"",
							Kind:  token.STRING,
						},
					},
				},
				Lparen: 0,
				Rparen: 10,
			},
		},
		file.Decls,
	)
}

func TestBuilder_addHTTPHandler(t *testing.T) {
	t.Run("HTTP handler exists", func(t *testing.T) {
		builder := NewBuilder(parserModels.ParsingInfo{
			HTTPHandler: parserModels.HTTPHandlerInfo{
				Name:   "handler",
				Exists: true,
			},
		})

		builder.addHTTPHandler()

		file := builder.fc.Compile()
		assert.Equal(t, []ast.Decl{}, file.Decls)
	})

	t.Run("HTTP handler doesn't exist", func(t *testing.T) {
		builder := NewBuilder(parserModels.ParsingInfo{
			HTTPHandler: parserModels.HTTPHandlerInfo{
				Name:   "handler",
				Exists: false,
			},
		})

		builder.addHTTPHandler()

		file := builder.fc.Compile()
		assert.Equal(
			t,
			[]ast.Decl{
				&ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{{Name: "handler"}},
							Type:  &ast.Ident{Name: models.HandlerToken.String()},
						},
					},
				},
			},
			file.Decls,
		)
	})
}

func TestBuilder_addTemplates(t *testing.T) {
	builder := NewBuilder(parserModels.ParsingInfo{})

	_ = builder.tm.Load()
	builder.addTemplates()

	file := builder.fc.Compile()
	assert.Len(t, file.Decls, len(templateModels.TemplateNames))
}

func TestBuilder_addInit(t *testing.T) {
	builder := NewBuilder(parserModels.ParsingInfo{
		Funcs: s2aModels.Functions{
			{
				Name: "function",
				Endpoint: s2aModels.Endpoint{
					Method: http.MethodGet,
					Path:   "path-to-function",
				},
			},
		},
		HTTPHandler: parserModels.HTTPHandlerInfo{
			Name: "handler",
		},
	})

	builder.addInit()

	file := builder.fc.Compile()
	assert.Equal(
		t,
		[]ast.Decl{
			&ast.FuncDecl{
				Name: &ast.Ident{Name: models.InitToken.String()},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								&ast.Ident{Name: models.MuxToken.String()},
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{Name: models.NewMuxToken.String()},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: models.MuxToken.String()},
									Sel: &ast.Ident{Name: models.HandleFuncToken.String()},
								},
								Args: []ast.Expr{
									&ast.Ident{Name: "\"GET path-to-function\""},
									&ast.Ident{Name: "_handler_function"},
								},
							},
						},
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{
								&ast.Ident{Name: "handler"},
							},
							Rhs: []ast.Expr{
								&ast.Ident{Name: models.MuxToken.String()},
							},
						},
					},
				},
			},
		},
		file.Decls,
	)
}

func TestBuilder_getFuncHandlerParamStmts(t *testing.T) {
	t.Run("zero body params", func(t *testing.T) {
		builder := NewBuilder(parserModels.ParsingInfo{})

		f := s2aModels.Function{
			Name: "function",
			Params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramString",
						TypeExpr: "string",
					},
					Type:   s2aModels.String,
					Origin: s2aModels.Path,
				},
			},
		}

		stmts, isBodyStruct := builder.getFuncHandlerParamStmts(f)
		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "paramStringPath"},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: models.PathValueGetToken.String()},
							Args: []ast.Expr{
								&ast.Ident{Name: "\"paramString\""},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "paramString"},
						&ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.UnmarshalString + "[string]"},
							Args: []ast.Expr{
								&ast.Ident{Name: "paramStringPath"},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
			},
			stmts,
		)
		assert.False(t, isBodyStruct)
	})

	t.Run("one body param", func(t *testing.T) {
		builder := NewBuilder(parserModels.ParsingInfo{})

		f := s2aModels.Function{
			Name: "function",
			Params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramString",
						TypeExpr: "string",
					},
					Type:   s2aModels.String,
					Origin: s2aModels.Body,
				},
			},
		}

		stmts, isBodyStruct := builder.getFuncHandlerParamStmts(f)
		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyToken.String()},
						&ast.Ident{Name: models.BodyReadErrToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.Ident{Name: models.BodyGetToken.String()},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyReadErrToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyReadErrToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: "paramString"},
						&ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.UnmarshalBytes + "[string]"},
							Args: []ast.Expr{
								&ast.Ident{Name: models.BodyToken.String()},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: "paramString" + models.UnmarshalErrSuffixToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
			},
			stmts,
		)
		assert.False(t, isBodyStruct)
	})

	t.Run("multiple body params", func(t *testing.T) {
		builder := NewBuilder(parserModels.ParsingInfo{})

		f := s2aModels.Function{
			Name: "function",
			Params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramString",
						TypeExpr: "string",
					},
					Type:   s2aModels.String,
					Origin: s2aModels.Body,
				},
				{
					Field: s2aModels.Field{
						Name:     "paramInt",
						TypeExpr: "int",
					},
					Type:   s2aModels.Int,
					Origin: s2aModels.Body,
				},
			},
		}

		stmts, isBodyStruct := builder.getFuncHandlerParamStmts(f)
		assert.Equal(
			t,
			astm.Stmts{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyToken.String()},
						&ast.Ident{Name: models.BodyReadErrToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.Ident{Name: models.BodyGetToken.String()},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyReadErrToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyReadErrToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: &ast.Ident{Name: models.BodyTypeToken.String()},
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{{Name: "ParamString"}},
												Type:  &ast.Ident{Name: "string"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"paramString\"`",
												},
											},
											{
												Names: []*ast.Ident{{Name: "ParamInt"}},
												Type:  &ast.Ident{Name: "int"},
												Tag: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "`json:\"paramInt\"`",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: models.BodyValueToken.String()},
						&ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: templateModels.UnmarshalBytes + fmt.Sprintf("[%s]", models.BodyTypeToken.String())},
							Args: []ast.Expr{
								&ast.Ident{Name: models.BodyToken.String()},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.NEQ,
						X:  &ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
						Y:  &ast.Ident{Name: models.NilToken.String()},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: templateModels.HandleBadRequest},
									Args: []ast.Expr{
										&ast.Ident{Name: models.WToken.String()},
										&ast.Ident{Name: models.BodyValueToken.String() + models.UnmarshalErrSuffixToken.String()},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.Ident{Name: models.ReturnToken.String()},
							},
						},
					},
				},
			},
			stmts,
		)
		assert.True(t, isBodyStruct)
	})
}

func TestBuilder_Build(t *testing.T) {
	_, _ = os.Create("./testdata/generated_actual.go")
	// defer func() {
	// 	_ = os.Remove("./testdata/generated_actual.go")
	// }()

	builder := NewBuilder(parserModels.ParsingInfo{
		PkgName: "testdata",
		Imports: []s2aModels.Import{
			{
				Path:  "github.com/KianIt/swag2api/statuses",
				Alias: "s2aStatuses",
			},
		},
		Funcs: s2aModels.Functions{
			{
				Name: "function1",
				Params: s2aModels.Params{
					{
						Field: s2aModels.Field{
							Name:     "paramString",
							TypeExpr: "string",
						},
						Type:   s2aModels.String,
						Origin: s2aModels.Path,
					},
					{
						Field: s2aModels.Field{
							Name:     "paramInt",
							TypeExpr: "int",
						},
						Type:   s2aModels.Int,
						Origin: s2aModels.Query,
					},
				},
				Results: s2aModels.Results{
					{
						Field: s2aModels.Field{
							Name:     "resultString",
							TypeExpr: "string",
						},
					},
					{
						Field: s2aModels.Field{
							Name:     "resultError",
							TypeExpr: "error",
						},
					},
				},
				Endpoint: s2aModels.Endpoint{
					Path:   "path-to-function1",
					Method: http.MethodGet,
				},
			},
			{
				Name: "function2",
				Params: s2aModels.Params{
					{
						Field: s2aModels.Field{
							Name:     "paramString",
							TypeExpr: "string",
						},
						Type:   s2aModels.String,
						Origin: s2aModels.Body,
					},
				},
				Results: s2aModels.Results{
					{
						Field: s2aModels.Field{
							Name:     "resultString",
							TypeExpr: "string",
						},
					},
					{
						Field: s2aModels.Field{
							Name:     "_",
							TypeExpr: "error",
						},
					},
				},
				Endpoint: s2aModels.Endpoint{
					Path:   "path-to-function2",
					Method: http.MethodPost,
				},
			},
			{
				Name: "function3",
				Params: s2aModels.Params{
					{
						Field: s2aModels.Field{
							Name:     "paramString",
							TypeExpr: "string",
						},
						Type:   s2aModels.String,
						Origin: s2aModels.Body,
					},
					{
						Field: s2aModels.Field{
							Name:     "paramInt",
							TypeExpr: "int",
						},
						Type:   s2aModels.Int,
						Origin: s2aModels.Body,
					},
				},
				Results: s2aModels.Results{
					{
						Field: s2aModels.Field{
							Name:     "resultError",
							TypeExpr: "bool",
						},
					},
				},
				Endpoint: s2aModels.Endpoint{
					Path:   "path-to-function3",
					Method: http.MethodPatch,
				},
			},
		},
		HTTPHandler: parserModels.HTTPHandlerInfo{
			Name:   "handler",
			Exists: false,
		},
	})

	require.NoError(t, builder.Build("./testdata/generated_actual.go"))

	expectedData, _ := os.ReadFile("./testdata/generated_expected.go")
	actualData, _ := os.ReadFile("./testdata/generated_actual.go")

	assert.Equal(t, expectedData, actualData)
}
