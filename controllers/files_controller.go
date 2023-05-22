package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
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

func (h *FilesController) Download(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())
	file := ctx.GetFile(r.Context())

	if err := h.repo.Files.AddDownload(r.Context(), file.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	file, err := h.repo.Files.Get(r.Context(), file.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	b, err := h.storage.Get(r.Context(), file.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	c.Hub.Send("folder."+file.FolderID,
		fmt.Sprintf(`"<span id="file-%s-downloads">%d</span>`, file.ID, file.Downloads),
	)

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", file.Name))

	io.Copy(w, b)
}

func (h *FilesController) Delete(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())
	file := ctx.GetFile(r.Context())

	if err := h.repo.Files.Delete(r.Context(), file.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	htmx.AddTrigger(w, "refresh-files")
}
