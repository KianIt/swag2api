package utils

import (
	"os"
	"path"
	"strings"
)

// IsGoSource checks that the file is a valid Golang source code file.
func IsGoSource(info os.FileInfo) bool {
	name := info.Name()
	return !info.IsDir() && path.Ext(name) == ".go" && !strings.HasSuffix(name, "_test.go") && path.Base(name) != "generated.go"
}
