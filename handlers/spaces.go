package handlers

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

var reSplitAdmins = regexp.MustCompile(`\s*[,;]\s*`)

type SpaceForm struct {
	Name   string `form:"name"`
	Admins string `form:"admins"`
}

func ListSpaces(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	var userSpaces []*models.Space
	var err error
	if c.Permissions.IsAdmin(c.User) {
		userSpaces, err = c.Repo.Spaces.GetAll(r.Context())
	} else {
		userSpaces, err = c.Repo.Spaces.GetAllByUsername(r.Context(), c.User.Username)
	}
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	// handle new empty installation
	if c.Permissions.IsAdmin(c.User) && len(userSpaces) == 0 {
		c.PersistFlash(w, ctx.Flash{
			Type: "info",
			Body: "Create an initial space to get started",
		})
		http.Redirect(w, r, c.PathTo("newSpace").String(), http.StatusSeeOther)
		return
	}

	// return forbidden if user is not an admin of anything
	if len(userSpaces) == 0 {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	space, err := c.Repo.Spaces.GetByName(r.Context(), userSpaces[0].Name)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowSpace(c, space, userSpaces, &models.Folder{}, validate.NewErrors()).Render(r.Context(), w)
}

func ShowSpace(w http.ResponseWriter, r *http.Request) {
	showSpace(w, r, &models.Folder{}, nil)
}

func NewSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.NewSpace(c, &models.Space{}, validate.NewErrors()).Render(r.Context(), w)
}

func CreateSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	// TODO add ToSpace() method to SpaceForm or add a ultity BindSpace function?
	space := &models.Space{
		Name:   b.Name,
		Admins: reSplitAdmins.Split(b.Admins, -1),
	}

	if err := c.Repo.Spaces.Create(r.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}

		views.NewSpace(c, space, validationErrors).Render(r.Context(), w)
		return
	}

	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Space created succesfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.PathTo("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func EditSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)
	views.EditSpace(c, space, validate.NewErrors()).Render(r.Context(), w)
}

func UpdateSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	space.Admins = reSplitAdmins.Split(b.Admins, -1)

	if err := c.Repo.Spaces.Update(r.Context(), space); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}
		views.EditSpace(c, space, validationErrors).Render(r.Context(), w)
		return
	}

	loc := c.PathTo("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func showSpace(w http.ResponseWriter, r *http.Request, folder *models.Folder, err error) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		c.HandleError(w, r, err)
		return
	}

	var userSpaces []*models.Space
	if c.Permissions.IsAdmin(c.User) {
		userSpaces, err = c.Repo.Spaces.GetAll(r.Context())
	} else {
		userSpaces, err = c.Repo.Spaces.GetAllByUsername(r.Context(), c.User.Username)
	}
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowSpace(c, space, userSpaces, folder, validationErrors).Render(r.Context(), w)
}