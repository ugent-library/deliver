package controllers

import (
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

func (h *Pages) Home(c Ctx) error {
	return h.homeView.Render(c.Res, c)
}

func (h *Pages) NotFound(c Ctx) error {
	return h.notFoundView.Render(c.Res, c)
}
