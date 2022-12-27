package controllers

import (
	"github.com/ugent-library/dilliver/view"
)

var (
	HomeView     = view.MustNew("page", "home")
	NotFoundView = view.MustNew("page", "not_found").Status(404)
)

func Home(c Ctx) error {
	return HomeView.Render(c.Res, c)
}

func NotFound(c Ctx) error {
	return NotFoundView.Render(c.Res, c)
}
