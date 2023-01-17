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
	showView view.View
	newView  view.View
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repo:     r,
		showView: view.MustNew("page", "show_space"),
		newView:  view.MustNew("page", "new_space"),
	}
}

type SpaceForm struct {
	Name string `form:"name"`
}

// TODO clean this up
func (h *Spaces) List(c *Ctx) error {
	var userSpaces []*models.Space
	allSpaces, err := h.repo.Spaces(c.Context())
	if err != nil {
		return err
	}
	if c.IsAdmin(c.User) {
		userSpaces = allSpaces
	} else {
		userSpaceIDs := c.UserSpaces(c.User)
		userSpaces = make([]*models.Space, len(userSpaceIDs))
		for i, id := range userSpaceIDs {
			for _, s := range allSpaces {
				if s.ID == id {
					userSpaces[i] = s
					break
				}
			}
		}
	}

	if len(userSpaces) == 0 {
		return httperror.Forbidden
	}

	space, err := h.repo.Space(c.Context(), userSpaces[0].ID)
	if err != nil {
		return err
	}

	return c.Render(h.showView, Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           &models.Folder{},
		"validationErrors": validate.NewErrors(),
	})
}

func (h *Spaces) Show(c *Ctx) error {
	return h.show(c, &models.Folder{}, nil)
}

func (h *Spaces) New(c *Ctx) error {
	return c.Render(h.newView, Map{
		"space":            &models.Space{},
		"validationErrors": validate.NewErrors(),
	})
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
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			return err
		}
		return c.Render(h.newView, Map{
			"space":            space,
			"validationErrors": validationErrors,
		})
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

// TODO clean this up
func (h *Spaces) show(c *Ctx, folder *models.Folder, err error) error {
	spaceID := c.Path("spaceID")

	if !c.IsSpaceAdmin(spaceID, c.User) {
		return httperror.Forbidden
	}

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	space, err := h.repo.Space(c.Context(), spaceID)
	if err != nil {
		return err
	}
	var userSpaces []*models.Space
	allSpaces, err := h.repo.Spaces(c.Context())
	if err != nil {
		return err
	}
	if c.IsAdmin(c.User) {
		userSpaces = allSpaces
	} else {
		userSpaceIDs := c.UserSpaces(c.User)
		userSpaces = make([]*models.Space, len(userSpaceIDs))
		for i, id := range userSpaceIDs {
			for _, s := range allSpaces {
				if s.ID == id {
					userSpaces[i] = s
					break
				}
			}
		}
	}

	return c.Render(h.showView, Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
