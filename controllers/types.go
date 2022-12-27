package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
)

type (
	Ctx = *handler.Ctx[models.User, handler.Unused, Flash]
	Map = map[string]any
)

type Flash struct {
	Type  string
	Title string
	Body  template.HTML
}

func RequireUser(c Ctx) error {
	if c.User() == nil {
		return httperror.Unauthorized
	}
	return nil
}

func HandleError(c Ctx, err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	switch httpErr.Code {
	case http.StatusUnauthorized:
		c.Redirect("login")
	case http.StatusNotFound:
		if err := NotFound(c); err != nil {
			HandleError(c, err)
		}
	default:
		http.Error(c.Res, http.StatusText(httpErr.Code), httpErr.Code)
	}
}
