package controllers

import (
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

func (c *Auth) Callback(ctx Ctx) error {
	claims := oidc.Claims{}
	if err := c.oidcAuth.CompleteAuth(ctx.Res, ctx.Req, &claims); err != nil {
		return err
	}
	if err := ctx.SetUser(&models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}); err != nil {
		return err
	}
	ctx.Redirect("spaces")
	return nil
}

func (c *Auth) Login(ctx Ctx) error {
	return c.oidcAuth.BeginAuth(ctx.Res, ctx.Req)
}

func (c *Auth) Logout(ctx Ctx) error {
	if err := ctx.DeleteUser(); err != nil {
		return err
	}
	ctx.Redirect("home")
	return nil
}
