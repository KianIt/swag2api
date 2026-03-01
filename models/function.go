package models

import "fmt"

type Function struct {
	Name     string
	Params   Params
	Results  Results
	Endpoint Endpoint
}

func (f Function) HandlerName() string {
	return fmt.Sprintf("_handler_%s", f.Name)
}

type Functions []Function

func (fs Functions) Map() map[string]Function {
	fMap := make(map[string]Function)

	for _, f := range fs {
		fMap[f.Name] = f
	}

	return fMap
}
