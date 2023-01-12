package controllers

import (
	"errors"
	"time"

	"github.com/ugent-library/deliver/bind"
	"github.com/ugent-library/deliver/httperror"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/view"
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
	return h.show(c, &models.Folder{}, nil)
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
		Type:         infoFlash,
		Body:         "Space created succesfully",
		DismissAfter: 3 * time.Second,
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
		SpaceID:   spaceID,
		Name:      b.Name,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	}

	if err := h.repo.CreateFolder(c.Context(), folder); err != nil {
		return h.show(c, folder, err)
	}

	c.Session.Append(flashKey, Flash{
		Type:         infoFlash,
		Body:         "Folder created succesfully",
		DismissAfter: 3 * time.Second,
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

func (h *Spaces) show(c *Ctx, folder *models.Folder, err error) error {
	spaceID := c.Path("spaceID")

	if !c.IsSpaceAdmin(spaceID, c.User) {
		return httperror.Forbidden
	}

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	// TODO clean this up
	// TODO don't eager load folders for all user spaces
	// TODO get all spaces in 1 query (admin sees all spaces anyway)
	var space *models.Space
	var userSpaces []*models.Space
	if c.IsAdmin(c.User) {
		allSpaces, err := h.repo.Spaces(c.Context())
		if err != nil {
			return err
		}
		userSpaces = allSpaces
	} else {
		userSpaceIDs := c.UserSpaces(c.User)
		userSpaces = make([]*models.Space, len(userSpaceIDs))
		for i, id := range userSpaceIDs {
			s, err := h.repo.Space(c.Context(), id)
			if err != nil {
				return err
			}
			userSpaces[i] = s
			if id == spaceID {
				space = s
			}
		}
	}
	for _, s := range userSpaces {
		if s.ID == spaceID {
			space = s
			break
		}
	}

	return c.Render(h.showView, Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
