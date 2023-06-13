package handlers

import (
	"net/http"
	"time"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/oidc"
)

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	claims := oidc.Claims{}
	if err := c.Auth.CompleteAuth(w, r, &claims); err != nil {
		c.HandleError(w, r, err)
		return
	}

	u := &models.User{
		Username: claims.PreferredUsername,
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
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	c := ctx.Get(r)

	if err := c.Auth.BeginAuth(w, r); err != nil {
		c.HandleError(w, r, err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, c.PathTo("home").String(), http.StatusSeeOther)
}
