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
	"github.com/ugent-library/httpx"
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

func (h *SpacesController) List(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	var userSpaces []*models.Space
	var err error
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces.GetAll(r.Context())
	} else {
		userSpaces, err = h.repo.Spaces.GetAllByUsername(r.Context(), c.User.Username)
	}
	if err != nil {
		c.HandleError(err)
		return
	}

	// handle new empty installation
	if c.IsAdmin(c.User) && len(userSpaces) == 0 {
		c.PersistFlash(ctx.Flash{
			Type: "info",
			Body: "Create an initial space to get started",
		})
		http.Redirect(w, r, c.PathTo("newSpace").String(), http.StatusSeeOther)
		return
	}

	// return forbidden if user is not an admin of anything
	if len(userSpaces) == 0 {
		c.HandleError(httperror.Forbidden)
		return
	}

	space, err := h.repo.Spaces.Get(r.Context(), userSpaces[0].ID)
	if err != nil {
		c.HandleError(err)
		return
	}

	httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.ShowSpace{
		Space:            space,
		UserSpaces:       userSpaces,
		Folder:           &models.Folder{},
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Show(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	h.show(w, r, c, &models.Folder{}, nil)
}

func (h *SpacesController) New(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.NewSpace{
		Space:            &models.Space{},
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Create(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(errors.Join(httperror.BadRequest, err))
		return
	}

	// TODO add ToSpace() method to SpaceForm or add a ultity BindSpace function?
	space := &models.Space{
		Name:   b.Name,
		Admins: reSplitAdmins.Split(b.Admins, -1),
	}

	if err := h.repo.Spaces.Create(r.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			c.HandleError(err)
			return
		}
		httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.NewSpace{
			Space:            space,
			ValidationErrors: validationErrors,
		}))
		return
	}

	c.PersistFlash(ctx.Flash{
		Type:         "info",
		Body:         "Space created succesfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.PathTo("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func (h *SpacesController) CreateFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	spaceName := c.PathParam("spaceName")

	space, err := h.repo.Spaces.GetByName(r.Context(), spaceName)
	if err != nil {
		c.HandleError(err)
		return
	}

	if !c.IsSpaceAdmin(c.User, space) {
		c.HandleError(httperror.Forbidden)
		return
	}

	b := FolderForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(errors.Join(httperror.BadRequest, err))
		return
	}

	// TODO constructor for new objects
	folder := &models.Folder{
		SpaceID:   space.ID,
		Name:      b.Name,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	}

	if err := h.repo.Folders.Create(r.Context(), folder); err != nil {
		h.show(w, r, c, folder, err)
		return
	}

	c.PersistFlash(ctx.Flash{
		Type:         "info",
		Body:         "Folder created succesfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.PathTo("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func (h *SpacesController) Edit(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	spaceName := c.PathParam("spaceName")

	space, err := h.repo.Spaces.GetByName(r.Context(), spaceName)
	if err != nil {
		c.HandleError(err)
		return
	}

	httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.EditSpace{
		Space:            space,
		ValidationErrors: validate.NewErrors(),
	}))
}

func (h *SpacesController) Update(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	spaceName := c.PathParam("spaceName")

	space, err := h.repo.Spaces.GetByName(r.Context(), spaceName)
	if err != nil {
		c.HandleError(err)
		return
	}

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(errors.Join(httperror.BadRequest, err))
		return
	}

	space.Admins = reSplitAdmins.Split(b.Admins, -1)

	if err := h.repo.Spaces.Update(r.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			c.HandleError(err)
			return
		}
		httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.EditSpace{
			Space:            space,
			ValidationErrors: validationErrors,
		}))
		return
	}

	loc := c.PathTo("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func (h *SpacesController) show(w http.ResponseWriter, r *http.Request, c *ctx.Ctx, folder *models.Folder, err error) {
	spaceName := c.PathParam("spaceName")

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		c.HandleError(err)
		return
	}

	space, err := h.repo.Spaces.GetByName(r.Context(), spaceName)
	if err != nil {
		c.HandleError(err)
		return
	}

	if !c.IsSpaceAdmin(c.User, space) {
		c.HandleError(httperror.Forbidden)
		return
	}

	var userSpaces []*models.Space
	if c.IsAdmin(c.User) {
		userSpaces, err = h.repo.Spaces.GetAll(r.Context())
	} else {
		userSpaces, err = h.repo.Spaces.GetAllByUsername(r.Context(), c.User.Username)
	}
	if err != nil {
		c.HandleError(err)
		return
	}

	httpx.RenderHTML(w, http.StatusOK, views.Page(c, &views.ShowSpace{
		Space:            space,
		UserSpaces:       userSpaces,
		Folder:           folder,
		ValidationErrors: validationErrors,
	}))
}
