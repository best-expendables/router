package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/best-expendables/logger"
)

func TestRecoverer(t *testing.T) {
	// Handler which throws the panic
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("panic")
	})

	// Request with context logger
	out := new(bytes.Buffer)
	entry := logger.NewLoggerFactory(logger.InfoLevel, logger.SetOut(out)).Logger(context.TODO())
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(logger.ContextWithEntry(entry, req.Context()))

	t.Run("Default panic handler", func(t *testing.T) {
		rw := httptest.NewRecorder()
		Recoverer(nil)(h).ServeHTTP(rw, req)

		if rw.Code != http.StatusInternalServerError {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusInternalServerError, rw.Code)
		}

		if out.Len() == 0 {
			t.Error("Log record is empty")
		}

		out.Reset()
	})

	t.Run("Custom panic handler", func(t *testing.T) {
		called := false

		rw := httptest.NewRecorder()
		Recoverer(PanicHandlerFunc(func(rw http.ResponseWriter, r *http.Request, err PanicRecoveredError) {
			called = true
			rw.WriteHeader(http.StatusOK)
		}))(h).ServeHTTP(rw, req)

		if !called {
			t.Error("Handler was not called")
		}

		if out.Len() == 0 {
			t.Error("Log record is empty")
		}

		out.Reset()

		if rw.Code != http.StatusOK {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusOK, rw.Code)
		}
	})

	t.Run("No panic", func(t *testing.T) {
		called := false

		rw := httptest.NewRecorder()
		Recoverer(PanicHandlerFunc(func(rw http.ResponseWriter, r *http.Request, err PanicRecoveredError) {
			called = true
		}))(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})).ServeHTTP(rw, req)

		if called {
			t.Error("Panic handler has been called")
		}
	})
}
