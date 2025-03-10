package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"starter/cmd/web"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes

	// API
	api := http.NewServeMux()
	api.HandleFunc("GET /", s.helloWorldHandler)
	api.HandleFunc("GET /health", s.healthHandler)
	mux.Handle("/api/", http.StripPrefix("/api", api))

	// FRONTEND
	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/", templ.Handler(web.HelloForm()))
	mux.HandleFunc("/hello", web.HelloWebHandler)

	return mux
}

func (s *Server) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Failed to write response", slog.Any("err", err))
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	health, err := s.db.Health(ctx)
	if err != nil {
		http.Error(w, "Failed to get db health", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(health)
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Failed to write response", slog.Any("err", err))
	}
}
