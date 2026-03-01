package templates

import (
	"encoding/json"
	"net/http"
)

func _writeResponse(w http.ResponseWriter, code int, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}
