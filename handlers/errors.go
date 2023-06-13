package handlers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httpx/render"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	render.HTML(w, http.StatusNotFound, views.PublicPage(c, &views.NotFound{}))
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	http.Redirect(w, r, c.PathTo("login").String(), http.StatusSeeOther)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	render.HTML(w, http.StatusForbidden, views.PublicPage(c, &views.Forbidden{}))
}
