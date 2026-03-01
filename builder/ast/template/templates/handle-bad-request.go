package templates

import (
	"net/http"
)

func _handleBadRequest(w http.ResponseWriter, err error) {
	var (
		code = http.StatusBadRequest
		msg  = http.StatusText(http.StatusBadRequest)
	)

	if err != nil {
		msg = err.Error()
	}

	response := _errorResponse{
		Error: msg,
	}

	_writeResponse(w, code, response)
}
