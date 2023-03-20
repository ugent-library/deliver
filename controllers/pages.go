package controllers

import (
	"net/http"

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
	return c.RenderHTML(http.StatusOK, views.Home(c))
}
