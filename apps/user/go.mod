module github.com/pranc1ngpegasus/coding-agent-playground/apps/user

go 1.25

require (
	github.com/pranc1ngpegasus/coding-agent-playground/logger v0.0.0-00010101000000-000000000000
	github.com/pranc1ngpegasus/coding-agent-playground/tracer v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	golang.org/x/sync v0.18.0
)

replace (
	github.com/pranc1ngpegasus/coding-agent-playground/logger => ../../logger
	github.com/pranc1ngpegasus/coding-agent-playground/tracer => ../../tracer
)
