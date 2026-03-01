package source

import (
	"go/ast"
	"go/token"
	"testing"

	s2aModels "github.com/KianIt/swag2api/models"
	"github.com/KianIt/swag2api/parser/models"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
)

func TestSourceParser_getParamType(t *testing.T) {
	tt := []struct {
		name      string
		typeExpr  string
		paramType s2aModels.ParamType
		isErr     bool
	}{
		{
			name:      "int type",
			typeExpr:  "int",
			paramType: s2aModels.Int,
		},
		{
			name:      "error type",
			typeExpr:  "error",
			paramType: s2aModels.Error,
		},
		{
			name:      "interface{} type",
			typeExpr:  "interface{}",
			paramType: s2aModels.Any,
		},
		{
			name:      "any type",
			typeExpr:  "any",
			paramType: s2aModels.Any,
		},
		{
			name:      "slice of string",
			typeExpr:  "[]string",
			paramType: s2aModels.ParamType("slice[]string"),
		},
		{
			name:      "map of string to int",
			typeExpr:  "map[string]int",
			paramType: s2aModels.ParamType("map[]int"),
		},
		{
			name:     "invalid map type",
			typeExpr: "map[invalid",
			isErr:    true,
		},
		{
			name:      "custom type with pkg",
			typeExpr:  "pkg.Type",
			paramType: s2aModels.ParamType("custom[]pkg.Type"),
		},
		{
			name:      "custom type without pkg",
			typeExpr:  "Type",
			paramType: s2aModels.ParamType("custom[]testPkg.Type"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()
			parser.PkgName = "testPkg"

			paramType, err := parser.getParamType(tc.typeExpr)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.paramType, paramType)
			}
		})
	}
}

func TestSourceParser_getType(t *testing.T) {
	parser := NewSourceParser()

	typeExpr, paramType, err := parser.getType(ast.NewIdent("[]pkg.Type"))
	require.NoError(t, err)
	assert.Equal(t, "[]pkg.Type", typeExpr)
	assert.Equal(t, s2aModels.ParamType("slice[]custom[]pkg.Type"), paramType)
}

func TestSourceParser_field2Results(t *testing.T) {
	tt := []struct {
		name    string
		field   *ast.Field
		idx     int
		results s2aModels.Results
		isErr   bool
	}{
		{
			name:    "nil field",
			field:   nil,
			idx:     0,
			results: s2aModels.Results{},
			isErr:   false,
		},
		{
			name: "unnamed field with string type",
			field: &ast.Field{
				Type: ast.NewIdent("string"),
			},
			idx: 0,
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "res0",
						TypeExpr: "string",
					},
				},
			},
			isErr: false,
		},
		{
			name: "named field with single name",
			field: &ast.Field{
				Names: []*ast.Ident{ast.NewIdent("result")},
				Type:  ast.NewIdent("string"),
			},
			idx: 0,
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "result",
						TypeExpr: "string",
					},
				},
			},
			isErr: false,
		},
		{
			name: "named field with multiple names",
			field: &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("result1"),
					ast.NewIdent("result2"),
				},
				Type: ast.NewIdent("int"),
			},
			idx: 0,
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "result1",
						TypeExpr: "int",
					},
				},
				{
					Field: s2aModels.Field{
						Name:     "result2",
						TypeExpr: "int",
					},
				},
			},
			isErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			results, err := parser.field2Results(tc.field, tc.idx)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.results, results)
			}
		})
	}
}

func TestSourceParser_fieldList2Results(t *testing.T) {
	tt := []struct {
		name      string
		fieldList *ast.FieldList
		results   s2aModels.Results
	}{
		{
			name:      "nil field list",
			fieldList: nil,
			results:   s2aModels.Results{},
		},
		{
			name: "empty field list",
			fieldList: &ast.FieldList{
				List: []*ast.Field{},
			},
			results: s2aModels.Results{},
		},
		{
			name: "single unnamed field",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "res0",
						TypeExpr: "string",
					},
				},
			},
		},
		{
			name: "multiple unnamed fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "res0",
						TypeExpr: "string",
					},
				},
				{
					Field: s2aModels.Field{
						Name:     "res1",
						TypeExpr: "error",
					},
				},
			},
		},
		{
			name: "single named field",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("result")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "result",
						TypeExpr: "string",
					},
				},
			},
		},
		{
			name: "multiple named fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type:  ast.NewIdent("string"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("err")},
						Type:  ast.NewIdent("error"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "data",
						TypeExpr: "string",
					},
				},
				{
					Field: s2aModels.Field{
						Name:     "err",
						TypeExpr: "error",
					},
				},
			},
		},
		{
			name: "field with multiple names",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("result1"),
							ast.NewIdent("result2"),
						},
						Type: ast.NewIdent("int"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "result1",
						TypeExpr: "int",
					},
				},
				{
					Field: s2aModels.Field{
						Name:     "result2",
						TypeExpr: "int",
					},
				},
			},
		},
		{
			name: "mixed named and unnamed fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type:  ast.NewIdent("string"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			results: s2aModels.Results{
				{
					Field: s2aModels.Field{
						Name:     "data",
						TypeExpr: "string",
					},
				},
				{
					Field: s2aModels.Field{
						Name:     "res1",
						TypeExpr: "error",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			results, err := parser.fieldList2Results(tc.fieldList)

			require.NoError(t, err)
			assert.Equal(t, tc.results, results)
		})
	}
}

