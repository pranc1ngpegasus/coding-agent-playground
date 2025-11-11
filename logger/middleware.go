package logger

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	http.Flusher
	status int
	size   int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		w,
		w.(http.Flusher),
		http.StatusOK,
		0,
	}
}

func (m *loggingResponseWriter) Status() int {
	return m.status
}

func (m *loggingResponseWriter) BytesWritten() int {
	return m.size
}

func (m *loggingResponseWriter) WriteHeader(status int) {
	m.status = status
	m.ResponseWriter.WriteHeader(status)
}

func (m *loggingResponseWriter) Write(b []byte) (int, error) {
	m.size = len(b)

	return m.ResponseWriter.Write(b)
}

func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ww := NewLoggingResponseWriter(w)
			t1 := time.Now()

			defer func() {
				logger.InfoContext(ctx, "Served",
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("duration", time.Since(t1).String()),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
