package middleware

import (
	"net/http"
	"strings"

	"github.com/best-expendables/user-service-client"
)

// Authentication looking for user in headers and inject user to context
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := userclient.User{
			Id:       r.Header.Get("X-Sm-User-ID"),
			Username: r.Header.Get("X-Sm-User-Name"),
			Email:    r.Header.Get("X-Sm-User-Email"),
			Roles:    strings.Split(r.Header.Get("X-Sm-User-Roles"), ","),
		}

		if user.Id == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := userclient.ContextWithUser(r.Context(), &user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
