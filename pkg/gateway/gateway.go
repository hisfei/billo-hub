package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Gateway encapsulates the configuration and startup logic for an HTTP server.
type Gateway struct {
	Addr            string
	Handler         http.Handler
	ShutdownTimeout time.Duration
	logger          *zap.Logger
}

// NewGateway creates a Gateway instance with default configuration.
func NewGateway(addr string, handler http.Handler, logger *zap.Logger) *Gateway {
	return &Gateway{
		Addr:            addr,
		Handler:         handler,
		ShutdownTimeout: 5 * time.Second, // Default shutdown timeout is 5 seconds
		logger:          logger,
	}
}

// Start starts the HTTP server. It blocks until the provided context is canceled or the server stops on its own.
// When the context is canceled, it attempts to gracefully shut down the server.
// The return value indicates whether the server started or shut down gracefully.
func (g *Gateway) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    g.Addr,
		Handler: g.Handler,
	}

	// Channel to signal that ListenAndServe has finished (either successfully or with an error)
	serverErr := make(chan error, 1)

	// Start the server in a goroutine to avoid blocking the main thread
	go func() {
		g.logger.Info("HTTP server is starting...", zap.String("addr", g.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// If it's not an http.ErrServerClosed error, it means the server stopped unexpectedly
			serverErr <- fmt.Errorf("failed to start HTTP server: %w", err)
		} else {
			// The server closed normally (e.g., via Shutdown)
			serverErr <- nil
		}
	}()

	// Wait for the external context to be canceled or for the server to stop on its own
	select {
	case err := <-serverErr:
		// The server stopped on its own (e.g., port in use, or normal exit after Shutdown)
		if err != nil {
			g.logger.Error("HTTP server stopped with error", zap.Error(err))
			return err
		}
		g.logger.Info("HTTP server stopped gracefully.")
		return nil // The server exited normally
	case <-ctx.Done():
		// Received an external signal, start graceful shutdown
		g.logger.Info("Shutting down HTTP server due to external signal...")

		// Create a context with a timeout to notify the server to finish current requests within the specified time
		shutdownCtx, cancel := context.WithTimeout(context.Background(), g.ShutdownTimeout)
		defer cancel()

		// Call Shutdown for graceful shutdown
		if err := srv.Shutdown(shutdownCtx); err != nil {
			g.logger.Error("Server forced to shutdown", zap.Error(err))
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		g.logger.Info("HTTP server exited gracefully.")
		return nil
	}
}
