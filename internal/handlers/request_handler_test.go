package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestPayloadHandler_HandlePayload(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	handler := NewRequestHandler(logger, 1<<20)

	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid json request",
			method:         "POST",
			contentType:    "application/json",
			body:           `{"message": "test log", "level": "info"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json request",
			method:         "POST",
			contentType:    "application/json",
			body:           `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			method:         "POST",
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			handler.HandlePayload(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				if response["status"] != "received" {
					t.Errorf("expected status 'received', got %v", response["status"])
				}
			}
		})
	}
}

func TestLogHandler_HandleHealth(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	handler := NewRequestHandler(logger, 1<<20)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.HandleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", response["status"])
	}
}