func TestSourceParser_field2Param(t *testing.T) {
	tt := []struct {
		name   string
		field  *ast.Field
		params s2aModels.Params
	}{
		{
			name:   "nil field",
			field:  nil,
			params: s2aModels.Params{},
		},
		{
			name: "unnamed field with string type",
			field: &ast.Field{
				Type: ast.NewIdent("string"),
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
		},
		{
			name: "named field with single name",
			field: &ast.Field{
				Names: []*ast.Ident{ast.NewIdent("result")},
				Type:  ast.NewIdent("string"),
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "result",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
		},
		{
			name: "named field with multiple names",
			field: &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("result1"),
					ast.NewIdent("result2"),
				},
				Type: ast.NewIdent("int"),
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "result1",
						TypeExpr: "int",
					},
					Type: s2aModels.Int,
				},
				{
					Field: s2aModels.Field{
						Name:     "result2",
						TypeExpr: "int",
					},
					Type: s2aModels.Int,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			results, err := parser.field2Params(tc.field)

			require.NoError(t, err)
			assert.Equal(t, tc.params, results)
		})
	}
}

func TestSourceParser_fieldList2Params(t *testing.T) {
	tt := []struct {
		name      string
		fieldList *ast.FieldList
		params    s2aModels.Params
	}{
		{
			name:      "nil field list",
			fieldList: nil,
			params:    s2aModels.Params{},
		},
		{
			name: "empty field list",
			fieldList: &ast.FieldList{
				List: []*ast.Field{},
			},
			params: s2aModels.Params{},
		},
		{
			name: "single unnamed field",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
		},
		{
			name: "multiple unnamed fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
				{
					Field: s2aModels.Field{
						Name:     "",
						TypeExpr: "error",
					},
					Type: s2aModels.Error,
				},
			},
		},
		{
			name: "single named field",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("result")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "result",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
		},
		{
			name: "multiple named fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type:  ast.NewIdent("string"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("err")},
						Type:  ast.NewIdent("error"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "data",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
				{
					Field: s2aModels.Field{
						Name:     "err",
						TypeExpr: "error",
					},
					Type: s2aModels.Error,
				},
			},
		},
		{
			name: "field with multiple names",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("result1"),
							ast.NewIdent("result2"),
						},
						Type: ast.NewIdent("int"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "result1",
						TypeExpr: "int",
					},
					Type: s2aModels.Int,
				},
				{
					Field: s2aModels.Field{
						Name:     "result2",
						TypeExpr: "int",
					},
					Type: s2aModels.Int,
				},
			},
		},
		{
			name: "mixed named and unnamed fields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type:  ast.NewIdent("string"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			params: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "data",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
				{
					Field: s2aModels.Field{
						Name:     "",
						TypeExpr: "error",
					},
					Type: s2aModels.Error,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			results, err := parser.fieldList2Params(tc.fieldList)

			require.NoError(t, err)
			assert.Equal(t, tc.params, results)
		})
	}
}

func TestSourceParser_parseFuncDecl(t *testing.T) {
	t.Run("nil decl", func(t *testing.T) {
		parser := NewSourceParser()

		err := parser.parseFuncDecl(nil)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		parser := NewSourceParser()

		err := parser.parseFuncDecl(&ast.FuncDecl{
			Name: ast.NewIdent("name"),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("param")},
							Type:  ast.NewIdent("string"),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("result")},
							Type:  ast.NewIdent("int"),
						},
					},
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(
			t,
			s2aModels.Functions{
				{
					Name: "name",
					Params: s2aModels.Params{
						{
							Field: s2aModels.Field{
								Name:     "param",
								TypeExpr: "string",
							},
							Type: s2aModels.String,
						},
					},
					Results: s2aModels.Results{
						{
							Field: s2aModels.Field{
								Name:     "result",
								TypeExpr: "int",
							},
						},
					},
				},
			},
			parser.Funcs,
		)
	})
}

func TestSourceParser_checkHTTPHandler(t *testing.T) {
	tt := []struct {
		name        string
		handlerName string
		decl        *ast.GenDecl
		exists      bool
	}{
		{
			name:        "nil declarataion",
			handlerName: "handler",
			decl:        nil,
			exists:      false,
		},
		{
			name:        "not VAR declaration",
			handlerName: "handler",
			decl:        &ast.GenDecl{Tok: token.IMPORT},
			exists:      false,
		},
		{
			name:        "no value specs",
			handlerName: "handler",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.TypeSpec{},
				},
			},
			exists: false,
		},
		{
			name:        "not handler declaration",
			handlerName: "handler",
			decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("unknown"),
						},
					},
				},
			},
			exists: false,
		},
		{
			name:        "not handler declaration",
			handlerName: "handler",
			decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("handler"),
						},
					},
				},
			},
			exists: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()
			parser.HTTPHandler.Name = tc.handlerName

			parser.checkHTTPHandler(tc.decl)
			assert.Equal(t, tc.exists, parser.HTTPHandler.Exists)
		})
	}
}

