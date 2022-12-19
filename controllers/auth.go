package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/oidc"
)

type Auth struct {
	oidcAuth *oidc.Auth
}

func NewAuth(oidcAuth *oidc.Auth) *Auth {
	return &Auth{
		oidcAuth: oidcAuth,
	}
}

func (c *Auth) Callback(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	claims := oidc.Claims{}
	if err := c.oidcAuth.CompleteAuth(w, r, &claims); err != nil {
		return err
	}
	if err := ctx.SetUser(w, r, &models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}); err != nil {
		return err
	}
	ctx.Redirect(w, r, "home")
	return nil
}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	return c.oidcAuth.BeginAuth(w, r)
}

func (c *Auth) Logout(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	if err := ctx.DeleteUser(w, r); err != nil {
		return err
	}
	ctx.Redirect(w, r, "home")
	return nil
}
