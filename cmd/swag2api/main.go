package main

import (
	"flag"
	"log"

	"github.com/KianIt/swag2api"
)

const (
	defaultPkgPath     = "."
	defaultMainFile    = "main.go"
	defaultAPIFile     = "generated.go"
	defaultHandlerName = "s2aHandler"
)

var (
	pkgPath     *string
	mainFile    *string
	apiFile     *string
	handlerName *string
)

func initFlags() {
	pkgPath = flag.String("pkg", defaultPkgPath, "Path to the Golang package")
	mainFile = flag.String("main", defaultMainFile, "Name of the swag main file")
	apiFile = flag.String("to", defaultAPIFile, "Name of the generated API file")
	handlerName = flag.String("handler", defaultHandlerName, "Name of the API HTTP handler")
	flag.Parse()
}

func main() {
	initFlags()

	if err := swag2api.Generate(*pkgPath, *mainFile, *apiFile, *handlerName); err != nil {
		log.Fatalf("swag2api error: %v", err)
	}
}