func TestSourceParser_parseImports(t *testing.T) {
	tt := []struct {
		name    string
		decl    *ast.GenDecl
		imports []s2aModels.Import
	}{
		{
			name:    "nil declarataion",
			decl:    nil,
			imports: make([]s2aModels.Import, 0),
		},
		{
			name:    "not IMPORT declaration",
			decl:    &ast.GenDecl{Tok: token.VAR},
			imports: make([]s2aModels.Import, 0),
		},
		{
			name: "no import specs",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.TypeSpec{},
				},
			},
			imports: make([]s2aModels.Import, 0),
		},
		{
			name: "couldn't unqiute path",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "path",
						},
					},
				},
			},
			imports: make([]s2aModels.Import, 0),
		},
		{
			name: "no alias",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "\"path\"",
						},
					},
				},
			},
			imports: []s2aModels.Import{
				{
					Path: "path",
				},
			},
		},
		{
			name: "import with alial",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "\"path\"",
						},
						Name: ast.NewIdent("alias"),
					},
				},
			},
			imports: []s2aModels.Import{
				{
					Path:  "path",
					Alias: "alias",
				},
			},
		},
		{
			name: "multiple imports",
			decl: &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "\"path1\"",
						},
						Name: ast.NewIdent("alias1"),
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "\"path2\"",
						},
					},
				},
			},
			imports: []s2aModels.Import{
				{
					Path:  "path1",
					Alias: "alias1",
				},
				{
					Path: "path2",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			parser.parseImports(tc.decl)
			assert.Equal(t, tc.imports, parser.Imports)
		})
	}
}

func TestSourceParser_Parse(t *testing.T) {
	tt := []struct {
		name          string
		pkgPath       string
		handlerName   string
		handlerExists bool
	}{
		{
			name:          "no API HTTP handler",
			pkgPath:       "../testdata/ok",
			handlerName:   "noHandler",
			handlerExists: false,
		},
		{
			name:          "has API HTTP handler",
			pkgPath:       "../testdata/ok",
			handlerName:   "handler",
			handlerExists: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSourceParser()

			err := parser.Parse(tc.pkgPath, tc.handlerName)

			require.NoError(t, err)
			assert.Equal(
				t,
				s2aModels.Functions{
					s2aModels.Function{
						Name: "method1",
						Params: s2aModels.Params{
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "pathString",
									TypeExpr: "string"},
								Type: "string", Origin: "",
							},
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "pathInt",
									TypeExpr: "int"},
								Type: "int", Origin: "",
							},
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "pathFloat64",
									TypeExpr: "float64"},
								Type: "float", Origin: "",
							},
						},
						Results: s2aModels.Results{
							s2aModels.Result{
								Field: s2aModels.Field{
									Name:     "result",
									TypeExpr: "string",
								},
							},
							s2aModels.Result{
								Field: s2aModels.Field{
									Name:     "err",
									TypeExpr: "error",
								},
							},
						},
						Endpoint: s2aModels.Endpoint{
							Method: "",
							Path:   "",
						},
					},
					s2aModels.Function{
						Name: "method7",
						Params: s2aModels.Params{
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "bodyModel",
									TypeExpr: "ThisPackageModel",
								},
								Type:   "custom[]ok.ThisPackageModel",
								Origin: "",
							},
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "bodyModelList",
									TypeExpr: "[]ThisPackageModel",
								},
								Type:   "slice[]custom[]ok.ThisPackageModel",
								Origin: "",
							},
							s2aModels.Param{
								Field: s2aModels.Field{
									Name:     "bodyModelMap",
									TypeExpr: "map[string]models.AnotherPackageModel",
								},
								Type:   "map[]custom[]models.AnotherPackageModel",
								Origin: "",
							},
						},
						Results: s2aModels.Results{
							s2aModels.Result{
								Field: s2aModels.Field{
									Name:     "code",
									TypeExpr: "int",
								},
							},
							s2aModels.Result{
								Field: s2aModels.Field{
									Name:     "err",
									TypeExpr: "error",
								},
							},
						},
						Endpoint: s2aModels.Endpoint{
							Method: "",
							Path:   "",
						},
					},
				},
				parser.Funcs,
			)
			assert.Equal(
				t,
				[]s2aModels.Import{
					{
						Path:  "net/http",
						Alias: "",
					},
					{
						Path:  "github.com/KianIt/swag2api/parser/testdata/ok/models",
						Alias: "models",
					},
				},
				parser.Imports,
			)
			assert.Equal(
				t,
				models.HTTPHandlerInfo{
					Name:   tc.handlerName,
					Exists: tc.handlerExists,
				},
				parser.HTTPHandler,
			)
		})
	}
}
