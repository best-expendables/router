package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/best-expendables/logger"
)

func TestNew(t *testing.T) {
	loggerBuf := new(bytes.Buffer)
	loggerFactory := logger.NewLoggerFactory(logger.InfoLevel, logger.SetOut(loggerBuf))
	config := Configuration{LoggerFactory: loggerFactory}

	mux, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	mux.Get("/ok", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Ok"))
		w.WriteHeader(http.StatusOK)
	}))

	mux.Get("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}))

	t.Run("Ok", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/ok", nil)
		AddAuthenticationHeaders(r)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		if w.Code != http.StatusOK || w.Body.String() != "Ok" {
			t.Error()
		}

		if loggerBuf.Len() < 50 {
			t.Error("Log is too short")
		}


		loggerBuf.Reset()
	})

	t.Run("Panic recovered", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/panic", nil)
		AddAuthenticationHeaders(r)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Error()
		}

		if loggerBuf.Len() < 50 {
			t.Error("Log is too short")
		}

		loggerBuf.Reset()
	})
}

func AddAuthenticationHeaders(r *http.Request) {
	r.Header.Set("X-LEL-User-ID", "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c")
	r.Header.Set("X-LEL-User-NAME", "admin")
	r.Header.Set("X-LEL-User-EMAIL", "admin@lel.com")
	r.Header.Set("X-LEL-User-Roles", "admin,super")
}
