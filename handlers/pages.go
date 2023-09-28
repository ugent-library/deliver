package handlers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if c.User != nil {
		http.Redirect(w, r, c.PathTo("spaces").String(), http.StatusSeeOther)
		return
	}

	views.HomePage(c).Render(r.Context(), w)
}
