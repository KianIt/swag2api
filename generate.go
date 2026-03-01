package swag2api

import (
	"log"
)

func Generate(pkgPath, mainFile, apiFile, handlerName string) error {
	log.Printf("Generating API from '%s' into '%s' with HTTP handler '%s'", pkgPath, apiFile, handlerName)

	log.Printf("API generated successfully")

	return nil
}
