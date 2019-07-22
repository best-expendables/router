package middleware

import (
	"net/http"
	"strings"

	"bitbucket.org/snapmartinc/user-service-client"
)

// Authentication looking for user in headers and inject user to context
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := userclient.User{
			Id:       r.Header.Get("X-LEL-User-ID"),
			Username: r.Header.Get("X-LEL-User-Name"),
			Email:    r.Header.Get("X-LEL-User-Email"),
			Roles:    strings.Split(r.Header.Get("X-LEL-User-Roles"), ","),
		}

		if user.Id == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := userclient.ContextWithUser(r.Context(), &user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
