package main

import (
	"net/http"
	"os"

	"github.com/joshhunt/postlog/internal/handlers"
	"github.com/joshhunt/postlog/internal/middleware"
	"github.com/joshhunt/postlog/internal/server"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	// Initialize handlers
	handlers := handlers.NewRequestHandler(sugarLogger, 1<<20) // max body size 1 MiB

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HandleHealth)
	mux.HandleFunc("/", handlers.HandlePayload)

	handler := middleware.LoggingMiddleware(sugarLogger)(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(port, handler, sugarLogger)

	if err := srv.Start(); err != nil {
		sugarLogger.Fatalw("server failed to start", "error", err)
	}
}
