package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint_HandlerPath(t *testing.T) {
	endpoint := Endpoint{
		Method: "METHOD",
		Path:   "path",
	}
	path := endpoint.HandlerPath()

	assert.Equal(t, "METHOD path", path)
}
