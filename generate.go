package swag2api

import (
	"fmt"
	"log"
	"path"

	"github.com/KianIt/swag2api/builder"
	"github.com/KianIt/swag2api/parser"
)

func Generate(pkgPath, mainFile, apiFile, handlerName string) error {
	log.Printf("Generating API from '%s' into '%s' with HTTP handler '%s'", pkgPath, apiFile, handlerName)

	apiPath := path.Join(pkgPath, apiFile)

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
