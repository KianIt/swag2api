package swag2api

import (
	"fmt"
	"log"
	"path"

	"github.com/KianIt/swag2api/builder"
	"github.com/KianIt/swag2api/parser"
	"github.com/KianIt/swag2api/validator"
)

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

	p := parser.NewParser()
	if err := p.Parse(pkgPath, mainFile, handlerName); err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	if err := builder.NewBuilder(p.GetInfo()).Build(apiPath); err != nil {
		return fmt.Errorf("building: %w", err)
	}

	log.Printf("API generated successfully")

	return nil
}
