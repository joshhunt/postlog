package handlers

import (
	"net/http"
	"time"
)

// HandleHealth provides a health check endpoint
func (h *RequestHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "postlog",
	}
	h.writeJSONResponse(w, response, http.StatusOK)
}
