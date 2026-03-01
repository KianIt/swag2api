package annot

import (
	"net/http"
	"testing"

	"github.com/KianIt/swag2api/models"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/swag"
)

func TestGetSchemaParamType(t *testing.T) {
	tt := []struct {
		name      string
		schema    *spec.Schema
		paramType models.ParamType
		isErr     bool
	}{
		{
			name:   "schema is nil",
			schema: nil,
			isErr:  true,
		},
		{
			name: "any type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: nil,
				},
			},
			paramType: models.Any,
		},
		{
			name: "string type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.STRING,
					},
				},
			},
			paramType: models.String,
		},
		{
			name: "interface type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.INTERFACE,
					},
				},
			},
			paramType: models.Any,
		},
		{
			name: "slice type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.ARRAY,
					},
				},
			},
			paramType: models.ParamType("slice[]"),
		},
		{
			name: "slice of integer type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.ARRAY,
					},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{
									swag.INTEGER,
								},
							},
						},
					},
				},
			},
			paramType: models.ParamType("slice[]int"),
		},
		{
			name: "map type",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.OBJECT,
					},
				},
			},
			paramType: models.ParamType("map[]"),
		},
		{
			name: "map of bool",
			schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{
						swag.OBJECT,
					},
					AdditionalProperties: &spec.SchemaOrBool{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{
									swag.BOOLEAN,
								},
							},
						},
					},
				},
			},
			paramType: models.ParamType("map[]bool"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paramType, err := getSchemaParamType(tc.schema)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.paramType, paramType)
			}
		})
	}
}

func TestGetParamType(t *testing.T) {
	tt := []struct {
		name      string
		parameter spec.Parameter
		paramType models.ParamType
	}{
		{
			name: "parameter with param props",
			parameter: spec.Parameter{
				SimpleSchema: spec.SimpleSchema{
					Type: "",
				},
				ParamProps: spec.ParamProps{
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{
								swag.STRING,
							},
						},
					},
				},
			},
			paramType: models.String,
		},
		{
			name: "parameter with simple schema",
			parameter: spec.Parameter{
				SimpleSchema: spec.SimpleSchema{
					Type: swag.STRING,
				},
			},
			paramType: models.String,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paramType, err := getParamType(tc.parameter)
			require.NoError(t, err)
			assert.Equal(t, tc.paramType, paramType)
		})
	}
}

func TestParameter2Param(t *testing.T) {
	tt := []struct {
		name      string
		parameter spec.Parameter
		param     models.Param
		isErr     bool
	}{
		{
			name: "unknown param origin",
			parameter: spec.Parameter{
				ParamProps: spec.ParamProps{
					In: "unknown",
				},
			},
			isErr: true,
		},
		{
			name: "unknown param type",
			parameter: spec.Parameter{
				ParamProps: spec.ParamProps{
					In: "path",
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{
								"unknown",
							},
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "ok",
			parameter: spec.Parameter{
				ParamProps: spec.ParamProps{
					Name: "name",
					In:   "path",
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{
								swag.STRING,
							},
						},
					},
				},
			},
			param: models.Param{
				Field: models.Field{
					Name: "name",
				},
				Type:   models.String,
				Origin: models.Path,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			param, err := parameter2Param(tc.parameter)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.param, param)
			}
		})
	}
}

func TestParameters2Param(t *testing.T) {
	parameters := []spec.Parameter{
		{
			ParamProps: spec.ParamProps{
				Name: "name",
				In:   "path",
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{
							swag.STRING,
						},
					},
				},
			},
		},
	}
	params, err := parameters2Params(parameters)
	require.NoError(t, err)

	assert.Equal(
		t,
		models.Params{
			{
				Field: models.Field{
					Name: "name",
				},
				Type:   models.String,
				Origin: models.Path,
			},
		},
		params,
	)
}

