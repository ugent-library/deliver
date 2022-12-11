package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/view"
)

type Pages struct {
	homeView view.View
}

func NewPages() *Pages {
	return &Pages{
		homeView: view.MustNew("page", "home"),
	}
}

func (c *Pages) Home(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	c.homeView.Render(w, ctx)
}
