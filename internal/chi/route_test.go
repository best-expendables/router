package chi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutePatternFromRequest(t *testing.T) {
	patterns := []struct {
		Pattern string
		URL     string
	}{
		{Pattern: "/", URL: "/"},
		{Pattern: "/user/{id}", URL: "/user/1000"},
		{Pattern: "/user/{id}/details", URL: "/user/1000/details"},
	}

	mux := chi.NewMux()

	for _, pattern := range patterns {
		mux.Get(pattern.Pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(RoutePatternFromRequest(r)))
			w.WriteHeader(http.StatusOK)
		}))
	}

	for _, pattern := range patterns {
		r, _ := http.NewRequest(http.MethodGet, pattern.URL, nil)
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, r)

		patternFromRequest := rw.Body.String()
		if patternFromRequest != pattern.Pattern {
			t.Errorf("Pattern from request not equals. Expected '%s', actual '%s'",
				pattern.Pattern, patternFromRequest)
		}
	}
}
