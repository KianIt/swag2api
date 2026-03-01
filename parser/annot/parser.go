package annot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"

	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"

	s2aModels "github.com/KianIt/swag2api/models"
)

const (
	parseFlags = swag.ParseAll
	parseDepth = 100
)

// AnnotParser is a swag annotations parser.
//
// Uses the swag.Parser to parse the annotations.
type AnnotParser struct {
	swag  *swag.Parser
	Funcs s2aModels.Functions
}

// NewAnnotParser returns a new swag annotations parser.
func NewAnnotParser() *AnnotParser {
	return &AnnotParser{
		swag:  swag.New(swag.SetParseDependency(parseFlags)),
		Funcs: make(s2aModels.Functions, 0),
	}
}

// Parse runs parsing.
//
// Runs the swag.Parser and uses its results to obtain
// information about functions.
func (p *AnnotParser) Parse(pkgPath, mainFile string) error {
	log.Printf("Parsing swag annotations from '%s'", pkgPath)

	if err := p.swag.ParseAPI(pkgPath, mainFile, parseDepth); err != nil {
		return fmt.Errorf("swag: %w", err)
	}

	// Parsing every method of every path.
	for path, item := range p.swag.GetSwagger().Paths.Paths {
		if err := p.parsePath(path, item); err != nil {
			return fmt.Errorf("path '%s': %w", path, err)
		}
	}

	sort.Slice(p.Funcs, func(i, j int) bool { return p.Funcs[i].Name < p.Funcs[j].Name })

	log.Printf("Swag annotations parsed succesfully")

	return nil
}

// parsePath parses path information.
func (p *AnnotParser) parsePath(path string, item spec.PathItem) error {
	if item.Get != nil {
		if err := p.parseMethod(http.MethodGet, path, item.Get); err != nil {
			return fmt.Errorf("%s: %w", http.MethodGet, err)
		}
	}

	if item.Post != nil {
		if err := p.parseMethod(http.MethodPost, path, item.Post); err != nil {
			return fmt.Errorf("%s: %w", http.MethodPost, err)
		}
	}

	if item.Put != nil {
		if err := p.parseMethod(http.MethodPut, path, item.Put); err != nil {
			return fmt.Errorf("%s: %w", http.MethodPut, err)
		}
	}

	if item.Patch != nil {
		if err := p.parseMethod(http.MethodPatch, path, item.Patch); err != nil {
			return fmt.Errorf("%s: %w", http.MethodPatch, err)
		}
	}

	if item.Delete != nil {
		if err := p.parseMethod(http.MethodDelete, path, item.Delete); err != nil {
			return fmt.Errorf("%s: %w", http.MethodDelete, err)
		}
	}

	if item.Options != nil {
		if err := p.parseMethod(http.MethodOptions, path, item.Options); err != nil {
			return fmt.Errorf("%s: %w", http.MethodOptions, err)
		}
	}

	if item.Head != nil {
		if err := p.parseMethod(http.MethodHead, path, item.Head); err != nil {
			return fmt.Errorf("%s: %w", http.MethodHead, err)
		}
	}

	return nil
}

// parseMethod parses method information.
func (p *AnnotParser) parseMethod(method, path string, op *spec.Operation) error {
	if op == nil {
		return fmt.Errorf("op is nil")
	}

	if op.ID == "" {
		return fmt.Errorf("ID is empty")
	}

	log.Printf("Parsing function '%s' (%s %s)", op.ID, method, path)

	params, err := parameters2Params(op.Parameters)
	if err != nil {
		return fmt.Errorf("parameters: %w", err)
	}

	p.Funcs = append(p.Funcs, s2aModels.Function{
		Name:   op.ID,
		Params: params,
		Endpoint: s2aModels.Endpoint{
			Method: method,
			Path:   path,
		},
	})

	return nil
}

