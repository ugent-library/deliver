package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/htmx"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
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
		http.Redirect(w, r, c.Path("newSpace").String(), http.StatusSeeOther)
		return
	}

	// return forbidden if user is not an admin of anything
	if len(userSpaces) == 0 {
		c.HandleError(w, r, httperror.Forbidden)
		return
	}

	http.Redirect(w, r, c.Path("space", "spaceName", userSpaces[0].Name).String(), http.StatusSeeOther)
}

func ShowSpace(w http.ResponseWriter, r *http.Request) {
	showSpace(w, r, &models.Folder{}, nil)
}

func GetFolders(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	pagination := getPagination(r)
	folders, err := c.Repo.Folders.GetBySpace(r.Context(), space.ID, pagination)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	htmx.PushURL(w, c.Path("space", "spaceName", space.Name, pagination.ToPairs()).String())

	views.Folders(c, space, folders, pagination).Render(r.Context(), w)
}

func NewSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	views.NewSpace(c, &models.Space{}, okay.NewErrors()).Render(r.Context(), w)
}

func CreateSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	b := SpaceForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	// TODO add ToSpace() method to SpaceForm or add a utility BindSpace function?
	space := &models.Space{
		Name:   b.Name,
		Admins: reSplitAdmins.Split(b.Admins, -1),
	}

	if err := c.Repo.Spaces.Create(r.Context(), space); err != nil {
		validationErrors := okay.NewErrors()
		if !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}

		views.NewSpace(c, space, validationErrors).Render(r.Context(), w)
		return
	}

	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Space created successfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.Path("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func EditSpace(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)
	views.EditSpace(c, space, okay.NewErrors()).Render(r.Context(), w)
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
		validationErrors := okay.NewErrors()
		if !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}
		views.EditSpace(c, space, validationErrors).Render(r.Context(), w)
		return
	}

	loc := c.Path("space", "spaceName", space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func showSpace(w http.ResponseWriter, r *http.Request, folder *models.Folder, err error) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	newFolderArgs := views.NewFolderArgs{
		Folder: folder,
		Errors: okay.NewErrors(),
	}
	if err != nil && !errors.As(err, &newFolderArgs.Errors) {
		c.HandleError(w, r, err)
		return
	}

	if len(newFolderArgs.Errors.Errors) > 0 || r.URL.Query().Get("focus") == "new-folder" {
		newFolderArgs.Autofocus = true
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

	pagination := getPagination(r)
	folders, err := c.Repo.Folders.GetBySpace(r.Context(), space.ID, pagination)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	views.ShowSpace(c, space, folders, pagination, userSpaces, newFolderArgs).Render(r.Context(), w)
}

func getPagination(r *http.Request) *models.Pagination {
	query := r.URL.Query()
	filters := make([]models.Filter, 0, len(query))
	q := strings.TrimSpace(query.Get("q"))
	if q != "" {
		filters = append(filters, models.Filter{Name: "q", Value: q})
	}

	return models.NewPagination(getQueryParamAsInt(query, "offset"), getQueryParamAsInt(query, "limit"), query.Get("sort"), filters...)
}

func getQueryParamAsInt(query url.Values, paramName string) int {
	param := query.Get(paramName)

	intValue, err := strconv.Atoi(param)
	if err != nil {
		return -1
	}

	return intValue
}
