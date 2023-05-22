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

func (h *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	claims := oidc.Claims{}
	if err := h.oidcAuth.CompleteAuth(w, r, &claims); err != nil {
		c.HandleError(w, r, err)
		return
	}

	u := &models.User{
		Username: claims.PreferredUsername,
		Name:     claims.Name,
		Email:    claims.Email,
	}
	if err := h.repo.Users.CreateOrUpdate(r.Context(), u); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ctx.RememberCookie,
		Value:    u.RememberToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusSeeOther)
}

func (h *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	if err := h.oidcAuth.BeginAuth(w, r); err != nil {
		c.HandleError(w, r, err)
	}
}

func (h *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r.Context())

	if err := h.repo.Users.RenewRememberToken(r.Context(), c.User.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ctx.RememberCookie,
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusSeeOther)
}
