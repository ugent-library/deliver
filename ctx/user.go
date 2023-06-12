package ctx

import (
	"net/http"

	"github.com/ugent-library/httperror"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r)

		if c.User == nil {
			c.HandleError(w, r, httperror.Unauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r)

		if c.User == nil {
			c.HandleError(w, r, httperror.Unauthorized)
			return
		}
		if !c.Permissions.IsAdmin(c.User) {
			c.HandleError(w, r, httperror.Forbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
