package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/bind"
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

func (h *Folders) Show(c Ctx) error {
	folderID := c.Path("folderID")
	folder, err := h.repo.Folder(c.Context(), folderID)
	if err != nil {
		return err
	}
	return c.Render(h.showView, Map{
		"folder": folder,
	})
}

func (h *Folders) Create(c Ctx) error {
	spaceID := c.Path("spaceID")
	b := FolderForm{}
	// TODO return ErrBadRequest
	if err := bind.Form(c.Req, &b); err != nil {
		return err
	}

	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}
	if err := h.repo.CreateFolder(c.Context(), folder); err != nil {
		return err
	}

	c.Session.Append(flashKey, Flash{
		Type: infoFlash,
		Body: "Folder created succesfully",
	})
	c.RedirectTo("folder", "folderID", folder.ID)

	return nil
}

// TODO remove files
func (h *Folders) Delete(c Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.Folder(c.Context(), folderID)
	if err != nil {
		return err
	}

	if err := h.repo.DeleteFolder(c.Context(), folderID); err != nil {
		return err
	}

	c.Session.Append(flashKey, Flash{
		Type: infoFlash,
		Body: "Folder deleted succesfully",
	})
	c.RedirectTo("space", "spaceID", folder.SpaceID)

	return nil
}

func (h *Folders) UploadFile(c Ctx) error {
	folderID := c.Path("folderID")

	// 2GB limit on request body
	c.Req.Body = http.MaxBytesReader(c.Res, c.Req.Body, 2_000_000_000)
	// buffer limit of 32MB
	if err := c.Req.ParseMultipartForm(32_000_000); err != nil {
		return err
	}

	for _, fileHeader := range c.Req.MultipartForm.File["file"] {
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
		md5, err := h.file.Add(c.Context(), file.ID, f)
		if err != nil {
			return err
		}

		file.Md5 = md5

		if err = h.repo.CreateFile(c.Context(), file); err != nil {
			return err
		}
	}

	c.Session.Append(flashKey, Flash{
		Type: infoFlash,
		Body: "File added succesfully",
	})

	c.RedirectTo("folder", "folderID", folderID)

	return nil
}