// parameters2Params converts a list of spec.Parameter into a s2aModels.Params.
func parameters2Params(parameters []spec.Parameter) (s2aModels.Params, error) {
	params := make(s2aModels.Params, 0, len(parameters))

	for _, parameter := range parameters {
		param, err := parameter2Param(parameter)
		if err != nil {
			return nil, fmt.Errorf("parameter '%s': %w", parameter.Name, err)
		}

		params = append(params, param)
	}

	return params, nil
}

// parameters2Params converts a spec.Parameter into a s2aModels.Param.
func parameter2Param(parameter spec.Parameter) (s2aModels.Param, error) {
	paramOrigin, err := s2aModels.GetParamOrigin(parameter.In)
	if err != nil {
		return s2aModels.Param{}, fmt.Errorf("origin: %w", err)
	}

	paramType, err := getParamType(parameter)
	if err != nil {
		return s2aModels.Param{}, fmt.Errorf("type: %w", err)
	}

	param := s2aModels.Param{
		Field: s2aModels.Field{
			Name: parameter.Name,
		},
		Type:   paramType,
		Origin: paramOrigin,
	}

	return param, nil
}

// getParamType parses a s2aModels.ParamType from a spec.Parameter.
//
// Uses data from the parameter schema if the parameter type is given.
// Otherwise constructs a custom schema from parameter data and uses it.
func getParamType(parameter spec.Parameter) (s2aModels.ParamType, error) {
	var schema *spec.Schema

	if parameter.Type == "" {
		schema = parameter.Schema
	} else {
		var (
			items                *spec.SchemaOrArray = nil
			additionalProperties *spec.SchemaOrBool  = nil
		)

		if parameter.Schema != nil {
			items = parameter.Schema.Items
			additionalProperties = parameter.Schema.AdditionalProperties
		}

		schema = &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 spec.StringOrArray([]string{parameter.Type}),
				Items:                items,
				AdditionalProperties: additionalProperties,
				Format:               parameter.Format,
			},
		}
	}

	return getSchemaParamType(schema)
}

// getSchemaParamType parses a s2aModels.ParamType from a *spec.Schema.
func getSchemaParamType(schema *spec.Schema) (s2aModels.ParamType, error) {
	if schema == nil {
		return "", errors.New("schema is nil")
	}

	// If the schema type is not given, then it's a custom or any type.
	if len(schema.Type) == 0 || schema.Type[0] == "" {
		// Getting the custom type name from the URL.
		if url := schema.Ref.GetURL(); url != nil {
			subType := s2aModels.ParamType(path.Base(url.String()))

			// Returning a custom type.
			return subType.CustomOf(), nil
		} else {
			return s2aModels.Any, nil
		}
	}

	// Parsing swag type.
	var paramType s2aModels.ParamType
	switch schema.Type[0] {
	case swag.STRING:
		paramType = s2aModels.String
	case swag.BOOLEAN:
		paramType = s2aModels.Bool
	case swag.INTEGER:
		paramType = s2aModels.Int
	case swag.NUMBER:
		paramType = s2aModels.Float
	case swag.ERROR:
		paramType = s2aModels.Error
	case swag.INTERFACE, swag.ANY:
		paramType = s2aModels.Any
	case swag.ARRAY:
		// Array is a slice.
		var (
			subType s2aModels.ParamType
			err     error
		)

		// If the items given, trying to get the item type.
		if schema.Items != nil {
			subType, err = getSchemaParamType(schema.Items.Schema)
			if err != nil {
				return "", err
			}
		}

		// Returning a slice.
		paramType = subType.SliceOf()
	case swag.OBJECT:
		// Object is a map.

		// Trying to get the value type.
		subType, err := getSchemaParamType(schema.AdditionalProperties.Schema)
		if err != nil {
			return "", err
		}

		// Returning a map.
		paramType = subType.MapOf()
	}

	if paramType != "" {
		return paramType, nil
	}

	return "", fmt.Errorf("unknown type: '%s'", schema.Type[0])
}
