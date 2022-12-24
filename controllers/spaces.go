package controllers

import (
	"context"

	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Spaces struct {
	repo     models.RepositoryService
	listView view.View
	showView view.View
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repo:     r,
		listView: view.MustNew("page", "list_spaces"),
		showView: view.MustNew("page", "show_space"),
	}
}

type SpaceForm struct {
	Name string `form:"name"`
}

func (c *Spaces) List(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	spaces, err := c.repo.Spaces(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.Render(c.listView, Map{
		"spaces": spaces,
	})
}

func (c *Spaces) Show(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	spaceID := ctx.Path("spaceID")
	space, err := c.repo.Space(ctx.Req.Context(), spaceID)
	if err != nil {
		return err
	}
	return ctx.Render(c.showView, Map{
		"space": space,
	})
}

func (c *Spaces) Create(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	b := SpaceForm{}
	if err := bindForm(ctx.Req, &b); err != nil {
		return err
	}

	space := &models.Space{
		Name: b.Name,
	}
	if err := c.repo.CreateSpace(context.TODO(), space); err != nil {
		return err
	}

	ctx.PersistFlash(handler.Flash{
		Type: "info",
		Body: "Space created succesfully",
	})
	ctx.Redirect("space", "spaceID", space.ID)

	return nil
}
