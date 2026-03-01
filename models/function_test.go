package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction_HandlerName(t *testing.T) {
	function := Function{
		Name: "name",
	}
	assert.Equal(t, "_handler_name", function.HandlerName())
}

func TestFunctions_Map(t *testing.T) {
	functions := Functions{
		{Name: "name1"},
		{Name: "name2"},
	}
	fMap := map[string]Function{
		"name1": functions[0],
		"name2": functions[1],
	}
	assert.Equal(t, fMap, functions.Map())
}
