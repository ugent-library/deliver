package controllers

import (
	"html/template"

	"github.com/ugent-library/dilliver/handler"
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
