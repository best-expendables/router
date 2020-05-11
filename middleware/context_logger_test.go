package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/best-expendables/logger"
	"github.com/best-expendables/trace"
)

func TestContextLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	loggerFactory := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if entry := logger.EntryFromContext(r.Context()); entry != nil {
			entry.Info("Log!")
		} else {
			t.Fatal("Logger not presented into context")
		}

		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	ctx := trace.ContextWithRequestID(req.Context(), "REQUEST-ID-TEST")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ContextLogger(loggerFactory)(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code not equals. Expected 200, actual %d", w.Code)
	}

	log := buf.String()

	if !strings.Contains(log, "\"input-url\":\"http://localhost\"") {
		t.Error("'input-url' not presented")
	}
}
