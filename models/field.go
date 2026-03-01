package models

import (
	"fmt"
	"strconv"
	"unicode"
)

type Field struct {
	Name     string
	TypeExpr string
}

func (f Field) IsError() bool {
	return f.TypeExpr == "error"
}

func (f Field) NameCapitalized() string {
	return string(unicode.ToTitle(rune(f.Name[0]))) + f.Name[1:]
}

func (f Field) JSONTag() string {
	return fmt.Sprintf("`json:%s`", strconv.Quote(f.Name))
}
