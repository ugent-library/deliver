package controllers

import (
	"errors"

	"github.com/ugent-library/dilliver/bind"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/validate"
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

func (h *Spaces) List(c *Ctx) error {
	return h.list(c, nil)
}

func (h *Spaces) Show(c *Ctx) error {
	return h.show(c, nil)
}

func (h *Spaces) Create(c *Ctx) error {
	b := SpaceForm{}
	// TODO return ErrBadRequest
	if err := bind.Form(c.Req, &b); err != nil {
		return err
	}

	space := &models.Space{
		Name: b.Name,
	}

	if err := h.repo.CreateSpace(c.Context(), space); err != nil {
		return h.list(c, err)
	}

	c.Session.Append(flashKey, Flash{
		Type: infoFlash,
		Body: "Space created succesfully",
	})
	c.RedirectTo("space", "spaceID", space.ID)

	return nil
}

func (h *Spaces) CreateFolder(c *Ctx) error {
	spaceID := c.Path("spaceID")

	if !c.IsSpaceAdmin(spaceID, c.User) {
		return httperror.Forbidden
	}

	b := FolderForm{}
	// TODO return ErrBadRequest
	if err := bind.Form(c.Req, &b); err != nil {
		return err
	}

	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}

	if err := h.repo.CreateFolder(c.Context(), folder); err != nil {
		return h.show(c, err)
	}

	c.Session.Append(flashKey, Flash{
		Type: infoFlash,
		Body: "Folder created succesfully",
	})
	c.RedirectTo("folder", "folderID", folder.ID)

	return nil
}

func (h *Spaces) list(c *Ctx, err error) error {
	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	spaces, err := h.repo.Spaces(c.Context())
	if err != nil {
		return err
	}

	return c.Render(h.listView, Map{
		"spaces":           spaces,
		"validationErrors": validationErrors,
	})
}

func (h *Spaces) show(c *Ctx, err error) error {
	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	spaceID := c.Path("spaceID")
	space, err := h.repo.Space(c.Context(), spaceID)
	if err != nil {
		return err
	}
	return c.Render(h.showView, Map{
		"space":            space,
		"validationErrors": validationErrors,
	})
}
