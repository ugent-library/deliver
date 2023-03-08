package controllers

import (
	"time"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/oidc"
)

type Auth struct {
	repo     models.RepositoryService
	oidcAuth *oidc.Auth
}

func NewAuth(repo models.RepositoryService, oidcAuth *oidc.Auth) *Auth {
	return &Auth{
		repo:     repo,
		oidcAuth: oidcAuth,
	}
}

func (h *Auth) Callback(c *Ctx) error {
	claims := oidc.Claims{}
	if err := h.oidcAuth.CompleteAuth(c.Res, c.Req, &claims); err != nil {
		return err
	}
	u := &models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}
	if err := h.repo.CreateOrUpdateUser(c.Context(), u); err != nil {
		return err
	}
	c.User = u
	c.Cookies.Set(rememberCookie, u.RememberToken, time.Now().Add(24*time.Hour*7))
	c.RedirectTo("home")
	return nil
}

func (h *Auth) Login(c *Ctx) error {
	return h.oidcAuth.BeginAuth(c.Res, c.Req)
}

func (h *Auth) Logout(c *Ctx) error {
	c.Cookies.Delete(rememberCookie)
	if err := h.repo.RenewUserRememberToken(c.Context(), c.User.ID); err != nil {
		return err
	}
	c.RedirectTo("home")
	return nil
}
