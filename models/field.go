package models

import (
	"fmt"
	"strconv"
	"unicode"
)

// Field is a named field of a specific type.
type Field struct {
	Name     string
	TypeExpr string
}

// IsError checks if the field is an error.
func (f Field) IsError() bool {
	return f.TypeExpr == "error"
}

// NameCapitalized returns the field's name with the first letter uppercased.
func (f Field) NameCapitalized() string {
	return string(unicode.ToTitle(rune(f.Name[0]))) + f.Name[1:]
}

// JSONTag returns the field's JSON tag.
func (f Field) JSONTag() string {
	return fmt.Sprintf("`json:%s`", strconv.Quote(f.Name))
}
