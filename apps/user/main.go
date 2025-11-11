package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var env struct {
	logLevel string
	port     string
}

func init() {
	flag.StringVar(&env.logLevel, "log-level", "info", "set the logging level")
	flag.StringVar(&env.port, "port", "8080", "set the listen on port")

	flag.VisitAll(func(f *flag.Flag) {
		env, exists := os.LookupEnv(strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))
		if exists {
			f.Value.Set(env)
		}
	})
	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := inject(ctx, env.logLevel)
	if err != nil {
		panic(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := app.tracerProvider.Shutdown(shutdownCtx); err != nil {
			app.logger.ErrorContext(shutdownCtx, "failed to shutdown tracer provider", slog.Any("err", err))
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", app.healthHandler())
	mux.HandleFunc("/example", app.exampleHandler())

	srv := &http.Server{
		Addr:    ":" + env.port,
		Handler: mux,
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.ErrorContext(ctx, "failed to start server", slog.Any("err", err))
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})
	eg.Go(func() error {
		<-egCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			app.logger.ErrorContext(shutdownCtx, "failed to stop server", slog.Any("err", err))
			return fmt.Errorf("failed to stop server: %w", err)
		}

		return nil
	})
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
