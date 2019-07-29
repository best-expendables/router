package middleware

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"bitbucket.org/snapmartinc/logger"
	"bitbucket.org/snapmartinc/newrelic-context/nrmock"
	"bitbucket.org/snapmartinc/trace"
	"github.com/go-chi/chi"
)

func TestNewrelic(t *testing.T) {
	patterns := []string{
		"/",
		"/user/{id}",

		"/location/*",
		"/location/{id}/details",
	}

	newrelic := &nrmock.NewrelicApp{}

	router := chi.NewMux()
	router.Use(Newrelic(newrelic))

	for _, pattern := range patterns {
		router.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}

	t.Run("Route exists", func(t *testing.T) {
		checks := [][2]string{
			// 1st element URL, 2td assertion for transaction name
			{"/", "GET /"},

			{"/user/1", "GET /user/{id}"},
			{"/user/08cb1330-729e-4538-b9d6-b9d46fa43716", "GET /user/{id}"},

			{"/location/1", "GET /location/*"},
			{"/location/08cb1330-729e-4538-b9d6-b9d46fa43716", "GET /location/*"},

			{"/location/1/details", "GET /location/{id}/details"},
			{"/location/08cb1330-729e-4538-b9d6-b9d46fa43716/details", "GET /location/{id}/details"},
		}

		for _, check := range checks {
			url := check[0]
			expectedTxName := check[1]
			expectedRequestID := "REQUEST-ID"

			r, _ := http.NewRequest(http.MethodGet, url, nil)
			ctx := trace.ContextWithRequestID(r.Context(), expectedRequestID)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, r.WithContext(ctx))

			if rw.Code != http.StatusOK {
				t.Errorf("Wrong response code - '%d', for url '%s'", rw.Code, url)
			}

			// TODO IMPLEMENT TEST FOR TRACE ID
			if txName := newrelic.Tnx.GetName(); txName != expectedTxName {
				t.Errorf("Transcation name not equals. Expected '%s', actual '%s'",
					expectedTxName, txName)
			}

			if val, ok := newrelic.Tnx.GetAttribute(logger.FieldRequestId); ok {
				if requestID, _ := val.(string); requestID != expectedRequestID {
					t.Errorf("Request ID not equals. Expected '%s', actual '%s'",
						expectedRequestID, requestID)
				}
			} else {
				t.Error("Request ID doesn't exists")
			}
		}
	})

	t.Run("Subroute", func(t *testing.T) {
		router.Route("/sub", func(r chi.Router) {
			r.Get("/details", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			})

			r.Get("/{id}/details", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			})
		})

		r, _ := http.NewRequest(http.MethodGet, "/sub/100/details", nil)

		router.ServeHTTP(httptest.NewRecorder(), r)

		expected := "GET /sub/{id}/details"
		if txName := newrelic.Tnx.GetName(); txName != expected {
			t.Errorf("Transaction name not equals. Expected '%s', actual '%s'", expected, txName)
		}
	})

	t.Run("Route not exists", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/not-exists", nil)
		router.ServeHTTP(httptest.NewRecorder(), r)

		expected := "GET route not found"
		if txName := newrelic.Tnx.GetName(); txName != expected {
			t.Errorf("Transaction name not equals. Expected '%s', actual '%s'", expected, txName)
		}
	})

	t.Run("Access denied", func(t *testing.T) {
		router := chi.NewMux()
		router.Use(Newrelic(newrelic), Authorized)
		router.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		r, _ := http.NewRequest(http.MethodGet, "/user/1", nil)

		router.ServeHTTP(httptest.NewRecorder(), r)

		expected := "GET /user/{id}"
		if txName := newrelic.Tnx.GetName(); txName != expected {
			t.Errorf("Transaction name not equals. Expected '%s', actual '%s'", expected, txName)
		}
	})

	t.Run("Panic recovered", func(t *testing.T) {
		router := chi.NewMux()
		router.Use(Newrelic(newrelic), Recoverer(nil))
		router.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
			panic("panic")
		})

		r, _ := http.NewRequest(http.MethodGet, "/user/1", nil)

		router.ServeHTTP(httptest.NewRecorder(), r)

		expected := "GET /user/{id}"
		if txName := newrelic.Tnx.GetName(); txName != expected {
			t.Errorf("Transaction name not equals. Expected '%s', actual '%s'", expected, txName)
		}
	})
}
