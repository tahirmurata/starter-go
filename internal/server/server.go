package server

import (
	"fmt"
	"net/http"

	"starter/internal/config"
	"starter/internal/database"
)

type Server struct {
	port string

	db database.ServiceManager
}

func New(cfg *config.Config) (*http.Server, error) {
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}

	newServer := &Server{
		port: cfg.App.Port,

		db: db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", newServer.port),
		Handler: newServer.RegisterRoutes(),
		// IdleTimeout:  time.Minute,
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 30 * time.Second,
	}

	return server, nil
}
