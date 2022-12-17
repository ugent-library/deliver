package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

type Spaces struct {
	repo     models.RepositoryService
	listView view.View
	showView view.View
}

func NewSpaces(r models.RepositoryService) *Spaces {
	return &Spaces{
		repo:     r,
		listView: view.MustNew("page", "list_spaces"),
		showView: view.MustNew("page", "show_space"),
	}
}

type SpaceForm struct {
	Name string `form:"name"`
}

func (c *Spaces) List(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	spaces, err := c.repo.Spaces(r.Context())
	if err != nil {
		return err
	}
	return c.listView.Render(w, ctx.Yield(Var{
		"spaces": spaces,
	}))
}

func (c *Spaces) Show(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	spaceID := mux.Vars(r)["spaceID"]
	space, err := c.repo.Space(context.TODO(), spaceID)
	if err != nil {
		return err
	}
	return c.showView.Render(w, ctx.Yield(Var{
		"space": space,
	}))
}

func (c *Spaces) Create(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	b := SpaceForm{}
	if err := bindForm(r, &b); err != nil {
		return err
	}

	space := &models.Space{
		Name: b.Name,
	}
	if err := c.repo.CreateSpace(context.TODO(), space); err != nil {
		return err
	}

	ctx.PersistFlash(w, r, Flash{
		Type: "info",
		Body: "Space created succesfully",
	})
	redirectURL := ctx.URLPath("space", "spaceID", space.ID).String()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	return nil
}
