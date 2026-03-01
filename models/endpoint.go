package models

import "fmt"

type Endpoint struct {
	Method string
	Path   string
}

func (e Endpoint) HandlerPath() string {
	return fmt.Sprintf("%s %s", e.Method, e.Path)
}
