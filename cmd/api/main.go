package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"starter/internal/config"
	"starter/internal/logger"
	"starter/internal/manager"
	"starter/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.LogAttrs(context.Background(), logger.LevelFatal, "Failed to get config", slog.Any("err", err))
		panic(1)
	}

	logger.Init(cfg)

	// Initialize service manager
	serviceManager := manager.Init()

	// Get the custom server and HTTP server
	httpSrvInstance, err := server.New(cfg, serviceManager)
	if err != nil {
		slog.LogAttrs(context.Background(), logger.LevelFatal, "Failed to get server", slog.Any("err", err))
		panic(1)
	}

	// Start the server in a goroutine
	go func() {
		slog.Info("Starting HTTP server", "addr", httpSrvInstance.Addr)
		err := httpSrvInstance.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.LogAttrs(context.Background(), logger.LevelFatal, "HTTP server error", slog.Any("err", err))
			serviceManager.WaitForStop()
			panic(1)
		}
	}()

	// Wait for signal to shutdown
	serviceManager.WaitForSignal()
}
