package templates

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		w        *httptest.ResponseRecorder
		err      error
		response _errorResponse
	}{
		{
			name: "Nil error",
			w:    httptest.NewRecorder(),
			err:  nil,
			response: _errorResponse{
				Error: http.StatusText(http.StatusBadRequest),
			},
		},
		{
			name: "Non-nil error",
			w:    httptest.NewRecorder(),
			err:  errors.New("test error"),
			response: _errorResponse{
				Error: "test error",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_handleBadRequest(tc.w, tc.err)

			assert.Equal(t, http.StatusBadRequest, tc.w.Result().StatusCode)

			var response _errorResponse
			_ = json.Unmarshal(tc.w.Body.Bytes(), &response)

			assert.Equal(t, tc.response, response)
		})
	}
}
