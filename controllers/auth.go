package controllers

import (
	"fmt"
	"log"
	"net/http"

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

// TODO store user in session
func (c *Auth) Callback(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	var claims oidc.Claims
	if err := c.oidcAuth.CompleteAuth(w, r, &claims); err != nil {
		log.Printf("%+v", err)
		return err
	}
	w.Write([]byte(fmt.Sprintf("%+v", claims)))
	return nil
}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	return c.oidcAuth.BeginAuth(w, r)
}

// TODO remove user from session
func (c *Auth) Logout(w http.ResponseWriter, r *http.Request, ctx Ctx) error {
	ctx.RedirectTo(w, r, "home")
	return nil
}
