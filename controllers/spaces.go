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

func (h *Spaces) List(c Ctx) error {
	spaces, err := h.repo.Spaces(c.Context())
	if err != nil {
		return err
	}
	return c.Render(h.listView, Map{
		"spaces": spaces,
	})
}

func (h *Spaces) Show(c Ctx) error {
	spaceID := c.Path("spaceID")
	space, err := h.repo.Space(c.Req.Context(), spaceID)
	if err != nil {
		return err
	}
	return c.Render(h.showView, Map{
		"space": space,
	})
}

func (h *Spaces) Create(c Ctx) error {
	b := SpaceForm{}
	if err := bindForm(c.Req, &b); err != nil {
		return err
	}

	space := &models.Space{
		Name: b.Name,
	}
	if err := h.repo.CreateSpace(context.TODO(), space); err != nil {
		return err
	}

	c.PersistFlash(Flash{
		Type: "info",
		Body: "Space created succesfully",
	})
	c.Redirect("space", "spaceID", space.ID)

	return nil
}
