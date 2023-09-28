package handlers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	w.WriteHeader(http.StatusNotFound)
	views.NotFound(c).Render(r.Context(), w)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	http.Redirect(w, r, c.PathTo("login").String(), http.StatusSeeOther)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	w.WriteHeader(http.StatusForbidden)
	views.Forbidden(c).Render(r.Context(), w)
}
