package models

import (
	"fmt"
	"strings"
)

type Param struct {
	Field
	Type   ParamType
	Origin ParamOrigin
}

func (p Param) NameOrigin() string {
	return p.Name + string(p.Origin)
}

type ParamType string

func (pt ParamType) Is(paramType ParamType) bool {
	return strings.Contains(string(pt), string(paramType))
}

func (pt ParamType) MapOf() ParamType {
	return Map + "[]" + pt
}

func (pt ParamType) SliceOf() ParamType {
	return Slice + "[]" + pt
}

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

type ParamOrigin string

const (
	Path  ParamOrigin = "Path"
	Query ParamOrigin = "Query"
	Body  ParamOrigin = "Body"
)

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

type Params []Param

func (ps Params) Map() map[string]Param {
	pMap := make(map[string]Param)

	for _, p := range ps {
		pMap[p.Name] = p
	}

	return pMap
}
