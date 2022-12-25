package controllers

import (
	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
)

type (
	Ctx = *handler.Ctx[models.User, handler.Empty]
	Map = map[string]any
)

func RequireUser(c Ctx) error {
	if c.User() == nil {
		return httperror.Unauthorized
	}
	return nil
}
