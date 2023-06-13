package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/httpx/render"
)

type FolderForm struct {
	Name string `form:"name"`
}

func ShowFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	if htmx.Request(r) {
		render.HTML(w, http.StatusOK, views.Files(c, folder.Files))
		return
	}

	render.HTML(w, http.StatusOK, views.Page(c, &views.ShowFolder{
		Folder:      folder,
		MaxFileSize: c.MaxFileSize,
	}))
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	b := FolderForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	// TODO constructor for new objects
	folder := &models.Folder{
		SpaceID:   space.ID,
		Name:      b.Name,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	}

	if err := c.Repo.Folders.Create(r.Context(), folder); err != nil {
		showSpace(w, r, folder, err)
		return
	}

	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Folder created succesfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.PathTo("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func EditFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	render.HTML(w, http.StatusOK, views.Page(c, &views.EditFolder{
		Folder:           folder,
		ValidationErrors: validate.NewErrors(),
	}))
}

func UpdateFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	b := FolderForm{}
	if err := bind.Form(r, &b); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	folder.Name = b.Name

	if err := c.Repo.Folders.Update(r.Context(), folder); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}

		render.HTML(w, http.StatusOK, views.Page(c, &views.EditFolder{
			Folder:           folder,
			ValidationErrors: validationErrors,
		}))
		return
	}

	loc := c.PathTo("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	if err := c.Repo.Folders.Delete(r.Context(), folder.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	c.Hub.Send("space."+folder.Space.ID, views.AddFlash(ctx.Flash{
		Type: "info",
		Body: fmt.Sprintf("%s just deleted the folder %s.", c.User.Name, folder.Name),
	}))
	c.Hub.Send("folder."+folder.ID, views.AddFlash(ctx.Flash{
		Type: "error",
		Body: fmt.Sprintf("%s just deleted this folder.", c.User.Name),
	}))
	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Folder deleted succesfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.PathTo("space", "spaceName", folder.Space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func ShareFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	render.HTML(w, http.StatusOK, views.PublicPage(c, &views.ShareFolder{
		Folder: folder,
	}))
}
