package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ram291/opamp-control-pane/internal/api"
	"github.com/ram291/opamp-control-pane/internal/embedfs"
	"github.com/ram291/opamp-control-pane/internal/supervisor"
	"github.com/ram291/opamp-control-pane/internal/version"
)

func main() {
	configPath := flag.String("config", "configs/supervisor.yaml", "Path to supervisor config file")
	flag.Parse()

	// Print version
	fmt.Println(version.Info())

	// Load configuration
	cfg, err := supervisor.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create logger
	logger := supervisor.NewLogger()

	// Create and start supervisor
	sup := supervisor.New(logger, cfg)
	if err := sup.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start supervisor: %v\n", err)
		os.Exit(1)
	}

	// Start the collector if executable is configured
	if cfg.Agent.Executable != "" {
		if err := sup.StartCollector(); err != nil {
			logger.Errorf(ctx, "Failed to start collector: %v", err)
		}
	}

	// Get frontend filesystem (may be nil if not built)
	var frontendFS = embedfs.GetFrontendFS()
	if frontendFS == nil {
		logger.Debugf(ctx, "React frontend not built. API-only mode.")
	}

	// Create API server
	apiServer := api.New(sup, frontendFS)

	httpServer := &http.Server{
		Addr:    cfg.API.Listen,
		Handler: apiServer.Handler(),
	}

	// Start HTTP server
	go func() {
		logger.Debugf(ctx, "Management API listening on %s", cfg.API.Listen)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf(ctx, "HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Debugf(ctx, "Shutting down...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	_ = httpServer.Shutdown(shutdownCtx)
	_ = sup.Stop(shutdownCtx)

	logger.Debugf(ctx, "Shutdown complete.")
}