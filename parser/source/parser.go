package source

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	s2aModels "github.com/KianIt/swag2api/models"
	"github.com/KianIt/swag2api/parser/models"
	"github.com/KianIt/swag2api/utils"
)

// Parser is a source code parser.
//
// Uses the built-in AST features to parse the functions.
// Also parses the main package name, file imports,
// API HTTP handler information.
type Parser struct {
	fs          *token.FileSet
	PkgName     string
	Imports     []s2aModels.Import
	Funcs       s2aModels.Functions
	HTTPHandler models.HTTPHandlerInfo
}

// NewSourceParser returns a new source code parser.
func NewSourceParser() *Parser {
	return &Parser{
		fs:      token.NewFileSet(),
		Imports: make([]s2aModels.Import, 0),
		Funcs:   make(s2aModels.Functions, 0),
	}
}

// Parse runs parsing.
//
// Walks over every file in the package (and subpackages)
// and reads definitions of all the visited functions.
func (p *Parser) Parse(pkgPath, handlerName string) error {
	log.Printf("Parsing source code from '%s'", pkgPath)

	p.HTTPHandler.Name = handlerName

	// Getting a file list.
	files, err := p.parseFiles(pkgPath)
	if err != nil {
		return fmt.Errorf("files: %w", err)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name.Name < files[j].Name.Name })

	// Parsing every file.
	for _, file := range files {
		ast.Walk(p, file)
	}

	sort.Slice(p.Funcs, func(i, j int) bool { return p.Funcs[i].Name < p.Funcs[j].Name })

	return nil
}

// parseFiles parses the source code package and returns a list of files.
func (p *Parser) parseFiles(pkgPath string) ([]*ast.File, error) {
	files := make([]*ast.File, 0)

	if err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, _ error) error {
		// Skipping not Golang source code files.
		if info.IsDir() || !utils.IsGoSource(info.Name()) {
			return nil
		}

		// Parsing a source code file.
		file, parseErr := parser.ParseFile(p.fs, path, nil, parser.ParseComments)
		if parseErr != nil {
			return fmt.Errorf("file '%s': %w", path, parseErr)
		}

		// File path relative to packagep path.
		relPath, pathErr := filepath.Rel(pkgPath, path)
		if pathErr != nil {
			return fmt.Errorf("file '%s' rel path: %w", path, pathErr)
		}

		// If relative path is package root then it's a main package.
		if p.PkgName == "" && relPath == filepath.Base(relPath) {
			p.PkgName = file.Name.Name
		}

		files = append(files, file)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// Visit implements the [ast.Visitor] interface.
//
// Allows the SourceParser to walk over nodes in a parsed file.
func (p *Parser) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch decl := node.(type) {
	case *ast.File:
		// Skipping files from non-main packages.
		if decl.Name.Name != p.PkgName {
			return nil
		}

		return p
	case *ast.GenDecl:
		// Trying to update the HTTP handler information.
		p.checkHTTPHandler(decl)
		// Parsing file imports.
		p.parseImports(decl)

		return p
	case *ast.FuncDecl:
		// Parsing a function.
		if err := p.parseFuncDecl(decl); err != nil {
			log.Printf("Error: parsing function '%s' declaration failed: %v", decl.Name.Name, err)
		}

		return nil
	}

	return nil
}

// checkHTTPHandler tries to update the HTTP handler information.
func (p *Parser) checkHTTPHandler(decl *ast.GenDecl) {
	if decl == nil || decl.Tok != token.VAR {
		return
	}

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, ident := range valueSpec.Names {
				if ident.Name == p.HTTPHandler.Name {
					p.HTTPHandler.Exists = true
					return
				}
			}
		}
	}
}

// parseImports tries to parse file imports.
func (p *Parser) parseImports(decl *ast.GenDecl) {
	if decl == nil || decl.Tok != token.IMPORT {
		return
	}

	for _, spec := range decl.Specs {
		if importSpec, ok := spec.(*ast.ImportSpec); ok {
			name, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				log.Printf("Error: couldn't unquote import '%s'", importSpec.Path.Value)
				continue
			}

			alias := ""
			if importSpec.Name != nil {
				alias = importSpec.Name.Name
			}

			p.Imports = append(p.Imports, s2aModels.Import{
				Path:  name,
				Alias: alias,
			})
		}
	}
}

// parseFuncDecl parsees a function declaration.
func (p *Parser) parseFuncDecl(decl *ast.FuncDecl) error {
	if decl == nil {
		return errors.New("declaration is nil")
	}

	log.Printf("Parsing function '%s'", decl.Name.Name)

	// Params.
	params, err := p.fieldList2Params(decl.Type.Params)
	if err != nil {
		return fmt.Errorf("params: %w", err)
	}

	// Results.
	results, err := p.fieldList2Results(decl.Type.Results)
	if err != nil {
		return fmt.Errorf("results: %w", err)
	}

	p.Funcs = append(p.Funcs, s2aModels.Function{
		Name:    decl.Name.Name,
		Params:  params,
		Results: results,
	})

	return nil
}

