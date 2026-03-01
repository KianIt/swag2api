package models

import (
	"fmt"
	"strconv"
)

// Import is a description of a single Golang import.
type Import struct {
	Path  string
	Alias string
}

// Render returns the import's source code form as a string.
func (i Import) Render() string {
	return fmt.Sprintf("%s %s", i.Alias, strconv.Quote(i.Path))
}
