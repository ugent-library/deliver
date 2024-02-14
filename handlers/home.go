package handlers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if c.User != nil {
		http.Redirect(w, r, c.Path("spaces").String(), http.StatusSeeOther)
		return
	}

	views.Home(c).Render(r.Context(), w)
}
