package models

import s2aModels "github.com/KianIt/swag2api/models"

type ParsingInfo struct {
	PkgName     string
	Imports     []s2aModels.Import
	Funcs       s2aModels.Functions
	HTTPHandler HTTPHandlerInfo
}

type HTTPHandlerInfo struct {
	Name   string
	Exists bool
}
