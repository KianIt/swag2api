package swag2api

import (
	"fmt"
	"log"

	"github.com/KianIt/swag2api/parser"
)

func Generate(pkgPath, mainFile, apiFile, handlerName string) error {
	log.Printf("Generating API from '%s' into '%s' with HTTP handler '%s'", pkgPath, apiFile, handlerName)

	p := parser.NewParser()
	if err := p.Parse(pkgPath, mainFile, handlerName); err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	log.Printf("API generated successfully")

	return nil
}
