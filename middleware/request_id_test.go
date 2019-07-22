package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bitbucket.org/snapmartinc/trace"
)

func TestRequestID_FromHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID := trace.RequestIDFromContext(r.Context()); requestID == "" {
			t.Error("Request id has not presented into context")
		}
		if requestID := trace.RequestIDFromHeader(r.Header); requestID == "" {
			t.Error("Request id has not presented into header")
		}

		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	trace.RequestIDToHeader(req.Header, "REQUEST-ID-TEST")

	w := httptest.NewRecorder()
	RequestID("")(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Response code is not 200")
	}

	if requestID := trace.RequestIDFromHeader(w.Header()); requestID != "REQUEST-ID-TEST" {
		t.Errorf("Request id not equals for response headers. Expected '%s', actual '%s'",
			"REQUEST-ID-TEST", requestID)
	}
}

func TestRequestID_CreateNew_NoPrefix(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID := trace.RequestIDFromContext(r.Context()); requestID == "" {
			t.Error("Request id has not presented into context")
		}
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)

	w := httptest.NewRecorder()
	RequestID("")(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Response code is not 200")
	}

	if requestID := trace.RequestIDFromHeader(w.Header()); requestID == "" {
		t.Error("Response headers doesn't contains request id")
	}
}

func TestRequestID_CreateNew_WithPrefix(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID := trace.RequestIDFromContext(r.Context()); requestID == "" {
			t.Error("Request id has not presented into context")
		}

		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)

	w := httptest.NewRecorder()
	RequestID("TEST")(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Response code is not 200")
	}

	requestID := trace.RequestIDFromHeader(w.Header())

	if requestID == "" {
		t.Fatal("Response headers doesn't contains request id")
	}

	if !strings.HasPrefix(requestID, "TEST") {
		t.Fatal("Request id doesn't has a 'TEST'")
	}
}
