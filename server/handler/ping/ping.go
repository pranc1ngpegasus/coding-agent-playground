package ping

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	pingv1 "github.com/example/multimodule/api/gen/example/ping/v1"
	"github.com/example/multimodule/api/gen/example/ping/v1/pingv1connect"
)

// Service implements the generated PingServiceHandler interface.
type Service struct {
	logger *slog.Logger
}

// NewService creates a new ping service.
func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

// Handler returns the connect HTTP handler for the ping service.
func Handler(svc *Service, opts ...connect.HandlerOption) (string, http.Handler) {
	return pingv1connect.NewPingServiceHandler(svc, opts...)
}

// Ping echoes the request message.
func (s *Service) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	start := time.Now()
	msg := req.Msg.GetMessage()
	if msg == "" {
		msg = "pong"
	}
	resp := connect.NewResponse(&pingv1.PingResponse{Message: msg})
	resp.Header().Set("X-Handled-Duration", time.Since(start).String())
	s.logger.InfoContext(ctx, "handled ping", "message", req.Msg.GetMessage(), "duration", time.Since(start))
	return resp, nil
}

var _ pingv1connect.PingServiceHandler = (*Service)(nil)
