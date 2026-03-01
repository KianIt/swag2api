package templates

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalString(t *testing.T) {
	t.Run("Single value", func(t *testing.T) {
		value, err := _unmarshalString[int]("123")
		assert.NoError(t, err)
		assert.Equal(t, 123, value)
	})

	t.Run("Multiple int values", func(t *testing.T) {
		values, err := _unmarshalString[[]int]("[123, 123]")
		assert.NoError(t, err)
		assert.Equal(t, []int{123, 123}, values)
	})

	t.Run("Multiple string values", func(t *testing.T) {
		values, err := _unmarshalString[[]string]("[\"123\", \"123\"]")
		assert.NoError(t, err)
		assert.Equal(t, []string{"123", "123"}, values)
	})
}

func TestUnmarshalBytes(t *testing.T) {
	var response = _errorResponse{Error: "error"}
	data, _ := json.Marshal(response)

	value, err := _unmarshalBytes[_errorResponse](data)
	assert.NoError(t, err)
	assert.Equal(t, response, value)
}
