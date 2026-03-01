package utils

import (
	"path"
	"strings"
)

// IsGoSource checks that the file is a valid Golang source code file.
func IsGoSource(name string) bool {
	return path.Ext(name) == ".go" && !strings.HasSuffix(name, "_test.go") && path.Base(name) != "generated.go"
}
