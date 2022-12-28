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
	c.Var.User = &models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}
	c.Session.Set(userKey, c.Var.User)
	c.Redirect("spaces")
	return nil
}

func (h *Auth) Login(c Ctx) error {
	return h.oidcAuth.BeginAuth(c.Res, c.Req)
}

func (h *Auth) Logout(c Ctx) error {
	c.Session.Delete(userKey)
	c.Redirect("home")
	return nil
}
