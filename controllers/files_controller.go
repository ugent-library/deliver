package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
)

type FilesController struct {
	repo    *repositories.Repo
	storage objectstore.Store
}

func NewFilesController(r *repositories.Repo, s objectstore.Store) *FilesController {
	return &FilesController{
		repo:    r,
		storage: s,
	}
}

func (h *FilesController) Download(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	fileID := c.Path("fileID")

	_, err := h.repo.Files.Get(c.Context(), fileID)
	if err != nil && errors.Is(err, models.ErrNotFound) {
		return httperror.NotFound
	} else if err != nil {
		return err
	}

	if err := h.repo.Files.AddDownload(c.Context(), fileID); err != nil {
		return err
	}

	file, err := h.repo.Files.Get(c.Context(), fileID)
	if err != nil {
		return err
	}

	b, err := h.storage.Get(c.Context(), file.ID)
	if err != nil {
		return err
	}

	c.Hub.Send("folder."+file.FolderID,
		fmt.Sprintf(`"<span id="file-%s-downloads">%d</span>`, file.ID, file.Downloads),
	)

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", file.Name))

	_, err = io.Copy(w, b)

	return err
}

func (h *FilesController) Delete(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	fileID := c.Path("fileID")

	file, err := h.repo.Files.Get(c.Context(), fileID)
	if err != nil {
		return httperror.NotFound
	}

	if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
		return httperror.Forbidden
	}

	if err := h.repo.Files.Delete(c.Context(), fileID); err != nil {
		return err
	}

	htmx.AddTrigger(w, "refresh-files")

	return nil
}
