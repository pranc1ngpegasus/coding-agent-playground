package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace/noop"
)

func NewNoopTracer(ctx context.Context) noop.TracerProvider {
	return noop.NewTracerProvider()
}
