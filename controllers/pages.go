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

func (c *Pages) Home(ctx Ctx) error {
	return c.homeView.Render(ctx.Res, ctx)
}

func (c *Pages) NotFound(ctx Ctx) error {
	return c.notFoundView.Render(ctx.Res, ctx)
}
