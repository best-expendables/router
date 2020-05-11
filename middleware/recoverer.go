package middleware

import (
	"net/http"
	"runtime/debug"

	"fmt"

	"github.com/best-expendables/logger"
	"github.com/newrelic/go-agent"
)

type (
	// PanicHandler interface
	// Triggered on application panic and passing as last argument (rvr) value from recover()
	PanicHandler interface {
		ServeHTTP(w http.ResponseWriter, req *http.Request, err PanicRecoveredError)
	}

	// PanicHandlerFunc implements PanicHandler interface
	PanicHandlerFunc func(http.ResponseWriter, *http.Request, PanicRecoveredError)
)

// Recoverer panic error handling
func Recoverer(handler PanicHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					err := PanicRecoveredError{Rvr: rvr, Stack: debug.Stack()}

					logger.EntryFromContextOrDefault(r.Context()).WithFields(logger.Fields{
						"input-url": r.URL.String(),
						"err":       rvr,
						"stack":     string(err.Stack),
						"component": "router.middleware.recoverer",
					}).Error("Recovered")

					if tnx, ok := w.(newrelic.Transaction); ok {
						_ = tnx.NoticeError(&err)
					}

					if handler == nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					handler.ServeHTTP(w, r, err)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// ServeHTTP implements PanicHandler interface
func (f PanicHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, err PanicRecoveredError) {
	f(w, r, err)
}

type PanicRecoveredError struct {
	Rvr   interface{}
	Stack []byte
}

func (err PanicRecoveredError) Error() string {
	return fmt.Sprintf("panic: '%+v'. Stack: %s", err.Rvr, err.Stack)
}
