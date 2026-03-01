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

var TemplateNames = []string{
	BaseResponse,
	UnmarshalString,
	UnmarshalBytes,
	HandleBadRequest,
	HandleResult,
	WriteResponse,
}

func IsExistingTemplate(name string) bool {
	return utils.Contains(TemplateNames, name)
}
