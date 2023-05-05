package controllers

import (
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/views"
)

type PagesController struct{}

func NewPagesController() *PagesController {
	return &PagesController{}
}

func (h *PagesController) Home(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	if c.User != nil {
		c.RedirectTo("spaces")
		return nil
	}
	return c.HTML(http.StatusOK, views.Page(c, &views.Home{}))
}
