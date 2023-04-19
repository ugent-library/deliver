package controllers

import (
	"fmt"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
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

func (h *Files) Download(c *ctx.Ctx) error {
	fileID := c.Path("fileID")

	if err := h.repo.AddFileDownload(c.Context(), fileID); err != nil {
		return err
	}

	file, err := h.repo.FileByID(c.Context(), fileID)
	if err != nil {
		return err
	}

	c.Hub.Send("folder."+file.FolderID,
		fmt.Sprintf(`"<span id="file-%s-downloads">%d</span>`, file.ID, file.Downloads),
	)

	c.Res.Header().Add("Content-Disposition", "attachment")

	return h.file.Get(c.Context(), fileID, c.Res)
}

func (h *Files) Delete(c *ctx.Ctx) error {
	fileID := c.Path("fileID")

	file, err := h.repo.FileByID(c.Context(), fileID)
	if err != nil {
		return httperror.NotFound
	}

	if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
		return httperror.Forbidden
	}

	if err := h.repo.DeleteFile(c.Context(), fileID); err != nil {
		return err
	}

	htmx.AddTrigger(c.Res, "refresh-files")

	return nil
}