func TestAnnotParser_parseMethod(t *testing.T) {
	parser := NewAnnotParser()

	t.Run("operation is nil", func(t *testing.T) {
		err := parser.parseMethod("method", "path", nil)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		err := parser.parseMethod("method", "path", &spec.Operation{
			OperationProps: spec.OperationProps{
				ID: "name",
				Parameters: []spec.Parameter{{
					ParamProps: spec.ParamProps{
						Name: "paramName",
						In:   "path",
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{
									swag.STRING,
								},
							},
						},
					},
				}},
			},
		})
		require.NoError(t, err)
		assert.Equal(
			t,
			models.Functions{
				{
					Name: "name",
					Params: models.Params{
						{
							Field: models.Field{
								Name: "paramName",
							},
							Type:   models.String,
							Origin: models.Path,
						},
					},
					Endpoint: models.Endpoint{
						Method: "method",
						Path:   "path",
					},
				},
			},
			parser.Funcs,
		)
	})
}

func TestAnnotParser_parsePath(t *testing.T) {
	parser := NewAnnotParser()

	t.Run("operation is nil", func(t *testing.T) {
		err := parser.parseMethod("method", "path", nil)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		err := parser.parsePath("path", spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Get: &spec.Operation{
					OperationProps: spec.OperationProps{
						ID: "name",
						Parameters: []spec.Parameter{{
							ParamProps: spec.ParamProps{
								Name: "paramName",
								In:   "path",
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{
											swag.STRING,
										},
									},
								},
							},
						}},
					},
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(
			t,
			models.Functions{
				{
					Name: "name",
					Params: models.Params{
						{
							Field: models.Field{
								Name: "paramName",
							},
							Type:   models.String,
							Origin: models.Path,
						},
					},
					Endpoint: models.Endpoint{
						Method: http.MethodGet,
						Path:   "path",
					},
				},
			},
			parser.Funcs,
		)
	})
}

func TestAnnotParser_Parse(t *testing.T) {
	tt := []struct {
		name     string
		pkgPath  string
		mainFile string
		isErr    bool
	}{
		{
			name:     "main file not found",
			pkgPath:  "../testdata/mainfilenotfound",
			mainFile: "main.go",
			isErr:    true,
		},
		{
			name:     "ok",
			pkgPath:  "../testdata/ok",
			mainFile: "main.go",
			isErr:    false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewAnnotParser()

			err := parser.Parse(tc.pkgPath, tc.mainFile)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(
					t,
					models.Functions{
						{
							Name: "method1",
							Params: models.Params{
								{
									Field: models.Field{
										Name: "pathString",
									},
									Type:   models.String,
									Origin: models.Path,
								},
								{
									Field: models.Field{
										Name: "pathInt",
									},
									Type:   models.Int,
									Origin: models.Query,
								},
								{
									Field: models.Field{
										Name: "pathFloat64",
									},
									Type:   models.Float,
									Origin: models.Body,
								},
							},
							Endpoint: models.Endpoint{
								Method: http.MethodGet,
								Path:   "/path-to-method1",
							},
						},
						models.Function{
							Name: "method7",
							Params: models.Params{
								models.Param{
									Field: models.Field{
										Name:     "bodyModel",
										TypeExpr: "",
									},
									Type:   "custom[]ok.ThisPackageModel",
									Origin: "Body"},
								models.Param{
									Field: models.Field{
										Name:     "bodyModelList",
										TypeExpr: "",
									},
									Type:   "slice[]custom[]ok.ThisPackageModel",
									Origin: "Body",
								},
								models.Param{
									Field: models.Field{
										Name:     "bodyModelMap",
										TypeExpr: "",
									},
									Type:   "map[]custom[]models.AnotherPackageModel",
									Origin: "Body",
								},
							},
							Results: models.Results(nil),
							Endpoint: models.Endpoint{
								Method: "PATCH",
								Path:   "/path-to-method7",
							},
						},
					},
					parser.Funcs,
				)
			}
		})
	}
}
