package tracer

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

func Middleware(tracer trace.TracerProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Tracer("").Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			defer span.End()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
