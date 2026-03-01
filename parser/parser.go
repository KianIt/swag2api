package parser

import (
	"fmt"
	"log"

	s2aModels "github.com/KianIt/swag2api/models"
	"github.com/KianIt/swag2api/parser/annot"
	"github.com/KianIt/swag2api/parser/models"
	"github.com/KianIt/swag2api/parser/source"
)

// Parser is a wrapper for the swag annotations and source code parsers.
//
// It uses the both parsers to read the complete information about
// the annotated functions and the future API.
type Parser struct {
	annot  *annot.AnnotParser
	source *source.SourceParser
	info   models.ParsingInfo
}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{
		annot:  annot.NewAnnotParser(),
		source: source.NewSourceParser(),
	}
}

// Parse runs parsing.
//
// Runs the swag annotations and source code parsers and
// joins their results into the complete parsing information.
func (p *Parser) Parse(pkgPath, mainFname, handlerName string) error {
	log.Printf("Parsing started")

	if err := p.annot.Parse(pkgPath, mainFname); err != nil {
		return fmt.Errorf("annot: %w", err)
	}

	if err := p.source.Parse(pkgPath, handlerName); err != nil {
		return fmt.Errorf("source: %w", err)
	}

	combinedFuncs, err := p.combineFuncs()
	if err != nil {
		return fmt.Errorf("combining functions: %w", err)
	}

	p.info.PkgName = p.source.PkgName
	p.info.Imports = p.source.Imports
	p.info.Funcs = combinedFuncs
	p.info.HTTPHandler = p.source.HTTPHandler

	log.Printf("Parsing finished successfully")

	return nil
}

// GetInfo returns the complete parsing information.
func (p *Parser) GetInfo() models.ParsingInfo {
	return p.info
}

// combineFuncs combines functions from the swag annotations and source code parsers.
func (p *Parser) combineFuncs() (combinedFuncs s2aModels.Functions, err error) {
	combinedFuncs = make(s2aModels.Functions, 0, len(p.source.Funcs))

	annotFuncMap := p.annot.Funcs.Map()
	sourceFuncMap := p.source.Funcs.Map()

	for _, annotFunc := range p.annot.Funcs {
		if _, ok := sourceFuncMap[annotFunc.Name]; !ok {
			log.Printf("Warning: function '%s' found in annotations and not found in source code", annotFunc.Name)
		}
	}

	for _, sourceFunc := range p.source.Funcs {
		if annotFunc, ok := annotFuncMap[sourceFunc.Name]; ok {
			sourceFunc.Endpoint = annotFunc.Endpoint

			sourceParams, err := combineParams(sourceFunc.Name, annotFunc.Params, sourceFunc.Params)
			if err != nil {
				return nil, fmt.Errorf("function '%s' params: %w", sourceFunc.Name, err)
			}
			sourceFunc.Params = sourceParams

			combinedFuncs = append(combinedFuncs, sourceFunc)
		}
	}

	return
}

// combineFuncs combines function params from the swag annotations and source code parsers.
func combineParams(funcName string, annotParams, sourceParams s2aModels.Params) (combinedParams s2aModels.Params, err error) {
	combinedParams = make(s2aModels.Params, 0, len(sourceParams))

	annotParamMap := annotParams.Map()
	sourceParamMap := sourceParams.Map()

	for _, annotParam := range annotParams {
		if _, ok := sourceParamMap[annotParam.Name]; !ok {
			log.Printf("Warning: function '%s' param '%s' found in annotations and not found in source code", funcName, annotParam.Name)
		}
	}

	for _, sourceParam := range sourceParams {
		if annotParam, ok := annotParamMap[sourceParam.Name]; !ok {
			return nil, fmt.Errorf("param '%s' found in source code and not found in annotations", sourceParam.Name)
		} else if !sourceParam.Type.Is(annotParam.Type) {
			return nil, fmt.Errorf("param '%s' source code type '%s' isn`t equal to annotation type '%s'", sourceParam.Name, sourceParam.Type, annotParam.Type)
		} else {
			sourceParam.Origin = annotParam.Origin
			combinedParams = append(combinedParams, sourceParam)
		}
	}

	return
}
