package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/view"
)

type Pages struct {
	homeView     view.View
	notFoundView view.View
}

func NewPages() *Pages {
	return &Pages{
		homeView:     view.MustNew("page", "home"),
		notFoundView: view.MustNew("page", "not_found").Status(404),
	}
}

func (c *Pages) Home(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	c.homeView.Render(w, ctx)
}

func (c *Pages) NotFound(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	c.notFoundView.Render(w, ctx)
}
