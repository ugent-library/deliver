package controllers

import (
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

func (h *Files) Download(c Ctx) error {
	fileID := c.Path("fileID")
	if _, err := h.repo.File(c.Context(), fileID); err != nil {
		return err
	}
	return h.file.Get(c.Context(), fileID, c.Res)
}

func (h *Files) Delete(c Ctx) error {
	fileID := c.Path("fileID")
	file, err := h.repo.File(c.Context(), fileID)
	if err != nil {
		return err
	}
	if err := h.repo.DeleteFile(c.Context(), fileID); err != nil {
		return err
	}
	if err := h.file.Delete(c.Context(), fileID); err != nil {
		return err
	}

	c.PersistFlash(Flash{
		Type: Info,
		Body: "File deleted succesfully",
	})
	c.Redirect("folder", "folderID", file.FolderID)

	return nil
}
