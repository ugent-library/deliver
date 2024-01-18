package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/htmx"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	file := ctx.GetFile(r)

	if err := c.Repo.Files.AddDownload(r.Context(), file.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	file, err := c.Repo.Files.Get(r.Context(), file.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	b, err := c.Storage.Get(r.Context(), file.ID)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	c.Hub.SendString("folder."+file.FolderID,
		fmt.Sprintf(`"<span id="file-%s-downloads">%d</span>`, file.ID, file.Downloads),
	)

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", file.Name))

	io.Copy(w, b)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	folder := ctx.GetFolder(r)

	// TODO: retrieve content type by content sniffing without interfering with streaming body
	size, _ := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

	// request header only accepts ISO-8859-1 so we had to escape it
	name, _ := url.QueryUnescape(r.Header.Get("X-Upload-Filename"))

	file := &models.File{
		FolderID:    folder.ID,
		ID:          ulid.Make().String(),
		Name:        name,
		ContentType: r.Header.Get("Content-Type"),
		Size:        size,
	}

	// TODO get size
	md5, err := c.Storage.Add(r.Context(), file.ID, r.Body)
	if err != nil {
		c.HandleError(w, r, err)
		return
	}

	file.MD5 = md5

	if err := c.Repo.Files.Create(r.Context(), file); err != nil {
		c.HandleError(w, r, err)
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)
	file := ctx.GetFile(r)

	if err := c.Repo.Files.Delete(r.Context(), file.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	htmx.AddTrigger(w, "refresh-files")
}
