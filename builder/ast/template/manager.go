package template

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"

	astModels "github.com/KianIt/swag2api/builder/ast/models"
	"github.com/KianIt/swag2api/builder/ast/template/models"
	"github.com/KianIt/swag2api/utils"
)

const templatesRelativePath = "./templates"

// Manager is a manager for source code templates.
//
// Loads, keeps and returns templates.
type Manager struct {
	fs      *token.FileSet
	declMap map[string]ast.Decl
}

// NewManager returns a new template manager.
func NewManager() *Manager {
	return &Manager{
		fs:      token.NewFileSet(),
		declMap: make(map[string]ast.Decl),
	}
}

// Load loads templates from the source code.
//
// Determines the template source code actual absoule path,
// parses the source code files and saves the templates.
func (m *Manager) Load() error {
	log.Printf("Loading templates")

	absPath, err := getTemplatesAbsPath()
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	log.Printf("Templates absolute path: '%s'", absPath)

	files, err := m.parseFiles(absPath)
	if err != nil {
		return fmt.Errorf("parsing files: %w", err)
	}

	for _, file := range files {
		ast.Walk(m, file)
	}

	log.Printf("Templates loaded successfully")

	return nil
}

// parseFiles parses files from the template source code.
func (m *Manager) parseFiles(dirPath string) ([]*ast.File, error) {
	files := make([]*ast.File, 0)

	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !utils.IsGoSource(info) {
			return nil
		}

		file, parseErr := parser.ParseFile(m.fs, path, nil, parser.ParseComments)
		if parseErr != nil {
			return fmt.Errorf("file '%s': %w", path, parseErr)
		}

		files = append(files, file)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// Visit implements the ast.Visitor interface.
//
// Allows the Manager to walk over nodes in a parsed file.
func (m *Manager) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch decl := node.(type) {
	case *ast.File:
		return m
	case *ast.GenDecl:
		// Parsing general declarations.
		for _, spec := range decl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok && models.IsExistingTemplate(typeSpec.Name.Name) {
				m.declMap[typeSpec.Name.Name] = decl
			}
		}

		return m
	case *ast.FuncDecl:
		// Parsing function declarations.
		if models.IsExistingTemplate(decl.Name.Name) {
			m.declMap[decl.Name.Name] = decl
		}

		return m
	}

	return nil
}

// GetTemplates returns loaded templates.
func (m *Manager) GetTemplates() ([]astModels.Decl, error) {
	if m.declMap == nil {
		return nil, errors.New("decl map is nil")
	}

	decls := make([]astModels.Decl, 0, len(models.TemplateNames))
	for _, name := range models.TemplateNames {
		if template, ok := m.declMap[name]; ok {
			decls = append(decls, template)
		}
	}

	return decls, nil
}

// getTemplatesAbsPath returns the template source code actual absolute path.
func getTemplatesAbsPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("can't get caller")
	}

	return filepath.Join(filepath.Dir(filename), templatesRelativePath), nil
}
