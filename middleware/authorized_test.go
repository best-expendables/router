package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/best-expendables/user-service-client"
)

func TestAuthorized(t *testing.T) {
	userExists := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := userclient.GetCurrentUserFromContext(r.Context()); user != nil {
			userExists = true
		}
		w.WriteHeader(http.StatusOK)
	})

	t.Run("User exists", func(t *testing.T) {
		user := &userclient.User{
			Id:       "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c",
			Username: "admin",
			Email:    "admin@lel.com",
			Roles:    []string{"admin", "super"},
		}

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(userclient.ContextWithUser(req.Context(), user))

		w := httptest.NewRecorder()
		Authorized(next).ServeHTTP(w, req)

		if !userExists {
			t.Error("User not exists")
		}

		userExists = false
	})

	t.Run("User not exists", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		Authorized(next).ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusUnauthorized, w.Code)
		}

		if userExists {
			t.Error("User exists")
		}

		userExists = false
	})
}
