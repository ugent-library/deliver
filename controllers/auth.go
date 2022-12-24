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

func (h *Auth) Callback(c Ctx) error {
	claims := oidc.Claims{}
	if err := h.oidcAuth.CompleteAuth(c.Res, c.Req, &claims); err != nil {
		return err
	}
	if err := c.SetUser(&models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}); err != nil {
		return err
	}
	c.Redirect("spaces")
	return nil
}

func (h *Auth) Login(c Ctx) error {
	return h.oidcAuth.BeginAuth(c.Res, c.Req)
}

func (h *Auth) Logout(c Ctx) error {
	if err := c.DeleteUser(); err != nil {
		return err
	}
	c.Redirect("home")
	return nil
}
