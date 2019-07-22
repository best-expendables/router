package middleware

import (
	"net/http"

	"bitbucket.org/snapmartinc/logger"
)

// ContextLogger initialize logger entry for request context
//
// Required for all routers as first middleware!
// At that moment we don't see reason pass loggerFactory to the each middleware -
// therefore all middlewares uses logger from context, otherwise will be used global logger and
// information could be incomplete.
func ContextLogger(factory logger.Factory) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			entry := factory.Logger(ctx)
			entry = entry.WithField("input-url", r.URL.String())
			ctx = logger.ContextWithEntry(entry, ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
