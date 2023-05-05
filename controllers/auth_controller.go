package controllers

import (
	"net/http"
	"time"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/oidc"
)

type AuthController struct {
	repo     *repositories.Repo
	oidcAuth *oidc.Auth
}

func NewAuthController(repo *repositories.Repo, oidcAuth *oidc.Auth) *AuthController {
	return &AuthController{
		repo:     repo,
		oidcAuth: oidcAuth,
	}
}

func (h *AuthController) Callback(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	claims := oidc.Claims{}
	if err := h.oidcAuth.CompleteAuth(w, r, &claims); err != nil {
		return err
	}
	u := &models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}
	if err := h.repo.Users.CreateOrUpdate(c.Context(), u); err != nil {
		return err
	}
	c.User = u
	c.Cookies.Set(rememberCookie, u.RememberToken, time.Now().Add(24*time.Hour*7))
	c.RedirectTo("home")
	return nil
}

func (h *AuthController) Login(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	return h.oidcAuth.BeginAuth(w, r)
}

func (h *AuthController) Logout(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	c.Cookies.Delete(rememberCookie)
	if err := h.repo.Users.RenewRememberToken(c.Context(), c.User.ID); err != nil {
		return err
	}
	c.RedirectTo("home")
	return nil
}
