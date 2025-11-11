package main

import (
	"context"

	loggerPkg "github.com/pranc1ngpegasus/coding-agent-playground/logger"
	tracerPkg "github.com/pranc1ngpegasus/coding-agent-playground/tracer"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
)

type app struct {
	logger         *slog.Logger
	tracerProvider *trace.TracerProvider
}

func inject(ctx context.Context, logLevel string) (*app, error) {
	logger := loggerPkg.NewLogger(logLevel)

	tracerProvider, err := tracerPkg.NewTracer(ctx)
	if err != nil {
		return nil, err
	}

	return &app{
		logger:         logger,
		tracerProvider: tracerProvider,
	}, nil
}
