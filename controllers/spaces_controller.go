package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httperror"
)

type SpacesController struct {
	repo *repositories.Repo
}

func NewSpacesController(r *repositories.Repo) *SpacesController {
	return &SpacesController{
		repo: r,
	}
}

var reSplitAdmins = regexp.MustCompile(`\s*[,;]\s*`)

type SpaceForm struct {
	Name   string `form:"name"`
	Admins string `form:"admins"`
}

func (h *SpacesController) List(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	var userSpaces []*models.Space
	var err error
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces.GetAll(c.Context())
	} else {
		userSpaces, err = h.repo.Spaces.GetAllByUsername(c.Context(), c.User.Username)
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
		c.RedirectTo("newSpace")
		return nil
	}

	// return forbidden if user is not an admin of anything
	if len(userSpaces) == 0 {
		return httperror.Forbidden
	}

	space, err := h.repo.Spaces.Get(c.Context(), userSpaces[0].ID)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, views.Page(c, &views.ShowSpace{
		Space:            space,
		UserSpaces:       userSpaces,
		Folder:           &models.Folder{},
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Show(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	return h.show(c, &models.Folder{}, nil)
}

func (h *SpacesController) New(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	return c.HTML(http.StatusOK, views.Page(c, &views.NewSpace{
		Space:            &models.Space{},
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Create(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	// TODO add ToSpace() method to SpaceForm
	// or add a ultity BindSpace function?
	space := &models.Space{
		Name:   b.Name,
		Admins: reSplitAdmins.Split(b.Admins, -1),
	}

	if err := h.repo.Spaces.Create(c.Context(), space); err != nil {
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

func (h *SpacesController) CreateFolder(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.Spaces.GetByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, space) {
		return httperror.Forbidden
	}

	b := FolderForm{}
	if err := bind.Form(r, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	// TODO constructor for new objects
	folder := &models.Folder{
		SpaceID:   space.ID,
		Name:      b.Name,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	}

	if err := h.repo.Folders.Create(c.Context(), folder); err != nil {
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

func (h *SpacesController) Edit(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.Spaces.GetByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, views.Page(c, &views.EditSpace{
		Space:            space,
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Update(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	spaceName := c.Path("spaceName")

	space, err := h.repo.Spaces.GetByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	space.Admins = reSplitAdmins.Split(b.Admins, -1)

	if err := h.repo.Spaces.Update(c.Context(), space); err != nil {
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

func (h *SpacesController) show(c *ctx.Ctx, folder *models.Folder, err error) error {
	spaceName := c.Path("spaceName")

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	space, err := h.repo.Spaces.GetByName(c.Context(), spaceName)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, space) {
		return httperror.Forbidden
	}

	var userSpaces []*models.Space
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces.GetAll(c.Context())
	} else {
		userSpaces, err = h.repo.Spaces.GetAllByUsername(c.Context(), c.User.Username)
	}
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, views.Page(c, &views.ShowSpace{
		Space:            space,
		UserSpaces:       userSpaces,
		Folder:           folder,
		ValidationErrors: validationErrors,
	}))
}
