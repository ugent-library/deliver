package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Folders struct {
	repository models.RepositoryService
	listView   view.View
}

func NewFolders(r models.RepositoryService) *Folders {
	return &Folders{
		repository: r,
		listView:   view.MustNew("page", "list_folders"),
	}
}

type FolderForm struct {
	Name string `form:"name"`
}

type ListFoldersVars struct {
	SpaceID string
	Folders []*models.Folder
}

func (c *Folders) List(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	spaceID := mux.Vars(r)["space_id"]
	folders, err := c.repository.Folders(context.Background(), spaceID)
	if err != nil {
		panic(err) // TODO
	}
	c.listView.Render(w, yield(ctx, ListFoldersVars{
		SpaceID: spaceID,
		Folders: folders,
	}))
}

func (c *Folders) Create(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	spaceID := mux.Vars(r)["space_id"]
	b := FolderForm{}
	if err := bindForm(r, &b); err != nil {
		panic(err) // TODO
	}
	folder := &models.Folder{
		SpaceID: spaceID,
		Name:    b.Name,
	}
	if err := c.repository.CreateFolder(context.Background(), folder); err != nil {
		panic(err) // TODO
	}
	ctx.PersistFlash(w, r, Flash{
		Type: Info,
		Body: "Folder created succesfully",
	})
	redirectURL := ctx.URLPath("folders", "space_id", spaceID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
