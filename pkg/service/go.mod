module github.com/pranc1ngpegasus/coding-agent-playground/pkg/service

go 1.23

require (
	connectrpc.com/connect v1.19.1
	github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto v0.0.0
	go.opentelemetry.io/otel/trace v1.38.0
)

require (
	go.opentelemetry.io/otel v1.38.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto => ../proto
