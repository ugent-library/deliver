package controllers

import (
	"errors"
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httperror"
)

type ErrorsController struct{}

func NewErrorsController() *ErrorsController {
	return &ErrorsController{}
}

func (h *ErrorsController) Forbidden(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	return c.HTML(http.StatusForbidden, views.PublicPage(c, &views.Forbidden{}))
}

func (h *ErrorsController) NotFound(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	return c.HTML(http.StatusNotFound, views.PublicPage(c, &views.NotFound{}))
}

func (h *ErrorsController) HandleError(w http.ResponseWriter, r *http.Request, c *ctx.Ctx, err error) {
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
		if err := h.Forbidden(w, r, c); err != nil {
			h.HandleError(w, r, c, err)
		}
	case http.StatusNotFound:
		if err := h.NotFound(w, r, c); err != nil {
			h.HandleError(w, r, c, err)
		}
	default:
		c.Log.Error(err)
		http.Error(w, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
	}
}
