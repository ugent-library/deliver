package controllers

import (
	"github.com/ugent-library/deliver/view"
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
	if c.User != nil {
		c.RedirectTo("spaces")
		return nil
	}
	return h.homeView.Render(c.Res, c)
}
