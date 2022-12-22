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

func (c *Files) Download(ctx *Ctx) error {
	fileID := ctx.Path("fileID")
	if _, err := c.repo.File(ctx.Context(), fileID); err != nil {
		return err
	}
	return c.file.Get(ctx.Context(), fileID, ctx.Res)
}

func (c *Files) Delete(ctx *Ctx) error {
	if ctx.User() == nil {
		return ErrUnauthorized
	}

	fileID := ctx.Path("fileID")
	file, err := c.repo.File(ctx.Context(), fileID)
	if err != nil {
		return err
	}
	if err := c.repo.DeleteFile(ctx.Context(), fileID); err != nil {
		return err
	}
	if err := c.file.Delete(ctx.Context(), fileID); err != nil {
		return err
	}

	ctx.PersistFlash(Flash{
		Type: Info,
		Body: "File deleted succesfully",
	})
	ctx.Redirect("folder", "folderID", file.FolderID)

	return nil
}
