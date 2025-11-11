module github.com/pranc1ngpegasus/coding-agent-playground/cmd/server

go 1.24.0

toolchain go1.24.7

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
	google.golang.org/grpc v1.75.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace (
	github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto => ../../pkg/proto
	github.com/pranc1ngpegasus/coding-agent-playground/pkg/service => ../../pkg/service
)
