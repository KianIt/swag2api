package models

import (
	"fmt"
	"strconv"
)

type Import struct {
	Path  string
	Alias string
}

func (i Import) Render() string {
	return fmt.Sprintf("%s %s", i.Alias, strconv.Quote(i.Path))
}
