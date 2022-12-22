package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/ulid"
	"github.com/ugent-library/dilliver/view"
)

type Folders struct {
	repo     models.RepositoryService
	file     models.FileService
	showView view.View
}

type FolderForm struct {
	Name string `form:"name"`
}

func NewFolders(r models.RepositoryService, f models.FileService) *Folders {
	return &Folders{
		repo:     r,
		file:     f,
		showView: view.MustNew("page", "show_folder"),
	}
}

func (c *Folders) Show(ctx Ctx) error {
	folderID := ctx.Path("folderID")
	folder, err := c.repo.Folder(ctx.Context(), folderID)
	if err != nil {
		return err
	}
	return c.showView.Render(ctx.Res, ctx.Yield(Var{
		"folder": folder,
	}))
}

func (c *Folders) Create(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	spaceID := ctx.Path("spaceID")
	b := FolderForm{}
	if err := bindForm(ctx.Req, &b); err != nil {
		return err
	}

	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}
	if err := c.repo.CreateFolder(ctx.Context(), folder); err != nil {
		return err
	}

	ctx.PersistFlash(handler.Flash{
		Type: handler.Info,
		Body: "Folder created succesfully",
	})
	ctx.Redirect("folder", "folderID", folder.ID)

	return nil
}

// TODO remove files
func (c *Folders) Delete(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	folderID := ctx.Path("folderID")

	folder, err := c.repo.Folder(ctx.Context(), folderID)
	if err != nil {
		return err
	}

	if err := c.repo.DeleteFolder(ctx.Context(), folderID); err != nil {
		return err
	}

	ctx.PersistFlash(handler.Flash{
		Type: handler.Info,
		Body: "Folder deleted succesfully",
	})
	ctx.Redirect("space", "spaceID", folder.SpaceID)

	return nil
}

func (c *Folders) UploadFile(ctx Ctx) error {
	if ctx.User() == nil {
		return httperror.Unauthorized
	}

	folderID := ctx.Path("folderID")

	// 2GB limit on request body
	ctx.Req.Body = http.MaxBytesReader(ctx.Res, ctx.Req.Body, 2_000_000_000)
	// buffer limit of 32MB
	if err := ctx.Req.ParseMultipartForm(32_000_000); err != nil {
		return err
	}

	for _, fileHeader := range ctx.Req.MultipartForm.File["file"] {
		f, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		contentType, err := detectContentType(f)
		if err != nil {
			return err
		}

		file := &models.File{
			FolderID:    folderID,
			ID:          ulid.MustGenerate(),
			Name:        fileHeader.Filename,
			ContentType: contentType,
			Size:        fileHeader.Size,
		}

		// TODO get size
		md5, err := c.file.Add(ctx.Context(), file.ID, f)
		if err != nil {
			return err
		}

		file.Md5 = md5

		if err = c.repo.CreateFile(ctx.Context(), file); err != nil {
			return err
		}
	}

	ctx.PersistFlash(handler.Flash{
		Type: handler.Info,
		Body: "File added succesfully",
	})
	ctx.Redirect("folder", "folderID", folderID)

	return nil
}
