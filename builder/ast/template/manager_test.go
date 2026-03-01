package template

import (
	"testing"

	"github.com/KianIt/swag2api/builder/ast/template/models"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	manager := NewManager()

	err := manager.Load()
	assert.NoError(t, err)

	templates, err := manager.GetTemplates()
	assert.NoError(t, err)

	assert.Equal(t, len(models.TemplateNames), len(templates))
}
