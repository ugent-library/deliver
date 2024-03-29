package handlers

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/htmx"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/okay"
)

type FolderForm struct {
	Name string `form:"name"`
}

func ShowFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	if htmx.Request(r) {
		views.Files(c, folder.Files).Render(r.Context(), w)
		return
	}

	views.ShowFolder(c, folder).Render(r.Context(), w)
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	space := ctx.GetSpace(r)

	b := FolderForm{}
	if err := bind.Form(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	folder := models.NewFolder(space.ID, b.Name)

	if err := c.Repo.Folders.Create(r.Context(), folder); err != nil {
		showSpace(w, r, folder, err)
		return
	}

	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Folder created successfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.Path("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func EditFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)
	views.EditFolder(c, folder, okay.NewErrors()).Render(r.Context(), w)
}

func UpdateFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	b := FolderForm{}
	if err := bind.Form(r, &b, bind.Vacuum); err != nil {
		c.HandleError(w, r, errors.Join(httperror.BadRequest, err))
		return
	}

	folder.Name = b.Name

	if err := c.Repo.Folders.Update(r.Context(), folder); err != nil {
		validationErrors := okay.NewErrors()
		if !errors.As(err, &validationErrors) {
			c.HandleError(w, r, err)
			return
		}

		views.EditFolder(c, folder, validationErrors).Render(r.Context(), w)
		return
	}

	loc := c.Path("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func PostponeFolderExpiration(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	folder.PostponeExpiration()

	if err := c.Repo.Folders.Update(r.Context(), folder); err != nil {
		c.PersistFlash(w, ctx.Flash{
			Type: "error",
			Body: fmt.Sprintf("Unexpected error: %s", err.Error()),
		})
	} else {
		c.PersistFlash(w, ctx.Flash{
			Type:         "success",
			Body:         fmt.Sprintf("New expiration date for %s: %s", folder.Name, folder.ExpiresAt.In(c.Timezone).Format("2006-01-02")),
			DismissAfter: 3 * time.Second,
		})
	}

	htmx.Refresh(w)
}

func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	if err := c.Repo.Folders.Delete(r.Context(), folder.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	c.Hub.SendString("space."+folder.Space.ID, views.String(views.AddFlash(ctx.Flash{
		Type: "info",
		Body: fmt.Sprintf("%s just deleted the folder %s.", c.User.Name, folder.Name),
	})))
	c.Hub.SendString("folder."+folder.ID, views.String(views.AddFlash(ctx.Flash{
		Type: "error",
		Body: fmt.Sprintf("%s just deleted this folder.", c.User.Name),
	})))
	c.PersistFlash(w, ctx.Flash{
		Type:         "info",
		Body:         "Folder deleted successfully",
		DismissAfter: 3 * time.Second,
	})

	loc := c.Path("space", "spaceName", folder.Space.Name).String()
	http.Redirect(w, r, loc, http.StatusSeeOther)
}

func ShareFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	views.ShareFolder(c, folder).Render(r.Context(), w)
}

func DownloadFolder(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s-%s.zip", folder.ID, folder.Slug()),
	)

	zipWriter := zip.NewWriter(bufio.NewWriter(w))
	defer zipWriter.Close()

	for _, file := range folder.Files {
		fileWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:     file.Name,
			Method:   zip.Store, // no compression for streaming
			Modified: time.Now(),
		})
		if err != nil {
			c.Log.Error(err)
			return
		}

		if err := c.Repo.Files.AddDownload(r.Context(), file.ID); err != nil {
			c.Log.Error(err)
		}

		fileReader, err := c.Storage.Get(r.Context(), file.ID)
		if err != nil {
			c.Log.Error(err)
			return
		}
		io.Copy(fileWriter, fileReader)
	}
}
