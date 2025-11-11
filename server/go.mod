module github.com/example/multimodule/server

go 1.21

require (
connectrpc.com/connect v1.14.0
github.com/example/multimodule/api v0.0.0
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.50.0
go.opentelemetry.io/otel v1.24.0
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0
go.opentelemetry.io/otel/sdk v1.24.0
golang.org/x/sync v0.6.0
)

require (
go.opentelemetry.io/otel/metric v1.24.0 // indirect
go.opentelemetry.io/otel/trace v1.24.0 // indirect
go.opentelemetry.io/proto/otlp v1.1.0 // indirect
)
