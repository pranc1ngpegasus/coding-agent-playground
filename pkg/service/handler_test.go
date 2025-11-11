package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto/api/v1"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestGreetServiceHandler_Greet(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tracer := noop.NewTracerProvider().Tracer("test")
	handler := NewGreetServiceHandler(logger, tracer)

	tests := []struct {
		name     string
		reqName  string
		wantMsg  string
	}{
		{
			name:    "greet with name",
			reqName: "World",
			wantMsg: "Hello, World!",
		},
		{
			name:    "greet with empty name",
			reqName: "",
			wantMsg: "Hello, !",
		},
		{
			name:    "greet with Japanese name",
			reqName: "世界",
			wantMsg: "Hello, 世界!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := connect.NewRequest(&apiv1.GreetRequest{
				Name: tt.reqName,
			})

			resp, err := handler.Greet(ctx, req)
			if err != nil {
				t.Errorf("Greet() error = %v", err)
				return
			}

			if resp.Msg.Message != tt.wantMsg {
				t.Errorf("Greet() message = %v, want %v", resp.Msg.Message, tt.wantMsg)
			}
		})
	}
}
