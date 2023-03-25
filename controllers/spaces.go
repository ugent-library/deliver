package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httperror"
)

type Spaces struct {
	repo models.RepositoryService
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repo: r,
	}
}

var reSplitAdmins = regexp.MustCompile(`\s*[,;]\s*`)

type SpaceForm struct {
	Name   string `form:"name"`
	Admins string `form:"admins"`
}

func (h *Spaces) List(c *ctx.Ctx) error {
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

	// handle new empty installation
	if c.IsAdmin(c.User) && len(userSpaces) == 0 {
		c.AddFlash(ctx.Flash{
			Type: "info",
			Body: "Create an initial space to get started",
		})
		c.RedirectTo("new_space")
		return nil
	}

	// return forbidden if user is not an admin of anything
	if len(userSpaces) == 0 {
		return httperror.Forbidden
	}

	space, err := h.repo.SpaceByID(c.Context(), userSpaces[0].ID)
	if err != nil {
		return err
	}

	return c.HTMLX(http.StatusOK, "layouts/page", "show_space", Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           &models.Folder{},
		"validationErrors": validate.NewErrors(),
	})
}

func (h *Spaces) Show(c *ctx.Ctx) error {
	return h.show(c, &models.Folder{}, nil)
}

func (h *Spaces) New(c *ctx.Ctx) error {
	return c.HTML(http.StatusOK, views.Page(c, &views.NewSpace{
		Space:            &models.Space{},
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *Spaces) Create(c *ctx.Ctx) error {
	b := SpaceForm{}
	if err := bind.Form(c.Req, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
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
		return c.HTML(http.StatusOK, views.Page(c, &views.NewSpace{
			Space:            space,
			ValidationErrors: validationErrors,
		}))
	}

	c.AddFlash(ctx.Flash{
		Type:         "info",
		Body:         "Space created succesfully",
		DismissAfter: 3 * time.Second,
	})
	c.RedirectTo("space", "spaceName", space.Name)
	return nil
}

func (h *Spaces) CreateFolder(c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.SpaceByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, space) {
		return httperror.Forbidden
	}

	b := FolderForm{}
	if err := bind.Form(c.Req, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	// TODO constructor for new objects
	folder := &models.Folder{
		SpaceID:   space.ID,
		Name:      b.Name,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	}

	if err := h.repo.CreateFolder(c.Context(), folder); err != nil {
		return h.show(c, folder, err)
	}

	c.AddFlash(ctx.Flash{
		Type:         "info",
		Body:         "Folder created succesfully",
		DismissAfter: 3 * time.Second,
	})
	c.RedirectTo("folder", "folderID", folder.ID)

	return nil
}

func (h *Spaces) Edit(c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.SpaceByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, views.Page(c, &views.EditSpace{
		Space:            space,
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *Spaces) Update(c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.SpaceByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	b := SpaceForm{}
	if err := bind.Form(c.Req, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	space.Admins = reSplitAdmins.Split(b.Admins, -1)

	if err := h.repo.UpdateSpace(c.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			return err
		}
		return c.HTML(http.StatusOK, views.Page(c, &views.EditSpace{
			Space:            space,
			ValidationErrors: validationErrors,
		}))
	}

	c.RedirectTo("space", "spaceName", space.Name)

	return nil
}

func (h *Spaces) show(c *ctx.Ctx, folder *models.Folder, err error) error {
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

	return c.HTMLX(http.StatusOK, "layouts/page", "show_space", Map{
		"space":            space,
		"userSpaces":       userSpaces,
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
