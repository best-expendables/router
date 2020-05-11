package middleware

import (
	"net/http"

	"github.com/best-expendables/logger"
	"github.com/best-expendables/router/internal/accesslog"
)

// AccessLogOptions for request or response body
type AccessLogOptions struct {
	IgnoreRequestBody  bool
	IgnoreResponseBody bool
}

// AccessLog records all requests processed by the server.
func AccessLog(opts AccessLogOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writerOpt := accesslog.AccessWriterOptions{
				IgnoreRequestBody:  opts.IgnoreRequestBody,
				IgnoreResponseBody: opts.IgnoreResponseBody,
			}

			writer := accesslog.New(writerOpt, w, r)
			next.ServeHTTP(writer, r)
			accessLog := writer.Entry()

			entry := logger.EntryFromContextOrDefault(r.Context()).WithFields(logger.Fields{
				"request":  accessLog.Request,
				"response": accessLog.Response,
			})

			code := accessLog.Response.StatusCode
			if code < 400 || code == http.StatusUnprocessableEntity {
				entry.Info(http.StatusText(code))
			} else {
				entry.Error(http.StatusText(code))
			}
		})
	}
}
