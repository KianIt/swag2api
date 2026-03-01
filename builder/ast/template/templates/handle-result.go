package templates

import (
	"net/http"
)

func _handleResult(w http.ResponseWriter, err error, response any) {
	if err != nil {
		if cu, ok := err.(interface {
			Code() int
			Unwrap() error
		}); ok {
			if e := cu.Unwrap(); e != nil {
				_writeResponse(w, cu.Code(), _errorResponse{Error: e.Error()})
			} else if response == nil {
				_writeResponse(w, cu.Code(), _errorResponse{Error: http.StatusText(cu.Code())})
			} else {
				_writeResponse(w, cu.Code(), response)
			}
			return
		} else {
			_writeResponse(w, http.StatusInternalServerError, _errorResponse{Error: err.Error()})
			return
		}
	}
	_writeResponse(w, http.StatusOK, response)
}
