package controllers

import (
	"github.com/a-h/templ"
	"github.com/ugent-library/deliver/controllers/ctx"
	"github.com/ugent-library/deliver/views"
)

type Pages struct {
}

func NewPages() *Pages {
	return &Pages{}
}

func (h *Pages) Home(c *ctx.Ctx) error {
	if c.User != nil {
		c.RedirectTo("spaces")
		return nil
	}
	templ.Handler(views.Home(c)).ServeHTTP(c.Res, c.Req)
	return nil
}
