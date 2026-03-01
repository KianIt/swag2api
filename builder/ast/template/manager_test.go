package asttemplate

import (
	"testing"

	"github.com/KianIt/swag2api/builder/ast/template/models"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	manager := NewManager()

	err := manager.Load()
	require.NoError(t, err)

	templates, err := manager.GetTemplates()
	require.NoError(t, err)

	assert.Len(t, templates, len(models.TemplateNames))
}
