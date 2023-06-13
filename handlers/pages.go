package handlers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httpx/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if c.User != nil {
		http.Redirect(w, r, c.PathTo("spaces").String(), http.StatusSeeOther)
		return
	}
	render.HTML(w, http.StatusOK, views.Page(c, &views.Home{}))
}
