package controllers

import (
	"net/http"
)

type Pages struct {
}

func NewPages() *Pages {
	return &Pages{}
}

func (h *Pages) Home(c *Ctx) error {
	if c.User != nil {
		c.RedirectTo("spaces")
		return nil
	}
	return c.HTML(http.StatusOK, "page", "home", nil)
}
