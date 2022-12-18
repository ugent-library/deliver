package controllers

import (
	"net/http"

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

func (c *Files) Download(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	fileID := ctx.Path("fileID")
	if _, err := c.repo.File(r.Context(), fileID); err != nil {
		return err
	}
	return c.file.Get(r.Context(), fileID, w)
}

func (c *Files) Delete(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	fileID := ctx.Path("fileID")
	file, err := c.repo.File(r.Context(), fileID)
	if err != nil {
		return err
	}
	if err := c.repo.DeleteFile(r.Context(), fileID); err != nil {
		return err
	}
	if err := c.file.Delete(r.Context(), fileID); err != nil {
		return err
	}
	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "File deleted succesfully",
	})
	redirectURL := ctx.URLPath("folder", "folderID", file.FolderID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	return nil
}
