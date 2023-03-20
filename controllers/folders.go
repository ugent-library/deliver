package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/controllers/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/turbo"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/deliver/views"
	"github.com/ugent-library/httperror"
)

type Folders struct {
	repo        models.RepositoryService
	file        models.FileService
	maxFileSize int64
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService, f models.FileService, maxFileSize int64) *Folders {
	return &Folders{
		repo:        r,
		file:        f,
		maxFileSize: maxFileSize,
	}
}

func (h *Folders) Show(c *ctx.Ctx) error {
	return h.show(c, nil)
}

func (h *Folders) Edit(c *ctx.Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	return c.HTML(http.StatusOK, "layouts/page", "edit_folder", Map{
		"folder":           folder,
		"validationErrors": validate.NewErrors(),
	})
}

func (h *Folders) Update(c *ctx.Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	b := FolderForm{}
	if err := bind.Form(c.Req, &b); err != nil {
		return errors.Join(httperror.BadRequest, err)
	}

	folder.Name = b.Name

	if err := h.repo.UpdateFolder(c.Context(), folder); err != nil {
		validationErrors := validate.NewErrors()
		if err != nil && !errors.As(err, &validationErrors) {
			return err
		}
		return c.HTML(http.StatusOK, "layouts/page", "edit_folder", Map{
			"folder":           folder,
			"validationErrors": validationErrors,
		})
	}

	c.RedirectTo("folder", "folderID", folder.ID)

	return nil
}

func (h *Folders) Delete(c *ctx.Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	if err := h.repo.DeleteFolder(c.Context(), folderID); err != nil {
		return err
	}

	c.AddFlash(ctx.Flash{
		Type:         "info",
		Body:         "Folder deleted succesfully",
		DismissAfter: 3 * time.Second,
	})
	c.RedirectTo("space", "spaceName", folder.Space.Name)

	return nil
}

func (h *Folders) UploadFile(c *ctx.Ctx) error {
	folderID := c.Path("folderID")

	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return httperror.NotFound
	}

	if !c.IsSpaceAdmin(c.User, folder.Space) {
		return httperror.Forbidden
	}

	/*
		TODO: retrieve content type by content sniffing
		without interfering with streaming body
	*/
	contentLength, _ := strconv.ParseInt(c.Req.Header.Get("Content-Length"), 10, 64)

	// request header only accepts ISO-8859-1 so we had to escape it
	uploadFilename, _ := url.QueryUnescape(c.Req.Header.Get("X-Upload-Filename"))

	file := &models.File{
		FolderID:    folderID,
		ID:          ulid.Make().String(),
		Name:        uploadFilename,
		ContentType: c.Req.Header.Get("Content-Type"),
		Size:        contentLength,
	}

	// TODO get size
	md5, err := h.file.Add(c.Context(), file.ID, c.Req.Body)
	if err != nil {
		return err
	}

	file.MD5 = md5

	if err = h.repo.CreateFile(c.Context(), file); err != nil {
		return h.show(c, err)
	}

	// reload folder
	folder, err = h.repo.FolderByID(c.Context(), file.FolderID)

	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "", "show_folder/files", Map{
		"folder":      folder,
		"maxFileSize": h.maxFileSize,
	})
}

func (h *Folders) Share(c *ctx.Ctx) error {
	folderID := c.Path("folderID")
	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, "layouts/public_page", "share_folder", Map{
		"folder": folder,
	})
}

func (h *Folders) show(c *ctx.Ctx, err error) error {
	validationErrors := validate.NewErrors()
	if err != nil && !errors.As(err, &validationErrors) {
		return err
	}

	folderID := c.Path("folderID")
	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}

	if turbo.Request(c.Req) {
		return turbo.Render(c.Res, c.Req, http.StatusOK,
			turbo.RemoveMatch(".modal.show, .modal-backdrop"),
			turbo.Update("files").Render(views.FolderFiles(c, folder)),
		)
	}

	return c.HTML(http.StatusOK, "layouts/page", "show_folder", Map{
		"folder":           folder,
		"validationErrors": validationErrors,
		"maxFileSize":      h.maxFileSize,
	})
}
