package templates

import "net/http"

func _handleResult(w http.ResponseWriter, err error, response any) {
	if err != nil {
		if cu, ok := err.(interface {
			Code() int
			Unwrap() error
		}); ok {
			if e := cu.Unwrap(); e != nil {
				_writeResponse(w, cu.Code(), _baseResponse{Code: cu.Code(), Msg: e.Error()})
				return
			} else {
				_writeResponse(w, cu.Code(), response)
				return
			}
		} else {
			_writeResponse(w, http.StatusInternalServerError, _baseResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
	}

	_writeResponse(w, http.StatusOK, response)
}
