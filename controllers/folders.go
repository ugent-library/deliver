package controllers

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Folders struct {
	repo     models.RepositoryService
	showView view.View
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService) *Folders {
	return &Folders{
		repo:     r,
		showView: view.MustNew("page", "show_folder"),
	}
}

func (c *Folders) Show(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	folderID := mux.Vars(r)["folderID"]
	c.showView.Render(w, ctx.Yield(Var{
		"folderID": folderID,
		"files":    nil,
	}))
}

func (c *Folders) Create(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	spaceID := mux.Vars(r)["spaceID"]
	b := FolderForm{}
	if err := bindForm(r, &b); err != nil {
		panic(err) // TODO
	}

	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}
	if err := c.repo.CreateFolder(context.TODO(), folder); err != nil {
		panic(err) // TODO
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "Folder created succesfully",
	})
	redirectURL := ctx.URLPath("folder", "folderID", folder.ID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (c *Folders) UploadFile(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	folderID := mux.Vars(r)["folderID"]

	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2_000_000_000)
	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32_000_000); err != nil {
		panic(err) // TODO
	}
	f, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err) // TODO
	}
	defer f.Close()

	// detect content type
	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		panic(err) // TODO
	}
	contentType := http.DetectContentType(buf)

	// and rewind
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err) // TODO
	}

	file := &models.File{
		FolderID:    folderID,
		Name:        fileHeader.Filename,
		ContentType: contentType,
	}

	if err = c.repo.CreateFile(context.TODO(), file, f); err != nil {
		panic(err) // TODO
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "File added succesfully",
	})
	redirectURL := ctx.URLPath("folder", "folderID", folderID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
