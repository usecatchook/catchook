package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/theotruvelot/catchook/internal/app"
	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/internal/server"
	"github.com/theotruvelot/catchook/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	appLogger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.Background()

	appLogger.Info(ctx, "Starting Catchook API",
		logger.String("version", "0.0.1"),
		logger.String("go_version", "1.24.5"),
		logger.String("env", os.Getenv("ENV")),
	)

	container, err := app.NewContainer(cfg, appLogger)
	if err != nil {
		appLogger.Fatal(ctx, "Failed to create application container",
			logger.Error(err),
		)
	}
	defer container.Close()

	// Create HTTP server
	httpServer := server.NewServer(container)

	// Channel to listen for interrupt signal to trigger shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			appLogger.Fatal(ctx, "Failed to start server",
				logger.Error(err),
			)
		}
	}()

	// Wait for interrupt signal
	<-quit
	appLogger.Info(ctx, "Shutting down server...")

	// Graceful shutdown
	if err := httpServer.Shutdown(); err != nil {
		appLogger.Error(ctx, "Server forced to shutdown",
			logger.Error(err),
		)
	} else {
		appLogger.Info(ctx, "Server exited gracefully")
	}
}
