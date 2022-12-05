package controllers

import (
	"context"
	"net/http"

	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Spaces struct {
	repository models.RepositoryService
	listView   view.View
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repository: r,
		listView:   view.MustNew("page", "list_spaces"),
	}
}

type SpaceForm struct {
	Name string `form:"name"`
}

type ListSpacesVars struct {
	Spaces []*models.Space
}

func (c *Spaces) List(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	spaces, err := c.repository.Spaces(context.Background())
	if err != nil {
		panic(err) // TODO
	}
	c.listView.Render(w, yield(ctx, ListSpacesVars{Spaces: spaces}))
}

func (c *Spaces) Create(w http.ResponseWriter, r *http.Request, ctx Ctx) {
	b := SpaceForm{}
	if err := bindForm(r, &b); err != nil {
		panic(err) // TODO
	}
	space := &models.Space{Name: b.Name}
	if err := c.repository.CreateSpace(context.Background(), space); err != nil {
		panic(err) // TODO
	}
	http.Redirect(w, r, ctx.URLPath("spaces").String(), http.StatusSeeOther)
}
