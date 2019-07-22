package chi

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RoutePatternFromRequest return the pattern
// Could be used only after request handling
func RoutePatternFromRequest(r *http.Request) string {
	if ctx := r.Context().Value(chi.RouteCtxKey); ctx != nil {
		if ctx, ok := ctx.(*chi.Context); ok {
			if len(ctx.RoutePatterns) == 0 {
				ctx.Routes.Match(ctx, r.Method, r.URL.Path)
			}

			if len(ctx.RoutePatterns) > 0 {
				return ctx.RoutePattern()
			}
		}
	}

	return ""
}
