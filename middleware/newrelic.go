package middleware

import (
	"net/http"

	"bitbucket.org/snapmartinc/router/internal/chi"

	"bitbucket.org/snapmartinc/logger"
	"bitbucket.org/snapmartinc/trace"
	"github.com/newrelic/go-agent"
)

// Newrelic create new transaction and set it to request context.
func Newrelic(app newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.URL.Path, w, r)

			// Newrelic and logger has identical attributes names
			if requestID := trace.RequestIDFromContext(r.Context()); requestID != "" {
				txn.AddAttribute(logger.FieldRequestId, requestID)
			}

			next.ServeHTTP(txn, newrelic.RequestWithTransactionContext(r, txn))

			if pattern := chi.RoutePatternFromRequest(r); pattern != "" {
				txn.SetName(r.Method + " " + pattern)
			} else {
				txn.SetName(r.Method + " route not found")
			}

			if err := txn.End(); err != nil {
				logger.EntryFromContextOrDefault(r.Context()).WithFields(logger.Fields{
					"err":       err.Error(),
					"component": "router.middleware.newrelic",
				}).Error("Can not finish NewRelic transaction")
			}
		})
	}
}
