package templates

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	s2aStatuses "github.com/KianIt/swag2api/statuses"
)

func TestHandleResult(t *testing.T) {
	tt := []struct {
		name        string
		w           *httptest.ResponseRecorder
		err         error
		code        int
		responseIn  any
		responseOut any
	}{
		{
			name: "Nil error",
			w:    httptest.NewRecorder(),
			err:  nil,
			code: http.StatusOK,
			responseIn: map[string]any{
				"key": "value",
			},
			responseOut: map[string]any{
				"key": "value",
			},
		},
		{
			name: "Simple error",
			w:    httptest.NewRecorder(),
			err:  errors.New("test error"),
			code: http.StatusInternalServerError,
			responseIn: map[string]any{
				"key": "value",
			},
			responseOut: map[string]any{
				"error": "test error",
			},
		},
		{
			name:       "Status wrapper for nil error and nil response",
			w:          httptest.NewRecorder(),
			err:        s2aStatuses.NotFoundError(nil),
			code:       http.StatusNotFound,
			responseIn: nil,
			responseOut: map[string]any{
				"error": http.StatusText(http.StatusNotFound),
			},
		},
		{
			name: "Status wrapper for nil error and non-nil response",
			w:    httptest.NewRecorder(),
			err:  s2aStatuses.NotFoundError(nil),
			code: http.StatusNotFound,
			responseIn: map[string]any{
				"key": "value",
			},
			responseOut: map[string]any{
				"key": "value",
			},
		},
		{
			name: "Status wrapper for non-nil error",
			w:    httptest.NewRecorder(),
			err:  s2aStatuses.NotFoundError(errors.New("test error")),
			code: http.StatusNotFound,
			responseIn: map[string]any{
				"key": "value",
			},
			responseOut: map[string]any{
				"error": "test error",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_handleResult(tc.w, tc.err, tc.responseIn)

			assert.Equal(t, tc.code, tc.w.Result().StatusCode)

			var response map[string]any
			_ = json.Unmarshal(tc.w.Body.Bytes(), &response)

			assert.Equal(t, tc.responseOut, response)
		})
	}
}
