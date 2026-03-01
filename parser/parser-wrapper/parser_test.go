package parserwrapper

import (
	"testing"

	s2aModels "github.com/KianIt/swag2api/models"
	"github.com/KianIt/swag2api/parser/models"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
)

func TestCombineParams(t *testing.T) {
	tt := []struct {
		name           string
		annotParams    s2aModels.Params
		sourceParams   s2aModels.Params
		combinedParams s2aModels.Params
		isErr          bool
	}{
		{
			name:           "no params",
			annotParams:    s2aModels.Params{},
			sourceParams:   s2aModels.Params{},
			combinedParams: s2aModels.Params{},
		},
		{
			name: "param found in annotaions and not found in source",
			annotParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "annotParam",
						TypeExpr: "string",
					},
					Type:   s2aModels.String,
					Origin: s2aModels.Path,
				},
			},
			sourceParams:   s2aModels.Params{},
			combinedParams: s2aModels.Params{},
			isErr:          false,
		},
		{
			name:        "param found in source and not found in annotaions",
			annotParams: s2aModels.Params{},
			sourceParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "sourceParam",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
			combinedParams: s2aModels.Params{},
			isErr:          true,
		},
		{
			name: "source param type is not annotations param type",
			annotParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramName",
						TypeExpr: "string",
					},
					Type:   s2aModels.ParamType("").MapOf(),
					Origin: s2aModels.Path,
				},
			},
			sourceParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramName",
						TypeExpr: "string",
					},
					Type: s2aModels.String,
				},
			},
			combinedParams: s2aModels.Params{},
			isErr:          true,
		},
		{
			name: "ok",
			annotParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramName",
						TypeExpr: "string",
					},
					Type:   s2aModels.ParamType("").MapOf(),
					Origin: s2aModels.Path,
				},
			},
			sourceParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramName",
						TypeExpr: "string",
					},
					Type: s2aModels.String.MapOf(),
				},
			},
			combinedParams: s2aModels.Params{
				{
					Field: s2aModels.Field{
						Name:     "paramName",
						TypeExpr: "string",
					},
					Type:   s2aModels.String.MapOf(),
					Origin: s2aModels.Path,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			combinedParams, err := combineParams("func", tc.annotParams, tc.sourceParams)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.combinedParams, combinedParams)
			}
		})
	}
}

func TestParser_CombineFunctions(t *testing.T) {
	tt := []struct {
		name          string
		annotFuncs    s2aModels.Functions
		sourceFuncs   s2aModels.Functions
		combinedFuncs s2aModels.Functions
		isErr         bool
	}{
		{
			name:          "no functions",
			annotFuncs:    s2aModels.Functions{},
			sourceFuncs:   s2aModels.Functions{},
			combinedFuncs: s2aModels.Functions{},
		},
		{
			name: "func found in annotaions and not found in source",
			annotFuncs: s2aModels.Functions{
				{
					Name: "annotFunc",
				},
			},
			sourceFuncs:   s2aModels.Functions{},
			combinedFuncs: s2aModels.Functions{},
			isErr:         false,
		},
		{
			name:       "param found in source and not found in annotaions",
			annotFuncs: s2aModels.Functions{},
			sourceFuncs: s2aModels.Functions{
				{
					Name: "sourceFunc",
				},
			},
			combinedFuncs: s2aModels.Functions{},
			isErr:         false,
		},
		{
			name: "ok",
			annotFuncs: s2aModels.Functions{
				{
					Name: "funcName",
					Params: s2aModels.Params{
						{
							Field: s2aModels.Field{
								Name:     "paramName",
								TypeExpr: "string",
							},
							Type:   s2aModels.ParamType("").MapOf(),
							Origin: s2aModels.Path,
						},
					},
					Endpoint: s2aModels.Endpoint{
						Method: "method",
						Path:   "path",
					},
				},
			},
			sourceFuncs: s2aModels.Functions{
				{
					Name: "funcName",
					Params: s2aModels.Params{
						{
							Field: s2aModels.Field{
								Name:     "paramName",
								TypeExpr: "string",
							},
							Type: s2aModels.String.MapOf(),
						},
					},
					Results: s2aModels.Results{
						{
							Field: s2aModels.Field{
								Name:     "resultName",
								TypeExpr: "int",
							},
						},
					},
				},
			},
			combinedFuncs: s2aModels.Functions{
				{
					Name: "funcName",
					Params: s2aModels.Params{
						{
							Field: s2aModels.Field{
								Name:     "paramName",
								TypeExpr: "string",
							},
							Type:   s2aModels.String.MapOf(),
							Origin: s2aModels.Path,
						},
					},
					Results: s2aModels.Results{
						{
							Field: s2aModels.Field{
								Name:     "resultName",
								TypeExpr: "int",
							},
						},
					},
					Endpoint: s2aModels.Endpoint{
						Method: "method",
						Path:   "path",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParserWrapper()
			parser.annot.Funcs = tc.annotFuncs
			parser.source.Funcs = tc.sourceFuncs

			combinedFuncs, err := parser.combineFuncs()

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.combinedFuncs, combinedFuncs)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	pkgPath := "../testdata/ok"
	mainFile := "main.go"
	handlerName := "handler"

	parser := NewParserWrapper()

	require.NoError(t, parser.Parse(pkgPath, mainFile, handlerName))

	assert.Equal(
		t,
		models.ParsingInfo{
			PkgName: "ok",
			Imports: []s2aModels.Import{
				{
					Path:  "net/http",
					Alias: "",
				},
				{
					Path:  "github.com/KianIt/swag2api/parser/testdata/ok/models",
					Alias: "models",
				},
			},
			Funcs: s2aModels.Functions{
				s2aModels.Function{
					Name: "method1",
					Params: s2aModels.Params{
						s2aModels.Param{
							Field: s2aModels.Field{
								Name:     "pathString",
								TypeExpr: "string",
							},
							Type:   "string",
							Origin: "Path",
						},
						s2aModels.Param{
							Field: s2aModels.Field{
								Name:     "pathInt",
								TypeExpr: "int",
							},
							Type:   "int",
							Origin: "Query",
						},
						s2aModels.Param{
							Field: s2aModels.Field{
								Name:     "pathFloat64",
								TypeExpr: "float64",
							},
							Type:   "float",
							Origin: "Body",
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
						Method: "GET",
						Path:   "/path-to-method1",
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
							Origin: "Body",
						},
						s2aModels.Param{
							Field: s2aModels.Field{
								Name:     "bodyModelList",
								TypeExpr: "[]ThisPackageModel",
							},
							Type:   "slice[]custom[]ok.ThisPackageModel",
							Origin: "Body",
						},
						s2aModels.Param{
							Field: s2aModels.Field{
								Name:     "bodyModelMap",
								TypeExpr: "map[string]models.AnotherPackageModel",
							},
							Type:   "map[]custom[]models.AnotherPackageModel",
							Origin: "Body",
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
						Method: "PATCH",
						Path:   "/path-to-method7",
					},
				}},
			HTTPHandler: models.HTTPHandlerInfo{Name: "handler", Exists: true},
		},
		parser.GetInfo(),
	)
}
