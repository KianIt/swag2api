package filecompiler

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/KianIt/swag2api/builder/ast/models"
)

// FileCompiler is an AST file compiler.
//
// Turns AST declaration into a single file.
type FileCompiler struct {
	file    *ast.File
	imports []ast.Spec
	decls   []ast.Decl
}

// NewFileCompiler returns a new file compiler.
func NewFileCompiler(pkgName string) *FileCompiler {
	return &FileCompiler{
		file: &ast.File{
			Name:  ast.NewIdent(pkgName),
			Decls: make([]ast.Decl, 0),
			Scope: ast.NewScope(nil),
		},
	}
}

// AddImport adds an import.
func (fc *FileCompiler) AddImport(path string, alias string) {
	spec := &ast.ImportSpec{
		Name: ast.NewIdent(alias),
		Path: models.GetStringBasicLit(strconv.Quote(path)),
	}
	fc.imports = append(fc.imports, spec)
}

// AddDecl adds a declaration.
func (fc *FileCompiler) AddDecl(decl models.Decl) {
	fc.decls = append(fc.decls, decl)
}

// Compile compiles a file.
func (fc *FileCompiler) Compile() *ast.File {
	if len(fc.imports) > 0 {
		fc.compileImports()
	}

	for _, decl := range fc.decls {
		fc.compileDecl(decl)
	}

	return fc.file
}

// compileImports compiles file imports.
func (fc *FileCompiler) compileImports() {
	lPar, rPar := token.NoPos, token.NoPos
	if len(fc.imports) > 1 {
		lPar = fc.imports[0].Pos()
		rPar = fc.imports[len(fc.imports)-1].End()
	}

	decl := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: lPar,
		Specs:  fc.imports,
		Rparen: rPar,
	}

	fc.compileDecl(decl)
}

// compileDecl compiles file declarations.
func (fc *FileCompiler) compileDecl(decl models.Decl) {
	fc.file.Decls = append(fc.file.Decls, decl)
}
