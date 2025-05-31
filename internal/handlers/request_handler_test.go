package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestPayloadHandler_HandlePayload(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	handler := NewRequestHandler(logger, 1<<20)

	tests := []struct {
		name            string
		method          string
		path            string
		contentType     string
		body            string
		expectedStatus  int
		expectedPayload map[string]any
	}{
		{
			name:           "valid json request",
			method:         "POST",
			contentType:    "application/json",
			path:           "/foo",
			body:           `{"message": "test log", "level": "info"}`,
			expectedStatus: http.StatusOK,
			expectedPayload: map[string]any{
				"payload_name":   "foo",
				"payload_method": "POST",
				"message":        "test log",
				"level":          "info",
			},
		},
		{
			name:           "valid json request with subpath",
			method:         "POST",
			contentType:    "application/json",
			path:           "/foo/bar",
			body:           `{"message": "test log", "level": "info"}`,
			expectedStatus: http.StatusOK,
			expectedPayload: map[string]any{
				"payload_name":   "foo/bar",
				"payload_method": "POST",
				"message":        "test log",
				"level":          "info",
			},
		},
		{
			name:           "valid get request",
			method:         "GET",
			contentType:    "application/json",
			path:           "/foo/bar?hello=world&foo=69",
			expectedStatus: http.StatusOK,
			expectedPayload: map[string]any{
				"payload_name":   "foo/bar",
				"payload_method": "GET",
				"hello":          "world",
				"foo":            "69",
			},
		},
		{
			name:           "invalid json request",
			method:         "POST",
			contentType:    "application/json",
			path:           "/foo",
			body:           `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			method:         "POST",
			contentType:    "application/json",
			path:           "/foo",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.method, testCase.path, bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", testCase.contentType)

			rr := httptest.NewRecorder()
			handler.HandlePayload(rr, req)

			assert.Equal(t, testCase.expectedStatus, rr.Code)

			if testCase.expectedStatus == http.StatusOK {
				var response map[string]any
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.Nil(t, err)

				assert.Equal(t, testCase.expectedPayload, response)
			}
		})
	}
}
