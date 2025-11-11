package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Container bundles application-wide dependencies.
type Container struct {
	Logger         *slog.Logger
	tracerProvider *sdktrace.TracerProvider
}

// NewContainer wires dependencies used by the HTTP server.
func NewContainer(ctx context.Context, cfg EnvConfig) (*Container, error) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	var (
		exporter *otlptracehttp.Exporter
		err      error
	)
	if cfg.OTLPEndpoint != "" {
		exporter, err = newOTLPExporter(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("create otlp exporter: %w", err)
		}
	}

	resources, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("build resource: %w", err)
	}

	opts := []sdktrace.TracerProviderOption{sdktrace.WithResource(resources)}
	if exporter != nil {
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}
	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Container{Logger: logger, tracerProvider: tp}, nil
}

func newOTLPExporter(ctx context.Context, cfg EnvConfig) (*otlptracehttp.Exporter, error) {
	clientOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
	}
	if cfg.OTLPTimeout > 0 {
		clientOpts = append(clientOpts, otlptracehttp.WithTimeout(cfg.OTLPTimeout))
	}
	if cfg.OTLPInsecure {
		clientOpts = append(clientOpts, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

// Shutdown releases resources managed by the container.
func (c *Container) Shutdown(ctx context.Context) error {
	var err error
	if c.tracerProvider != nil {
		shutdownErr := c.tracerProvider.Shutdown(ctx)
		if shutdownErr != nil && !errors.Is(shutdownErr, context.DeadlineExceeded) {
			err = errors.Join(err, shutdownErr)
		}
	}
	return err
}
