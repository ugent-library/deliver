package controllers

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/bind"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/validate"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/httphelpers"
)

type Folders struct {
	repo models.RepositoryService
	file models.FileService
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService, f models.FileService) *Folders {
	return &Folders{
		repo: r,
		file: f,
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

	return c.HTML(http.StatusOK, "page", "edit_folder", Map{
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
		return c.HTML(http.StatusOK, "page", "edit_folder", Map{
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

		mediatype, err := httphelpers.DetectContentType(f)
		if err != nil {
			return err
		}

		file := &models.File{
			FolderID:    folderID,
			ID:          ulid.Make().String(),
			Name:        fileHeader.Filename,
			ContentType: mediatype,
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

func (h *Folders) Share(c *Ctx) error {
	folderID := c.Path("folderID")
	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return err
	}
	var folderSize int64 = 0
	for _, file := range folder.Files {
		folderSize += file.Size
	}
	return c.HTML(http.StatusOK, "simple_page", "share_folder", Map{
		"folder":     folder,
		"folderSize": folderSize,
	})
}

func (h *Folders) Download(c *Ctx) error {

	fmt.Fprintf(os.Stderr, "download folder\n")
	folderID := c.Path("folderID")
	folder, err := h.repo.FolderByID(c.Context(), folderID)
	if err != nil {
		return httperror.NotFound
	}

	zipFileName := fmt.Sprintf(
		"%s-%s.zip",
		folder.ID,
		folder.Slug(),
	)
	c.Res.Header().Add("Content-Type", "application/zip")
	c.Res.Header().Add(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename*=UTF-8''%s", zipFileName),
	)

	zipFh := zip.NewWriter(bufio.NewWriter(c.Res))
	defer zipFh.Close()

	for _, file := range folder.Files {

		fileWriter, err := zipFh.CreateHeader(&zip.FileHeader{
			Name:     file.Name,
			Method:   zip.Store, //no compression for streaming
			Modified: time.Now(),
		})
		if err != nil {
			return err
		}

		// increment file download count
		if err := h.repo.AddFileDownload(c.Context(), file.ID); err != nil {
			return err
		}

		// add file contents to zip
		h.file.Get(c.Context(), file.ID, fileWriter)
	}

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
	return c.HTML(http.StatusOK, "page", "show_folder", Map{
		"folder":           folder,
		"validationErrors": validationErrors,
	})
}
