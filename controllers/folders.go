package controllers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/ulid"
	"github.com/ugent-library/dilliver/view"
)

type Folders struct {
	repo     models.RepositoryService
	file     models.FileService
	showView view.View
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService, f models.FileService) *Folders {
	return &Folders{
		repo:     r,
		file:     f,
		showView: view.MustNew("page", "show_folder"),
	}
}

func (c *Folders) Show(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	folderID := mux.Vars(r)["folderID"]
	folder, err := c.repo.Folder(context.TODO(), folderID)
	if errors.Is(err, models.ErrNotFound) {
		ctx.Router.NotFoundHandler.ServeHTTP(w, r)
		return
	}
	if err != nil {
		panic(err) // TODO
	}
	c.showView.Render(w, ctx.Yield(Var{
		"folder": folder,
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

// TODO remove files
func (c *Folders) Delete(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	folderID := mux.Vars(r)["folderID"]

	folder, err := c.repo.Folder(context.TODO(), folderID)
	if errors.Is(err, models.ErrNotFound) {
		ctx.Router.NotFoundHandler.ServeHTTP(w, r)
		return
	}
	if err != nil {
		panic(err) // TODO
	}

	if err := c.repo.DeleteFolder(context.TODO(), folderID); err != nil {
		panic(err) // TODO
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "Folder deleted succesfully",
	})
	redirectURL := ctx.URLPath("space", "spaceID", folder.SpaceID).String()
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

	for _, fileHeader := range r.MultipartForm.File["file"] {
		f, err := fileHeader.Open()
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
			ID:          ulid.MustGenerate(),
			Name:        fileHeader.Filename,
			ContentType: contentType,
			Size:        fileHeader.Size,
		}

		// TODO get size
		md5, err := c.file.Add(context.TODO(), file.ID, f)
		if err != nil {
			panic(err) // TODO
		}

		file.Md5 = md5

		if err = c.repo.CreateFile(context.TODO(), file); err != nil {
			panic(err) // TODO
		}
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "File added succesfully",
	})
	redirectURL := ctx.URLPath("folder", "folderID", folderID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
