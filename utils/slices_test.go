package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tt := []struct {
		name      string
		slice     []string
		value     string
		containts bool
	}{
		{
			name:      "Slice contains value",
			slice:     []string{"x"},
			value:     "x",
			containts: true,
		},
		{
			name:      "Slice doesn't contain value",
			slice:     []string{"x"},
			value:     "y",
			containts: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.containts, Contains(tc.slice, tc.value))
		})
	}
}
