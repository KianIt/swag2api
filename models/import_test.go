package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport_Render(t *testing.T) {
	imp := Import{
		Path:  "path",
		Alias: "alias",
	}
	assert.Equal(t, "alias \"path\"", imp.Render())
}
