package controllers

import (
	"context"

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

func (c *Spaces) List(ctx *Ctx) error {
	if ctx.User() == nil {
		return ErrUnauthorized
	}

	spaces, err := c.repo.Spaces(ctx.Context())
	if err != nil {
		return err
	}
	return c.listView.Render(ctx.Res, ctx.Yield(Var{
		"spaces": spaces,
	}))
}

func (c *Spaces) Show(ctx *Ctx) error {
	if ctx.User() == nil {
		return ErrUnauthorized
	}

	spaceID := ctx.Path("spaceID")
	space, err := c.repo.Space(ctx.Req.Context(), spaceID)
	if err != nil {
		return err
	}
	return c.showView.Render(ctx.Res, ctx.Yield(Var{
		"space": space,
	}))
}

func (c *Spaces) Create(ctx *Ctx) error {
	if ctx.User() == nil {
		return ErrUnauthorized
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

	ctx.PersistFlash(Flash{
		Type: "info",
		Body: "Space created succesfully",
	})
	ctx.Redirect("space", "spaceID", space.ID)

	return nil
}
