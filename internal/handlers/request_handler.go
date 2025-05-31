package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type RequestHandler struct {
	logger      *zap.SugaredLogger
	maxBodySize int64
}

func NewRequestHandler(logger *zap.SugaredLogger, maxBodySize int64) *RequestHandler {
	return &RequestHandler{
		logger:      logger,
		maxBodySize: maxBodySize,
	}
}

// writeJSONResponse writes a JSON response
func (h *RequestHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Errorw("error encoding JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes an error response
func (h *RequestHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errorResp := map[string]string{
		"error": message,
		"code":  http.StatusText(statusCode),
	}
	h.writeJSONResponse(w, errorResp, statusCode)
}
