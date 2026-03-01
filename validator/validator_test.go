package validator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePkg(t *testing.T) {
	tt := []struct {
		name    string
		pkgPath string
		err     error
	}{
		{
			name:    "Not a dir",
			pkgPath: "./testdata/pkg/notadir",
			err:     errNotDirectory,
		},
		{
			name:    "Not found",
			pkgPath: "./testdata/pkg/notfound",
			err:     errPathNotExist,
		},
		{
			name:    "Ok",
			pkgPath: "./testdata/pkg/ok",
			err:     nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePkg(tc.pkgPath)

			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestValidateMainFile(t *testing.T) {
	tt := []struct {
		name     string
		mainPath string
		err      error
	}{
		{
			name:     "Not a file",
			mainPath: "./testdata/main/notafile",
			err:      errNotFile,
		},
		{
			name:     "Not found",
			mainPath: "./testdata/main/notfound",
			err:      errPathNotExist,
		},
		{
			name:     "Ok",
			mainPath: "./testdata/main/ok",
			err:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMainFile(tc.mainPath)

			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestValidateAPIFile(t *testing.T) {
	// Creating the file for test.
	_, _ = os.Create("./testdata/api/ok")

	tt := []struct {
		name    string
		apiPath string
		err     error
	}{
		{
			name:    "Not a file",
			apiPath: "./testdata/api/notafile",
			err:     errNotFile,
		},
		{
			name:    "Not found",
			apiPath: "./testdata/api/notfound",
			err:     nil,
		},
		{
			name:    "Ok",
			apiPath: "./testdata/api/ok",
			err:     nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAPIFile(tc.apiPath)

			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}
