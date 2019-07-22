package middleware

import (
	"fmt"
	"net/http"

	"bitbucket.org/snapmartinc/trace"
	"github.com/go-chi/chi/middleware"
)

// RequestID middleware takes request id from headers, set it into the context and pass to response headers.
func RequestID(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := trace.RequestIDFromHeader(r.Header)

			if requestID == "" {
				if prefix != "" {
					requestID = prefix + "/"
				}

				requestID += fmt.Sprintf("%06d", middleware.NextRequestID())
			}
			ctx := trace.ContextWithRequestID(r.Context(), requestID)
			// Send old or new response id back to the client
			trace.RequestIDToHeader(w.Header(), requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
