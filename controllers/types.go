package controllers

import (
	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/models"
)

type (
	Ctx = *handler.Ctx[models.User]

	Var = map[string]any
)
