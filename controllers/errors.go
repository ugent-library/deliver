package controllers

import (
	"errors"
	"net/http"

	"github.com/ugent-library/deliver/controllers/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
)

type Errors struct {
}

func NewErrors() *Errors {
	return &Errors{}
}

func (h *Errors) Forbidden(c *ctx.Ctx) error {
	return c.HTML(http.StatusForbidden, "layouts/public_page", "errors/forbidden", nil)
}

func (h *Errors) NotFound(c *ctx.Ctx) error {
	return c.HTML(http.StatusNotFound, "layouts/public_page", "errors/not_found", nil)
}

func (h *Errors) HandleError(c *ctx.Ctx, err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	switch httpErr.StatusCode {
	case http.StatusUnauthorized:
		c.RedirectTo("login")
	case http.StatusForbidden:
		if err := h.Forbidden(c); err != nil {
			h.HandleError(c, err)
		}
	case http.StatusNotFound:
		if err := h.NotFound(c); err != nil {
			h.HandleError(c, err)
		}
	default:
		c.Log.Error(err)
		http.Error(c.Res, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
	}
}
