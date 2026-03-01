package models

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
)

func TestParam_NameOrigin(t *testing.T) {
	param := Param{
		Field: Field{
			Name: "name",
		},
		Origin: "origin",
	}
	assert.Equal(t, "nameorigin", param.NameOrigin())
}

func TestParam_Map(t *testing.T) {
	params := Params{
		{Field: Field{Name: "name1"}},
		{Field: Field{Name: "name2"}},
	}
	fParam := map[string]Param{
		"name1": params[0],
		"name2": params[1],
	}
	assert.Equal(t, fParam, params.Map())
}

func TestParamType_Is(t *testing.T) {
	tt := []struct {
		name       string
		paramType1 ParamType
		paramType2 ParamType
		is         bool
	}{
		{
			name:       "Param types different",
			paramType1: ParamType("slice[]"),
			paramType2: ParamType("map[]"),
			is:         false,
		},
		{
			name:       "Param type 2 is more concrete",
			paramType1: ParamType("map[]"),
			paramType2: ParamType("map[]string"),
			is:         false,
		},
		{
			name:       "Param type 1 is more concrete",
			paramType1: ParamType("map[]string"),
			paramType2: ParamType("map[]"),
			is:         true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.is, tc.paramType1.Is(tc.paramType2))
		})
	}
}

func TestParamType_MapOf(t *testing.T) {
	assert.Equal(t, ParamType("map[]type"), ParamType("type").MapOf())
}

func TestParamType_SliceOf(t *testing.T) {
	assert.Equal(t, ParamType("slice[]type"), ParamType("type").SliceOf())
}

func TestParamType_CustomOf(t *testing.T) {
	assert.Equal(t, ParamType("custom[]type"), ParamType("type").CustomOf())
}

func TestGetParamOrigin(t *testing.T) {
	tt := []struct {
		name       string
		originName string
		origin     ParamOrigin
		isErr      bool
	}{
		{
			name:       "Path origin",
			originName: "path",
			origin:     Path,
		},
		{
			name:       "Query origin",
			originName: "query",
			origin:     Query,
		},
		{
			name:       "Body origin",
			originName: "body",
			origin:     Body,
		},
		{
			name:  "Unknown origin",
			isErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			origin, err := GetParamOrigin(tc.originName)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.origin, origin)
			}
		})
	}
}
