package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ugent-library/deliver/bind"
	"github.com/ugent-library/deliver/httperror"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/ulid"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/view"
)

type Folders struct {
	repo     models.RepositoryService
	file     models.FileService
	showView view.View
	editView view.View
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService, f models.FileService) *Folders {
	return &Folders{
		repo:     r,
		file:     f,
		showView: view.MustNew("page", "show_folder"),
		editView: view.MustNew("page", "edit_folder"),
	}
}

func (h *Folders) Show(c *Ctx) error {
	return h.show(c, nil)
}

func (h *Folders) Edit(c *Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	return c.Render(h.editView, Map{
		"folder":           folder,
		"validationErrors": validate.NewErrors(),
	})
}

func (h *Folders) Update(c *Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	b := FolderForm{}
	// TODO return ErrBadRequest
	if err := bind.Form(c.Req, &b); err != nil {
		return err
	}

	folder.Name = b.Name

	if err := h.repo.UpdateFolder(c.Context(), folder); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			return err
		}
		return c.Render(h.editView, Map{
			"folder":           folder,
			"validationErrors": validationErrors,
		})
	}

	c.RedirectTo("folder", "folderID", folder.ID)

	return nil
}

func (h *Folders) Delete(c *Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	for _, f := range folder.Files {
		if err := h.file.Delete(c.Context(), f.ID); err != nil {
			return err
		}
	}
	if err := h.repo.DeleteFolder(c.Context(), folderID); err != nil {
		return err
	}

	c.Session.Append(flashKey, Flash{
		Type:         infoFlash,
		Body:         "Folder deleted succesfully",
		DismissAfter: 3 * time.Second,
	})
	c.RedirectTo("space", "spaceName", folder.Space.Name)

	return nil
}

func (h *Folders) UploadFile(c *Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

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

		file.MD5 = md5

		if err = h.repo.CreateFile(c.Context(), file); err != nil {
			return h.show(c, err)
		}
	}

	c.Session.Append(flashKey, Flash{
		Type:         infoFlash,
		Body:         "File added succesfully",
		DismissAfter: 3 * time.Second,
	})

	c.RedirectTo("folder", "folderID", folderID)

	return nil
}

func (h *Folders) show(c *Ctx, err error) error {
	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	folderID := c.Path("folderID")
	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}
	return c.Render(h.showView, Map{
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
