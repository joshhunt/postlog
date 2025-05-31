package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// HandlePayload processes incoming log requests
func (h *RequestHandler) HandlePayload(w http.ResponseWriter, r *http.Request) {

	// Read request body
	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, h.maxBodySize))
	if err != nil {
		h.logger.Errorw("failed to read request body", "error", err)
		h.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		h.logger.Errorw("invalid json in request body", "error", err)
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Add request metadata
	payload["req_path"] = r.URL.Path
	payload["req_method"] = r.Method
	payload["received_at"] = time.Now().UTC()

	// Log the request with structured fields
	h.logRequest(payload)

	// Create response
	response := map[string]interface{}{
		"status":       "received",
		"received_at":  time.Now().UTC(),
		"request_path": r.URL.Path,
		"method":       r.Method,
		"data":         payload,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// logRequest logs the incoming request with flattened fields
func (h *RequestHandler) logRequest(payload map[string]interface{}) {
	fields := make([]interface{}, 0, len(payload)*2)
	for key, value := range payload {
		fields = append(fields, key, value)
	}

	h.logger.Infow("received payload", fields...)
}
