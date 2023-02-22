package controllers

import (
	"time"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
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

func (h *Files) Download(c *Ctx) error {
	fileID := c.Path("fileID")
	if err := h.repo.AddFileDownload(c.Context(), fileID); err != nil {
		return err
	}
	c.Res.Header().Add("Content-Disposition", "attachment")
	return h.file.Get(c.Context(), fileID, c.Res)
}

func (h *Files) Delete(c *Ctx) error {
	fileID := c.Path("fileID")

	file, err := h.repo.FileByID(c.Context(), fileID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
		return httperror.Forbidden
	}

	if err := h.repo.DeleteFile(c.Context(), fileID); err != nil {
		return err
	}

	c.Session.Append(flashKey, Flash{
		Type:         infoFlash,
		Body:         "File deleted succesfully",
		DismissAfter: 3 * time.Second,
	})
	c.RedirectTo("folder", "folderID", file.FolderID)

	return nil
}
