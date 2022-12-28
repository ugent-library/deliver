package controllers

import (
	"errors"
	"net/http"

	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Errors struct {
	notFoundView view.View
}

func NewErrors() *Errors {
	return &Errors{
		notFoundView: view.MustNew("page", "not_found").Status(404),
	}
}

func (h *Errors) NotFound(c Ctx) error {
	return h.notFoundView.Render(c.Res, c)
}

func (h *Errors) HandleError(c Ctx, err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	switch httpErr.Code {
	case http.StatusUnauthorized:
		c.RedirectTo("login")
	case http.StatusNotFound:
		if err := h.NotFound(c); err != nil {
			h.HandleError(c, err)
		}
	default:
		c.Log.Error(err)
		http.Error(c.Res, http.StatusText(httpErr.Code), httpErr.Code)
	}
}
