package models

import "github.com/KianIt/swag2api/utils"

const (
	BaseResponse     = "_baseResponse"
	UnmarshalString  = "_unmarshalString"
	UnmarshalBytes   = "_unmarshalBytes"
	HandleBadRequest = "_handleBadRequest"
	HandleResult     = "_handleResult"
	WriteResponse    = "_writeResponse"
)

// TemplateNames is a list of all templates.
var TemplateNames = []string{
	BaseResponse,
	UnmarshalString,
	UnmarshalBytes,
	HandleBadRequest,
	HandleResult,
	WriteResponse,
}

// IsExistingTemplate checks if a template with a given name exists.
func IsExistingTemplate(name string) bool {
	return utils.Contains(TemplateNames, name)
}
