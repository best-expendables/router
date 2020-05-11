package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/best-expendables/user-service-client"
	"github.com/go-chi/chi"
)

func TestAuthorized_Success(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "net://example.com/foo", nil)
	req = req.WithContext(userclient.ContextWithUser(context.TODO(), &userclient.User{}))

	w := httptest.NewRecorder()
	Authorized(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code doesn't equals. Expected %d, actual %d", http.StatusOK, w.Code)
	}
}

func TestAuthorized_Failed(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "net://example.com/foo", nil)

	w := httptest.NewRecorder()
	Authorized(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code doesn't equals. Expected %d, actual %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOnlyRoles_Success(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	user := userclient.User{
		Roles: []string{
			userclient.RoleAdmin,
			userclient.RoleShippingProvider,
		},
	}

	req := httptest.NewRequest("GET", "net://example.com/foo", nil)
	req = req.WithContext(userclient.ContextWithUser(context.TODO(), &user))

	w := httptest.NewRecorder()
	OnlyRoles(userclient.RoleAdmin)(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code doesn't equals. Expected %d, actual %d", http.StatusOK, w.Code)
	}
}

func TestOnlyRoles_Failed_WrongRoles(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	user := userclient.User{
		Roles: []string{
			userclient.RoleAdmin,
			userclient.RoleShippingProvider,
		},
	}

	req := httptest.NewRequest("GET", "net://example.com/foo", nil)
	req = req.WithContext(userclient.ContextWithUser(context.TODO(), &user))

	w := httptest.NewRecorder()
	OnlyRoles(userclient.RolePlatformWarehouse)(next).ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code doesn't equals. Expected %d, actual %d", http.StatusOK, w.Code)
	}
}

func TestOnlyRoles_Failed_Unauthorized(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "net://example.com/foo", nil)

	w := httptest.NewRecorder()
	OnlyRoles(userclient.RolePlatformWarehouse)(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code doesn't equals. Expected %d, actual %d", http.StatusOK, w.Code)
	}
}

func TestOnlyRoles_WithRouter(t *testing.T) {
	router := chi.NewMux()

	router.Get("/guest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router.Group(func(r chi.Router) {
		r.Use(OnlyRoles(userclient.RoleAdmin, userclient.RoleAccessAdmin))
	}).Get("/role-admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Access granted for non-restricted
	{
		r, _ := http.NewRequest(http.MethodGet, "/guest", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusOK, w.Code)
		}
	}

	// Access granted for restricted route
	{
		user := userclient.User{
			Roles: []string{userclient.RoleAdmin},
		}

		r, _ := http.NewRequest(http.MethodGet, "/role-admin", nil)
		r = r.WithContext(userclient.ContextWithUser(r.Context(), &user))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusOK, w.Code)
		}
	}

	//  Guest not allowed
	{
		r, _ := http.NewRequest(http.MethodGet, "/role-admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusUnauthorized, w.Code)
		}
	}

	// User doesn't contains required role
	{
		user := userclient.User{
			Roles: []string{
				userclient.RolePlatformOverview,
				userclient.RoleReturn,
			},
		}

		r, _ := http.NewRequest(http.MethodGet, "/role-admin", nil)
		r = r.WithContext(userclient.ContextWithUser(r.Context(), &user))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("Response code not equals. Expected %d, actual %d", http.StatusForbidden, w.Code)
		}
	}
}
