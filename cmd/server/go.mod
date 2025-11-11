module github.com/pranc1ngpegasus/coding-agent-playground/cmd/server

go 1.23

require (
	connectrpc.com/connect v1.19.1
	connectrpc.com/otelconnect v0.8.0
	github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto v0.0.0
	github.com/pranc1ngpegasus/coding-agent-playground/pkg/service v0.0.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	golang.org/x/net v0.46.0
	golang.org/x/sync v0.17.0
	google.golang.org/grpc v1.70.0
)
