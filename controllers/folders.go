package controllers

import (
	"net/http"

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

func (c *Folders) Show(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	folderID := ctx.Path("folderID")
	folder, err := c.repo.Folder(r.Context(), folderID)
	if err != nil {
		return err
	}
	return c.showView.Render(w, ctx.Yield(Var{
		"folder": folder,
	}))
}

func (c *Folders) Create(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	spaceID := ctx.Path("spaceID")
	b := FolderForm{}
	if err := bindForm(r, &b); err != nil {
		return err
	}

	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}
	if err := c.repo.CreateFolder(r.Context(), folder); err != nil {
		return err
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "Folder created succesfully",
	})
	ctx.Redirect(w, r, "folder", "folderID", folder.ID)

	return nil
}

// TODO remove files
func (c *Folders) Delete(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	folderID := ctx.Path("folderID")

	folder, err := c.repo.Folder(r.Context(), folderID)
	if err != nil {
		return err
	}

	if err := c.repo.DeleteFolder(r.Context(), folderID); err != nil {
		return err
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "Folder deleted succesfully",
	})
	ctx.Redirect(w, r, "space", "spaceID", folder.SpaceID)

	return nil
}

func (c *Folders) UploadFile(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	folderID := ctx.Path("folderID")

	// 2GB limit on request body
	r.Body = http.MaxBytesReader(w, r.Body, 2_000_000_000)
	// buffer limit of 32MB
	if err := r.ParseMultipartForm(32_000_000); err != nil {
		return err
	}

	for _, fileHeader := range r.MultipartForm.File["file"] {
		f, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		contentType, err := detectContentType(f)
		if err != nil {
			return err
		}

		file := &models.File{
			FolderID:    folderID,
			ID:          ulid.MustGenerate(),
			Name:        fileHeader.Filename,
			ContentType: contentType,
			Size:        fileHeader.Size,
		}

		// TODO get size
		md5, err := c.file.Add(r.Context(), file.ID, f)
		if err != nil {
			return err
		}

		file.Md5 = md5

		if err = c.repo.CreateFile(r.Context(), file); err != nil {
			return err
		}
	}

	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "File added succesfully",
	})
	ctx.Redirect(w, r, "folder", "folderID", folderID)

	return nil
}
