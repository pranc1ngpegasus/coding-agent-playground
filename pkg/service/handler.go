package service

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	apiv1 "github.com/pranc1ngpegasus/coding-agent-playground/pkg/proto/api/v1"
	"go.opentelemetry.io/otel/trace"
)

type GreetServiceHandler struct {
	logger *slog.Logger
	tracer trace.Tracer
}

func NewGreetServiceHandler(logger *slog.Logger, tracer trace.Tracer) *GreetServiceHandler {
	return &GreetServiceHandler{
		logger: logger,
		tracer: tracer,
	}
}

func (h *GreetServiceHandler) Greet(
	ctx context.Context,
	req *connect.Request[apiv1.GreetRequest],
) (*connect.Response[apiv1.GreetResponse], error) {
	ctx, span := h.tracer.Start(ctx, "GreetServiceHandler.Greet")
	defer span.End()

	h.logger.InfoContext(ctx, "received greet request",
		slog.String("name", req.Msg.Name),
	)

	message := fmt.Sprintf("Hello, %s!", req.Msg.Name)

	resp := connect.NewResponse(&apiv1.GreetResponse{
		Message: message,
	})

	h.logger.InfoContext(ctx, "sending greet response",
		slog.String("message", message),
	)

	return resp, nil
}
