package models

import (
	"fmt"
	"strings"
)

// Param is a param of a function.
type Param struct {
	Field
	Type   ParamType
	Origin ParamOrigin
}

// NameOrigin returns the param's name joined with its origin.
func (p Param) NameOrigin() string {
	return p.Name + string(p.Origin)
}

// Params is a list of Param.
type Params []Param

// Map returns a map of Param by name.
func (ps Params) Map() map[string]Param {
	pMap := make(map[string]Param)

	for _, p := range ps {
		pMap[p.Name] = p
	}

	return pMap
}

// ParamType is a type of Param.
type ParamType string

// ParamType checks if the param type equals to another type or is more specific.
func (pt ParamType) Is(paramType ParamType) bool {
	return strings.Contains(string(pt), string(paramType))
}

// MapOf returns a new map type with values of the current type.
func (pt ParamType) MapOf() ParamType {
	return Map + "[]" + pt
}

// MapOf returns a new slice type with values of the current type.
func (pt ParamType) SliceOf() ParamType {
	return Slice + "[]" + pt
}

// MapOf returns a new custom type that equals to the current type.
func (pt ParamType) CustomOf() ParamType {
	return Custom + "[]" + pt
}

const (
	Int    ParamType = "int"
	Float  ParamType = "float"
	String ParamType = "string"
	Bool   ParamType = "bool"
	Error  ParamType = "error"
	Any    ParamType = "any"
	Map    ParamType = "map"
	Slice  ParamType = "slice"
	Custom ParamType = "custom"
)

// ParamOrigin is an origin of Param.
type ParamOrigin string

const (
	Path  ParamOrigin = "Path"
	Query ParamOrigin = "Query"
	Body  ParamOrigin = "Body"
)

// GetParamOrigin parses a param origin from string.
func GetParamOrigin(originName string) (ParamOrigin, error) {
	switch originName {
	case "path":
		return Path, nil
	case "query":
		return Query, nil
	case "body":
		return Body, nil
	default:
		return "", fmt.Errorf("unsupported param origin: '%s'", originName)
	}
}