// fieldList2Params converts an *[ast.FieldList] into a [s2aModels.Params].
func (p *Parser) fieldList2Params(fieldList *ast.FieldList) (s2aModels.Params, error) {
	if fieldList == nil {
		return make(s2aModels.Params, 0), nil
	}

	params := make(s2aModels.Params, 0)
	for i, field := range fieldList.List {
		fieldParams, err := p.field2Params(field)
		if err != nil {
			return nil, fmt.Errorf("field #%d: %w", i, err)
		}

		params = append(params, fieldParams...)
	}

	return params, nil
}

// field2Params converts an *[ast.Field] into a [s2aModels.Params].
func (p *Parser) field2Params(field *ast.Field) (s2aModels.Params, error) {
	if field == nil {
		return make(s2aModels.Params, 0), nil
	}

	typeExpr, paramType, err := p.getType(field.Type)
	if err != nil {
		return nil, fmt.Errorf("type: %w", err)
	}

	// if no names, then it's an unnamed field.
	if len(field.Names) == 0 {
		return s2aModels.Params{{Field: s2aModels.Field{Name: "", TypeExpr: typeExpr}, Type: paramType}}, nil
	}

	// Several fields with different names and the same type.
	params := make(s2aModels.Params, 0, len(field.Names))
	for _, name := range field.Names {
		params = append(params, s2aModels.Param{Field: s2aModels.Field{Name: name.Name, TypeExpr: typeExpr}, Type: paramType})
	}

	return params, nil
}

// fieldList2Results converts an *[ast.FieldList] into a [s2aModels.Results].
func (p *Parser) fieldList2Results(fieldList *ast.FieldList) (s2aModels.Results, error) {
	if fieldList == nil {
		return make(s2aModels.Results, 0), nil
	}

	results := make(s2aModels.Results, 0)
	idx := 0
	for _, field := range fieldList.List {
		fieldResults, err := p.field2Results(field, idx)
		if err != nil {
			return nil, fmt.Errorf("field #%d: %w", idx, err)
		}

		results = append(results, fieldResults...)
		idx += len(fieldResults)
	}

	return results, nil
}

// field2Results converts an *[ast.Field] into a s2aModels.Results
//
// Gives a default name to the unnamed results.
func (p *Parser) field2Results(field *ast.Field, idx int) (s2aModels.Results, error) {
	if field == nil {
		return make(s2aModels.Results, 0), nil
	}

	typeExpr, _, err := p.getType(field.Type)
	if err != nil {
		return nil, fmt.Errorf("type: %w", err)
	}

	// if no names, then it's an unnamed field.
	if len(field.Names) == 0 {
		return s2aModels.Results{{Field: s2aModels.Field{Name: fmt.Sprintf("res%d", idx), TypeExpr: typeExpr}}}, nil
	}

	// Several fields with different names and the same type.
	results := make(s2aModels.Results, 0, len(field.Names))
	for _, name := range field.Names {
		results = append(results, s2aModels.Result{Field: s2aModels.Field{Name: name.Name, TypeExpr: typeExpr}})
	}

	return results, nil
}

// getType parses a type expression and a [s2aModels.ParamType] from an [ast.Expr].
func (p *Parser) getType(expr ast.Expr) (string, s2aModels.ParamType, error) {
	var buf bytes.Buffer

	if err := printer.Fprint(&buf, token.NewFileSet(), expr); err != nil {
		return "", "", fmt.Errorf("print: %w", err)
	}

	typeExpr := buf.String()

	paramType, err := p.getParamType(typeExpr)
	if err != nil {
		return "", "", fmt.Errorf("param type: %w", err)
	}

	return typeExpr, paramType, nil
}

// getParamType parses a [s2aModels.ParamType] from a type expression.
func (p *Parser) getParamType(typeExpr string) (s2aModels.ParamType, error) {
	var paramType s2aModels.ParamType

	// Trying to parse a simple type.
	switch typeExpr {
	case "uint8", "int8", "uint16", "int16", "byte", "int32", "uint32", "rune", "uint64", "int64", "int", "uint":
		paramType = s2aModels.Int
	case "float32", "float64":
		paramType = s2aModels.Float
	case "bool":
		paramType = s2aModels.Bool
	case "string":
		paramType = s2aModels.String
	case "error":
		paramType = s2aModels.Error
	case "interface{}", "any":
		paramType = s2aModels.Any
	}

	if paramType != "" {
		return paramType, nil
	}

	// Prefix [] means a slice.
	if strings.HasPrefix(typeExpr, "[]") {
		// Trying to parse the slice item type.
		subType, err := p.getParamType(typeExpr[2:])
		if err != nil {
			return "", err
		}

		// Returning a slice.
		return subType.SliceOf(), nil
	}

	// Prefix map[ means a slice.
	if strings.HasPrefix(typeExpr, "map[") {
		_, after, found := strings.Cut(typeExpr, "]")

		if !found {
			return "", fmt.Errorf("invalid map type: '%s'", typeExpr)
		}

		// Trying to parse the map value type.
		subType, err := p.getParamType(after)
		if err != nil {
			return "", err
		}

		// Returning a slice.
		return subType.MapOf(), nil
	}

	// If custom type, then adding the package name.
	paramType = s2aModels.ParamType(typeExpr)
	if !strings.Contains(typeExpr, ".") {
		paramType = s2aModels.ParamType(p.PkgName) + "." + paramType
	}

	// Returning a custom type.
	return paramType.CustomOf(), nil
}
