package controllers

import (
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

func (h *Pages) Home(c *Ctx) error {
	return h.homeView.Render(c.Res, c)
}
