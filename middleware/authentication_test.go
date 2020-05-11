package middleware

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/best-expendables/user-service-client"
)

func TestAuthentication_Success(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := userclient.GetCurrentUserFromContext(r.Context())
		if user.Id != "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c" {
			t.Error("Should get user id from request context")
		}
		if user.Username != "admin" {
			t.Error("Should get user name")
		}
		if user.Email != "admin@lel.com" {
			t.Error("Should get user email")
		}
		if user.Roles[0] != "admin" && user.Roles[1] != "super" {
			t.Error("Should get user roles")
		}
	})

	req := httptest.NewRequest("GET", "net://example.com/foo", nil)
	req.Header.Set("X-LEL-User-ID", "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c")
	req.Header.Set("X-LEL-User-NAME", "admin")
	req.Header.Set("X-LEL-User-EMAIL", "admin@lel.com")
	req.Header.Set("X-LEL-User-Roles", "admin,super")

	w := httptest.NewRecorder()
	Authentication(next).ServeHTTP(w, req)
}
