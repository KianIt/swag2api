package models

import "fmt"

// Function is the comlete infomation about a function and a corresponding HTTP endpoint.
type Function struct {
	Name     string
	Params   Params
	Results  Results
	Endpoint Endpoint
}

// HandlerName returns the function's endpoint handler name.
func (f Function) HandlerName() string {
	return fmt.Sprintf("_handler_%s", f.Name)
}

// Functions is a list of Function.
type Functions []Function

// Map returns a map of Function by their names.
func (fs Functions) Map() map[string]Function {
	fMap := make(map[string]Function)

	for _, f := range fs {
		fMap[f.Name] = f
	}

	return fMap
}
