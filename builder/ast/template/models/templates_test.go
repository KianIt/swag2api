package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistingTemplate(t *testing.T) {
	tt := []struct {
		name         string
		templateName string
		exists       bool
	}{
		{
			name:         "Template exists",
			templateName: ErrorResponse,
			exists:       true,
		},
		{
			name:         "Template doesn't exist",
			templateName: "unknown",
			exists:       false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exists, IsExistingTemplate(tc.templateName))
		})
	}
}
