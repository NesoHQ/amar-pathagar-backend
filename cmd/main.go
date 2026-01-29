package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/online-library/internal/config"
	"github.com/yourusername/online-library/internal/infrastructure/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Server.Mode)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("starting amar pathagar backend",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

	// Create context with cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run server
	if err := run(ctx, cfg, log); err != nil {
		log.Error("server exited with error", zap.Error(err))
		os.Exit(1)
	}

	log.Info("server shutdown complete")
}
