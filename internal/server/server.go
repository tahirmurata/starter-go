package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"starter/internal/config"
	"starter/internal/database"
	"starter/internal/manager"
)

type Server struct {
	port string

	db database.Database
}

// New creates and returns both the custom Server instance and the http.Server
func New(cfg *config.Config, manager manager.Manager) (*http.Server, error) {
	db, err := database.New(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}

	// Add database shutdown function
	manager.Add(func() {
		slog.Info("closing database connection")
		db.Close()
	})

	newServer := &Server{
		port: cfg.App.Port,
		db:   db,
	}

	// Declare Server config
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", newServer.port),
		Handler: newServer.RegisterRoutes(),
		// IdleTimeout:  time.Minute,
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 30 * time.Second,
	}

	// Add HTTP server shutdown function
	manager.Add(func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		slog.Info("shutting down HTTP server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Warn("HTTP server forced to shutdown with error", "err", err)
		}
		slog.Info("HTTP server shutdown complete")
	})

	return httpServer, nil
}

// DB returns the database service manager
func (s *Server) DB() database.Database {
	return s.db
}
