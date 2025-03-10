package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"starter/internal/config"
	"starter/internal/logger"
	"starter/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.LogAttrs(context.Background(), slog.LevelInfo, "Shutting down gracefully, press ctrl+c again to force")
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "Server forced to shutdown with error", slog.Any("err", err))
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	cfg := config.New()

	logger.Init(cfg)

	newServer := server.New(cfg)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(newServer, done)

	slog.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Listening on %s", newServer.Addr))
	err := newServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.LogAttrs(context.Background(), logger.LevelFatal, "HTTP newServer error", slog.Any("err", err))
		panic(1)
	}

	// Wait for the graceful shutdown to complete
	<-done

	slog.LogAttrs(context.Background(), slog.LevelInfo, "Graceful shutdown complete")
}
