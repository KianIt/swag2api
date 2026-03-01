package ast

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/KianIt/swag2api/builder/ast/models"
	"github.com/stretchr/testify/assert"
)

func TestFileCompiler(t *testing.T) {
	fc := NewFileCompiler("pkg")

	fc.AddImport("path1", "alias1")
	fc.AddImport("path2", "alias2")

	fc.AddDecl(models.GetVarDecl("var1", "type1"))
	fc.AddDecl(models.GetVarDecl("var2", "type2"))

	file := fc.Compile()
	assert.Equal(t,
		[]ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Name: ast.NewIdent("alias1"), Path: models.GetStringBasicLit("\"path1\"")},
					&ast.ImportSpec{Name: ast.NewIdent("alias2"), Path: models.GetStringBasicLit("\"path2\"")},
				},
				Lparen: 0,
				Rparen: 7,
			},
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("var1")},
						Type:  ast.NewIdent("type1"),
					},
				},
			},
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("var2")},
						Type:  ast.NewIdent("type2"),
					},
				},
			},
		},
		file.Decls,
	)
}
