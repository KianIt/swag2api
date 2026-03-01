package templates

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]any{"key": "value"}

	_writeResponse(w, http.StatusOK, response)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	var bodyResponse map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &bodyResponse)

	assert.Equal(t, response, bodyResponse)
}
