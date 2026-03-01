package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoSource(t *testing.T) {
	tt := []struct {
		name       string
		fileName   string
		isGoSource bool
	}{
		{
			name:       "Go source file",
			fileName:   "abc.go",
			isGoSource: true,
		},
		{
			name:       "Not Go file",
			fileName:   "abc",
			isGoSource: false,
		},
		{
			name:       "Go test file",
			fileName:   "abc_test.go",
			isGoSource: false,
		},
		{
			name:       "Go generated file",
			fileName:   "generated.go",
			isGoSource: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isGoSource, IsGoSource(tc.fileName))
		})
	}
}
