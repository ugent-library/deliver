package controllers

import (
	"errors"
	"regexp"
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
	editView view.View
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repo:     r,
		showView: view.MustNew("page", "show_space"),
		newView:  view.MustNew("page", "new_space"),
		editView: view.MustNew("page", "edit_space"),
	}
}

var reSplitAdmins = regexp.MustCompile(`\s*[,;]\s*`)

type SpaceForm struct {
	Name   string `form:"name"`
	Admins string `form:"admins"`
}

func (h *Spaces) List(c *Ctx) error {
	var userSpaces []*models.Space
	var err error
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces(c.Context())
	} else {
		userSpaces, err = h.repo.SpacesByUsername(c.Context(), c.User.Username)
	}
	if err != nil {
		return err
	}

	if len(userSpaces) == 0 {
		return httperror.Forbidden
	}

	space, err := h.repo.SpaceByID(c.Context(), userSpaces[0].ID)
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

	// TODO add ToSpace() method to SpaceForm
	// or add a ultity BindSpace function?
	space := &models.Space{
		Name:   b.Name,
		Admins: reSplitAdmins.Split(b.Admins, -1),
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
	c.RedirectTo("space", "spaceName", space.Name)

	return nil
}

func (h *Spaces) CreateFolder(c *Ctx) error {
	spaceID := c.Path("spaceID")

	space, err := h.repo.SpaceByID(c.Context(), spaceID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, space) {
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

func (h *Spaces) Edit(c *Ctx) error {
	spaceID := c.Path("spaceID")

	space, err := h.repo.SpaceByID(c.Context(), spaceID)
	if err != nil {
		return err
	}

	return c.Render(h.editView, Map{
		"space":            space,
		"validationErrors": validate.NewErrors(),
	})
}

func (h *Spaces) Update(c *Ctx) error {
	spaceID := c.Path("spaceID")

	space, err := h.repo.SpaceByID(c.Context(), spaceID)
	if err != nil {
		return err
	}

	b := SpaceForm{}
	// TODO return ErrBadRequest
	if err := bind.Form(c.Req, &b); err != nil {
		return err
	}

	space.Admins = reSplitAdmins.Split(b.Admins, -1)

	if err := h.repo.UpdateSpace(c.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			return err
		}
		return c.Render(h.editView, Map{
			"space":            space,
			"validationErrors": validationErrors,
		})
	}

	c.RedirectTo("space", "spaceName", space.Name)

	return nil
}

func (h *Spaces) show(c *Ctx, folder *models.Folder, err error) error {
	spaceName := c.Path("spaceName")

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	space, err := h.repo.SpaceByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, space) {
		return httperror.Forbidden
	}

	var userSpaces []*models.Space
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces(c.Context())
	} else {
		userSpaces, err = h.repo.SpacesByUsername(c.Context(), c.User.Username)
	}
	if err != nil {
		return err
	}

	return c.Render(h.showView, Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
