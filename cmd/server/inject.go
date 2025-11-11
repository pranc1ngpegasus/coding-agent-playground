package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto/api/v1/serviceconnect"
	"github.com/pranc1ngpegasus/coding-agent-playground/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// initLogger initializes and returns a slog.Logger
func initLogger() *slog.Logger {
	var level slog.Level
	switch strings.ToLower(env.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}

// initTracer initializes and returns an OpenTelemetry tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(env.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		env.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider, nil
}

// initHandler initializes and returns the HTTP handler
func initHandler(logger *slog.Logger) http.Handler {
	tracer := otel.Tracer(env.ServiceName)

	// Create service handler
	greetHandler := service.NewGreetServiceHandler(logger, tracer)

	// Create mux
	mux := http.NewServeMux()

	// Register Connect handlers with OpenTelemetry interceptor
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		logger.Error("failed to create otel interceptor", slog.Any("error", err))
		// Fall back to no interceptor
		path, handler := serviceconnect.NewGreetServiceHandler(greetHandler)
		mux.Handle(path+"/", handler)
	} else {
		interceptors := connect.WithInterceptors(otelInterceptor)
		path, handler := serviceconnect.NewGreetServiceHandler(greetHandler, interceptors)
		mux.Handle(path+"/", handler)
	}

	// Wrap with h2c for HTTP/2 without TLS
	return h2c.NewHandler(mux, &http2.Server{})
}
