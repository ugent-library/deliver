package controllers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/httperror"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := ctx.Get(r.Context())

		if c.User == nil {
			c.HandleError(httperror.Unauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := ctx.Get(r.Context())

		if c.User == nil {
			c.HandleError(httperror.Unauthorized)
			return
		}
		if !c.IsAdmin(c.User) {
			c.HandleError(httperror.Forbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
