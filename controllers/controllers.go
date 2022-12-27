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

	flashTypeInfo = "info"
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
	if user := c.Session.Get(userKey); user != nil {
		c.Var.User = user.(*models.User)
	}
	if flash := c.Session.Pop(flashKey); flash != nil {
		c.Var.Flash = flash.([]Flash)
	}
	return c.Session.Save()
}

func RequireUser(c Ctx) error {
	if c.Var.User == nil {
		return httperror.Unauthorized
	}
	return nil
}
