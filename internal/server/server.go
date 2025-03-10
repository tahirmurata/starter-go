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

func New(cfg *config.Config) *http.Server {
	newServer := &Server{
		port: cfg.App.Port,

		db: database.New(cfg),
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", newServer.port),
		Handler: newServer.RegisterRoutes(),
		// IdleTimeout:  time.Minute,
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 30 * time.Second,
	}

	return server
}
