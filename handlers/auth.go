package handlers

import (
	"net/http"
	"time"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/oidc"
)

type AuthHandler struct {
	auth       *oidc.Auth
	matchClaim string
}

func NewAuthHandler(auth *oidc.Auth, matchClaim string) *AuthHandler {
	return &AuthHandler{
		auth:       auth,
		matchClaim: matchClaim,
	}
}

func (h *AuthHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	claims := oidc.Claims{}
	if err := h.auth.CompleteAuth(w, r, &claims); err != nil {
		c.HandleError(w, r, err)
		return
	}

	u := &models.User{
		Username: claims.GetString(h.matchClaim),
		Name:     claims.Name,
		Email:    claims.Email,
	}
	if err := c.Repo.Users.CreateOrUpdate(r.Context(), u); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ctx.RememberCookie,
		Value:    u.RememberToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	http.Redirect(w, r, c.Path("home").String(), http.StatusSeeOther)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if err := h.auth.BeginAuth(w, r); err != nil {
		c.HandleError(w, r, err)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if err := c.Repo.Users.RenewRememberToken(r.Context(), c.User.ID); err != nil {
		c.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ctx.RememberCookie,
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	http.Redirect(w, r, c.Path("home").String(), http.StatusSeeOther)
}
