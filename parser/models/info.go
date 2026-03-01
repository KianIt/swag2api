package models

import s2aModels "github.com/KianIt/swag2api/models"

// ParsingInfo is complete parsing information.
type ParsingInfo struct {
	PkgName     string
	Imports     []s2aModels.Import
	Funcs       s2aModels.Functions
	HTTPHandler HTTPHandlerInfo
}

// HTTPHandlerInfo is HTTP handler parsing information.
type HTTPHandlerInfo struct {
	Name   string
	Exists bool
}
