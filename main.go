package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	zap "go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Register an HTTP handler for all routes with the logger passed in
	http.HandleFunc("/", requestHandler(sugar))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil, // nil means use DefaultServeMux
		// (Optionally you could set ReadTimeout, WriteTimeout, IdleTimeout here.)
	}

	go func() {
		sugar.Infow("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("ListenAndServe()", "error", err)
		}
	}()

	// Create a channel to receive OS signals.
	quit := make(chan os.Signal, 1)
	// relay incoming SIGINT and SIGTERM to quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal arrives.
	sig := <-quit
	sugar.Infow("Received signal, shutting down server", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("Server forced to shut down", "error", err)
	}

	sugar.Infow("Server exited cleanly")
}

// requestHandler returns a handler function that has access to the logger
func requestHandler(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // limit to 1 MiB
		if err != nil {
			log.Errorw("Error reading body", "error", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var payload map[string]any
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			log.Errorw("Invalid JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		payload["req_path"] = r.URL.Path
		payload["req_method"] = r.Method

		log.Infow("Request received", flattenMap(payload)...)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		responseBytes, err := json.Marshal(payload)
		if err != nil {
			log.Errorw("Error marshaling response", "error", err)
			http.Error(w, "Failed to create response", http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
	}
}

// flattenMap converts a map[string]any to a slice of alternating keys and values
// for structured logging with zap
func flattenMap(m map[string]any) []any {
	result := make([]any, 0, len(m)*2)
	for key, value := range m {
		result = append(result, key, value)
	}
	return result
}
