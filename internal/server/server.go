package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.SugaredLogger
}

func New(port string, handler http.Handler, logger *zap.SugaredLogger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  time.Duration(10) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second,
			IdleTimeout:  time.Duration(10) * time.Second,
		},
		logger: logger,
	}
}

// Start starts the server and handles graceful shutdown
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		s.logger.Infow("starting HTTP server",
			"address", s.httpServer.Addr,
		)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalw("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	return s.waitForShutdown()
}

// waitForShutdown waits for interrupt signals and performs graceful shutdown
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-quit
	s.logger.Infow("received shutdown signal", "signal", sig)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Errorw("server forced to shutdown", "error", err)
		return err
	}

	s.logger.Infow("server shutdown completed gracefully")
	return nil
}
