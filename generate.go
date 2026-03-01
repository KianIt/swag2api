package swag2api

import (
	"fmt"
	"log"
	"path"

	"github.com/KianIt/swag2api/builder"
	parserwrapper "github.com/KianIt/swag2api/parser/parser-wrapper"
	"github.com/KianIt/swag2api/validator"
)

// Generate is an entry point for the API generation.
//
// It parses information from swag annotations and source code
// and uses it to build the API.
//
// On the validation stage it checks that
// the source code directory exists,
// the swag main file exists, and
// the generated API file doesn't exists (deletes if existed).
//
// On the parsing stage it reads the function definitions
// from the source code and joins them with the corresponding
// swag annotations to obtain the complete information
// about the functions and the future API.
//
// On the building stage it uses the parsing information
// to generate HTTP handlers for all the annotated function,
// every function handler uses an HTTP request to properly call
// the annotated function and return its results.
// All the function handlers are available via a
// signle HTTP handler that dispatches requests to the API.
//
// Params:
//   - pkgPath: path to the source code directory;
//   - mainFile: name of the swag main file;
//   - apiFile: name of the generated API file;
//   - handlerName: name of the API HTTP handler.
func Generate(pkgPath, mainFile, apiFile, handlerName string) error {
	log.Printf("Generating API from '%s' into '%s' with HTTP handler '%s'", pkgPath, apiFile, handlerName)

	if err := validator.ValidatePkg(pkgPath); err != nil {
		return fmt.Errorf("pkg: %w", err)
	}

	mainPath := path.Join(pkgPath, mainFile)
	if err := validator.ValidateMainFile(mainPath); err != nil {
		return fmt.Errorf("main file: %w", err)
	}

	apiPath := path.Join(pkgPath, apiFile)
	if err := validator.ValidateAPIFile(apiPath); err != nil {
		return fmt.Errorf("API file: %w", err)
	}

	p := parserwrapper.NewParserWrapper()
	if err := p.Parse(pkgPath, mainFile, handlerName); err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	if err := builder.NewBuilder(p.GetInfo()).Build(apiPath); err != nil {
		return fmt.Errorf("building: %w", err)
	}

	log.Printf("API generated successfully")

	return nil
}
