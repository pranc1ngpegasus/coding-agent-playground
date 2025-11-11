package main

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (a *app) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("user-service")

		// Create a span for this request
		ctx, span := tracer.Start(ctx, "health_check")
		defer span.End()

		// Add some attributes to the span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		// Simulate some work
		time.Sleep(10 * time.Millisecond)

		// Log the health check
		a.logger.InfoContext(ctx, "health check called",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		// Set span status to OK
		span.SetStatus(codes.Ok, "health check successful")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func (a *app) exampleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("user-service")

		// Create a span for this request
		ctx, span := tracer.Start(ctx, "example_operation")
		defer span.End()

		// Add attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)

		// Create a child span to demonstrate nested spans
		_, childSpan := tracer.Start(ctx, "process_data")
		time.Sleep(20 * time.Millisecond)
		childSpan.SetAttributes(attribute.String("operation", "data_processing"))
		childSpan.End()

		// Another child span
		_, childSpan2 := tracer.Start(ctx, "validate_data")
		time.Sleep(15 * time.Millisecond)
		childSpan2.SetAttributes(attribute.String("operation", "data_validation"))
		childSpan2.End()

		a.logger.InfoContext(ctx, "example operation completed")

		span.SetStatus(codes.Ok, "example operation successful")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"example operation completed"}`))
	}
}
