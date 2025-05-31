package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (h *RequestHandler) handlePost(r *http.Request) (map[string]any, *AppError) {
	payload := map[string]any{}

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		// Read request body
		bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, h.maxBodySize))
		if err != nil {
			appErr := NewAppError("failed to read request body", http.StatusBadRequest, "max_body_size", h.maxBodySize)
			return nil, appErr
		}

		// Parse JSON
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			appErr := NewAppError("invalid json in request body", http.StatusBadRequest, "body_size", len(bodyBytes))
			return nil, appErr
		}
	default:
		// TODO - handle other content types
		err := NewAppError("invalid content type", http.StatusUnsupportedMediaType, "content_type", contentType)
		return nil, err
	}

	return payload, nil
}

func (h *RequestHandler) handleGet(r *http.Request) (map[string]any, *AppError) {
	payload := map[string]any{}

	// Extract query parameters
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		if len(values) == 1 {
			// Single value - store as string
			payload[key] = values[0]
		} else if len(values) > 1 {
			// Multiple values - store as array
			payload[key] = values
		}
	}

	return payload, nil
}

// HandlePayload processes incoming log requests
func (h *RequestHandler) HandlePayload(res http.ResponseWriter, req *http.Request) {
	var payload map[string]any
	var appErr *AppError

	switch req.Method {
	case http.MethodPost:
		payload, appErr = h.handlePost(req)
		if appErr != nil {
			h.HandleAppError(res, appErr)
			return
		}
	case http.MethodGet:
		payload, appErr = h.handleGet(req)
		if appErr != nil {
			h.HandleAppError(res, appErr)
			return
		}
	default:
		err := NewAppError("Method not allowed", http.StatusMethodNotAllowed, "method", req.Method)
		h.HandleAppError(res, err)
		return
	}

	// Add request metadata
	path := strings.TrimPrefix(strings.TrimSuffix(req.URL.Path, "/"), "/")
	payload["payload_name"] = path
	payload["payload_method"] = req.Method

	// Log the request with structured fields
	h.logger.Infow("received payload", flattenPayload(payload)...)

	h.writeJSONResponse(res, payload, http.StatusOK)
}

// logRequest logs the incoming request with flattened fields
func flattenPayload(payload map[string]any) []any {
	fields := make([]any, 0, len(payload)*2)
	for key, value := range payload {
		fields = append(fields, key, value)
	}

	return fields
}
