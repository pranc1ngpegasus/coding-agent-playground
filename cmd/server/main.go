package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger := initLogger()
	logger.Info("starting server",
		slog.String("addr", env.ServerAddr),
		slog.String("service_name", env.ServiceName),
	)

	// Initialize tracer
	tracerProvider, err := initTracer(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shutdown tracer provider", slog.Any("error", err))
		}
	}()

	// Initialize HTTP handler
	handler := initHandler(logger)

	// Create HTTP server
	server := &http.Server{
		Addr:    env.ServerAddr,
		Handler: handler,
	}

	// Setup errgroup for concurrent execution
	eg, egCtx := errgroup.WithContext(ctx)

	// Start HTTP server
	eg.Go(func() error {
		logger.Info("HTTP server listening", slog.String("addr", env.ServerAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	// Handle signals for graceful shutdown
	eg.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		select {
		case sig := <-sigCh:
			logger.Info("received signal", slog.String("signal", sig.String()))
		case <-egCtx.Done():
			logger.Info("context cancelled")
		}

		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			time.Duration(env.ShutdownTimeout)*time.Second,
		)
		defer shutdownCancel()

		logger.Info("shutting down server gracefully",
			slog.Int("timeout_seconds", env.ShutdownTimeout),
		)

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		logger.Info("server shutdown complete")
		return nil
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
