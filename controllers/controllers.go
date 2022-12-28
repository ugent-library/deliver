package controllers

import (
	"html/template"

	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
)

const (
	userKey  = "user"
	flashKey = "flash"

	infoFlash = "info"
)

type (
	Var struct {
		User  *models.User
		Flash []Flash
	}

	Flash struct {
		Type  string
		Title string
		Body  template.HTML
	}

	Ctx = *handler.Ctx[Var]
	Map = map[string]any
)

func LoadSession(c Ctx) error {
	if val := c.Session.Get(userKey); val != nil {
		c.Var.User = val.(*models.User)
	}
	if vals := c.Session.Pop(flashKey); vals != nil {
		flash := make([]Flash, len(vals.([]any)))
		for i, v := range vals.([]any) {
			flash[i] = v.(Flash)
		}
		c.Var.Flash = flash
	}
	return nil
}

func RequireUser(c Ctx) error {
	if c.Var.User == nil {
		return httperror.Unauthorized
	}
	return nil
}
