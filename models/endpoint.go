package models

import "fmt"

// Endpoint is a description of a function HTTP endpoint.
type Endpoint struct {
	Method string
	Path   string
}

// HandlerPath returns the endpoint's method and path as a string.
func (e Endpoint) HandlerPath() string {
	return fmt.Sprintf("%s %s", e.Method, e.Path)
}
