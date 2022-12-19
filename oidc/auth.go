package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
)

type Config struct {
	URL              string
	ClientID         string
	ClientSecret     string
	RedirectURL      string
	CookieSecret     []byte
	CookieName       string
	AdditionalScopes []string
}

type Auth struct {
	oauthClient   *oauth2.Config
	tokenVerifier *oidc.IDTokenVerifier
	secureCookie  *securecookie.SecureCookie
	cookieName    string
}

func NewAuth(ctx context.Context, c Config) (*Auth, error) {
	oidcProvider, err := oidc.NewProvider(ctx, c.URL)
	if err != nil {
		return nil, err
	}

	// configure an oidc aware oauth2 client
	oauthClient := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     oidcProvider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	if c.AdditionalScopes != nil {
		oauthClient.Scopes = append(oauthClient.Scopes, c.AdditionalScopes...)
	} else {
		oauthClient.Scopes = append(oauthClient.Scopes, "profile")
	}

	tokenVerifier := oidcProvider.Verifier(&oidc.Config{ClientID: c.ClientID})

	auth := &Auth{
		oauthClient:   oauthClient,
		tokenVerifier: tokenVerifier,
		secureCookie:  securecookie.New(c.CookieSecret, nil),
		cookieName:    c.CookieName,
	}

	if auth.cookieName == "" {
		auth.cookieName = "oidc.state"
	}

	return auth, nil
}

func (a *Auth) BeginAuth(w http.ResponseWriter, r *http.Request) error {
	state, err := generateRandomState()
	if err != nil {
		return err
	}

	h, err := a.secureCookie.Encode(a.cookieName, state)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.cookieName,
		Value:    h,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, a.oauthClient.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return nil
}

func (a *Auth) CompleteAuth(w http.ResponseWriter, r *http.Request, claims any) error {
	cookie, err := r.Cookie(a.cookieName)
	if err != nil {
		return err
	}

	var state string
	if err = a.secureCookie.Decode(a.cookieName, cookie.Value, &state); err != nil {
		return err
	}

	if r.URL.Query().Get("state") != state {
		return errors.New("oidc: invalid state parameter")
	}

	oauthToken, err := a.oauthClient.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return err
	}

	// extract id token from oauth2 token
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return errors.New("oidc: id token missing")
	}

	// verify id token
	idToken, err := a.tokenVerifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return err
	}

	return idToken.Claims(claims)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
