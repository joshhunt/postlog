package handlers

import (
	"net/http"
)

// AppError is a simple error with HTTP status and logging properties
type AppError struct {
	Message    string
	StatusCode int
	Properties []any
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a simple app error
func NewAppError(message string, statusCode int, keysAndValues ...any) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Properties: keysAndValues,
	}
}

// Handle the error by logging it and writing HTTP response
func (h *RequestHandler) HandleAppError(w http.ResponseWriter, err *AppError) {
	// Log it
	properties := append([]any{"error_message", err.Message}, err.Properties...)
	h.logger.Errorw("request error", properties...)

	// Write HTTP response
	response := map[string]string{
		"error":  err.Message,
		"status": http.StatusText(err.StatusCode),
	}
	h.writeJSONResponse(w, response, err.StatusCode)
}
