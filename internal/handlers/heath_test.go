package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestLogHandler_HandleHealth(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	handler := NewRequestHandler(logger, 1<<20)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.HandleHealth(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, "healthy", response["status"])
}
