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

type AnnotParser struct {
	swag  *swag.Parser
	Funcs s2aModels.Functions
}

func NewAnnotParser() *AnnotParser {
	return &AnnotParser{
		swag:  swag.New(swag.SetParseDependency(parseFlags)),
		Funcs: make(s2aModels.Functions, 0),
	}
}

func (p *AnnotParser) Parse(pkgPath, mainFile string) error {
	log.Printf("Parsing swag annotations from '%s'", pkgPath)

	if err := p.swag.ParseAPI(pkgPath, mainFile, parseDepth); err != nil {
		return fmt.Errorf("swag: %w", err)
	}

	for path, item := range p.swag.GetSwagger().Paths.Paths {
		if err := p.parsePath(path, item); err != nil {
			return fmt.Errorf("path '%s': %w", path, err)
		}
	}

	sort.Slice(p.Funcs, func(i, j int) bool { return p.Funcs[i].Name < p.Funcs[j].Name })

	log.Printf("Swag annotations parsed succesfully")

	return nil
}

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

func getSchemaParamType(schema *spec.Schema) (s2aModels.ParamType, error) {
	if schema == nil {
		return "", errors.New("schema is nil")
	}

	if len(schema.Type) == 0 || schema.Type[0] == "" {
		if url := schema.Ref.GetURL(); url != nil {
			subType := s2aModels.ParamType(path.Base(url.String()))
			return subType.CustomOf(), nil
		} else {
			return s2aModels.Any, nil
		}
	}

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
		var (
			subType s2aModels.ParamType
			err     error
		)

		if schema.Items != nil {
			subType, err = getSchemaParamType(schema.Items.Schema)
			if err != nil {
				return "", err
			}
		}

		paramType = subType.SliceOf()
	case swag.OBJECT:
		subType, err := getSchemaParamType(schema.AdditionalProperties.Schema)
		if err != nil {
			return "", err
		}

		paramType = subType.MapOf()
	}

	if paramType != "" {
		return paramType, nil
	}

	return "", fmt.Errorf("unknown type: '%s'", schema.Type[0])
}
