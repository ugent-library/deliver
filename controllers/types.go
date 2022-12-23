package controllers

import (
	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/models"
)

type (
	Ctx = *handler.Ctx[models.User]

	Data = map[string]any
)
