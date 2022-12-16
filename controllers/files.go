package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/models"
)

type Files struct {
	repo models.RepositoryService
	file models.FileService
}

func NewFiles(r models.RepositoryService, f models.FileService) *Files {
	return &Files{
		repo: r,
		file: f,
	}
}

func (c *Files) Download(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	fileID := mux.Vars(r)["fileID"]

	_, err := c.repo.File(context.TODO(), fileID)
	if errors.Is(err, models.ErrNotFound) {
		ctx.Router.NotFoundHandler.ServeHTTP(w, r)
		return
	}
	if err != nil {
		panic(err) // TODO
	}
	c.file.Get(context.TODO(), fileID, w)
}

func (c *Files) Delete(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	fileID := mux.Vars(r)["fileID"]

	file, err := c.repo.File(context.TODO(), fileID)
	if errors.Is(err, models.ErrNotFound) {
		ctx.Router.NotFoundHandler.ServeHTTP(w, r)
		return
	}
	if err != nil {
		panic(err) // TODO
	}

	if err := c.repo.DeleteFile(context.TODO(), fileID); err != nil {
		panic(err) // TODO
	}

	if err := c.file.Delete(context.TODO(), fileID); err != nil {
		panic(err) // TODO
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "File deleted succesfully",
	})
	redirectURL := ctx.URLPath("folder", "folderID", file.FolderID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
