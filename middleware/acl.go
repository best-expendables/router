package middleware

import (
	"net/http"

	"bitbucket.org/snapmartinc/user-service-client"
)

type (
	// AccessChecker checks that's user has rights for access
	Authorizer interface {
		// Authorize request
		// FYI: don't need to parse request path.
		// Request already has information about matched route chi.RouteContext(r.Context())
		Authorize(user *userclient.User, r *http.Request) AccessCheck
	}

	// AuthorizerFn functional implementation of AccessChecker
	AccessCheckerFn func(*userclient.User, *http.Request) AccessCheck

	// AccessCheck result of access check
	AccessCheck struct {
		code int
	}
)

var (
	// AccessGranted successful authorization
	AccessGranted = AccessCheck{}
	// AccessUnauthorized user unauthorized
	AccessUnauthorized = AccessCheck{http.StatusUnauthorized}
	// AccessForbidden users doesn't have required missing role
	AccessForbidden = AccessCheck{http.StatusForbidden}
)

// Authorized grants access for non-guest users
func Authorized(next http.Handler) http.Handler {
	return ACL(AccessCheckerFn(func(u *userclient.User, _ *http.Request) AccessCheck {
		if u == nil {
			return AccessUnauthorized
		}
		return AccessGranted
	}))(next)
}

// OnlyRoles grants access for given roles
func OnlyRoles(roles ...string) func(http.Handler) http.Handler {
	return ACL(AccessCheckerFn(func(u *userclient.User, _ *http.Request) AccessCheck {
		if u == nil {
			return AccessUnauthorized
		}

		if u.HasRole(roles...) {
			return AccessGranted
		}

		return AccessForbidden
	}))
}

// ACL grants access for custom rules
func ACL(authorizer Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := userclient.GetCurrentUserFromContext(r.Context())
			if access := authorizer.Authorize(u, r); access == AccessGranted {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(access.code), access.code)
			}
		})
	}
}

func (fn AccessCheckerFn) Authorize(user *userclient.User, r *http.Request) AccessCheck {
	return fn(user, r)
}
