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

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"

	"github.com/example/multimodule/api/gen/example/ping/v1/pingv1connect"
	pinghandler "github.com/example/multimodule/server/handler/ping"
)

// Env holds the runtime configuration for the service.
var Env EnvConfig

// EnvConfig enumerates all configuration sourced from flags or the environment.
type EnvConfig struct {
	ListenAddress   string
	ShutdownTimeout time.Duration
	ServiceName     string
	ServiceVersion  string
	OTLPEndpoint    string
	OTLPInsecure    bool
	OTLPTimeout     time.Duration
}

func init() {
	Env = EnvConfig{
		ListenAddress:   ":8080",
		ShutdownTimeout: 10 * time.Second,
		ServiceName:     "connectrpc-server",
		ServiceVersion:  "0.1.0",
		OTLPEndpoint:    "",
		OTLPInsecure:    true,
		OTLPTimeout:     5 * time.Second,
	}

	flag.StringVar(&Env.ListenAddress, "listen", Env.ListenAddress, "HTTP listen address")
	flag.DurationVar(&Env.ShutdownTimeout, "shutdown-timeout", Env.ShutdownTimeout, "graceful shutdown timeout")
	flag.StringVar(&Env.ServiceName, "service-name", Env.ServiceName, "service name for telemetry")
	flag.StringVar(&Env.ServiceVersion, "service-version", Env.ServiceVersion, "service version for telemetry")
	flag.StringVar(&Env.OTLPEndpoint, "otlp-endpoint", Env.OTLPEndpoint, "OTLP endpoint address")
	flag.BoolVar(&Env.OTLPInsecure, "otlp-insecure", Env.OTLPInsecure, "disable TLS for OTLP exporter")
	flag.DurationVar(&Env.OTLPTimeout, "otlp-timeout", Env.OTLPTimeout, "timeout for OTLP exporter communication")

	applyEnvDefaults("APP_")

	flag.Parse()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := NewContainer(context.Background(), Env)
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
		logger.Error("failed to construct container", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), Env.ShutdownTimeout)
		defer cancel()
		if err := container.Shutdown(shutdownCtx); err != nil {
			container.Logger.Error("failed to shutdown container", "error", err)
		}
	}()

	mux := http.NewServeMux()

	pingSvc := pinghandler.NewService(container.Logger)
	svcPath, svcHandler := pinghandler.Handler(pingSvc)
	mux.Handle(svcPath, otelhttp.NewHandler(svcHandler, pingv1connect.PingServiceName))

	srv := &http.Server{Addr: Env.ListenAddress, Handler: mux}

	var eg errgroup.Group

	eg.Go(func() error {
		container.Logger.Info("starting http server", "addr", Env.ListenAddress)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		container.Logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), Env.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		container.Logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}

	container.Logger.Info("server exited cleanly")
}

func applyEnvDefaults(prefix string) {
	flag.VisitAll(func(f *flag.Flag) {
		envName := prefix + strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		if value, ok := os.LookupEnv(envName); ok {
			if err := flag.Set(f.Name, value); err != nil {
				fmt.Fprintf(os.Stderr, "invalid value for %s: %v\n", envName, err)
			}
		}
	})
}
