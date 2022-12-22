package controllers

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/dilliver/models"
	"go.uber.org/zap"
)

const (
	Info = "info"

	flashSessionKey = "flash"
	userSessionKey  = "user"
)

// TODO turn wrapper into an object
func Wrapper(c Config) func(func(*Ctx) error) http.HandlerFunc {
	return func(fn func(*Ctx) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := &Ctx{
				Config:    c,
				Res:       w,
				Req:       r,
				path:      mux.Vars(r),
				CSRFToken: csrf.Token(r),
				CSRFTag:   csrf.TemplateField(r),
			}
			if err := ctx.loadSession(); err != nil {
				ctx.handleError(err)
				return
			}
			if err := fn(ctx); err != nil {
				ctx.handleError(err)
			}
		}
	}
}

var (
	ErrUnauthorized = &HTTPError{http.StatusUnauthorized}
)

type HTTPError struct {
	Code int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.Code, http.StatusText(e.Code))
}

type Config struct {
	Log          *zap.SugaredLogger
	SessionStore sessions.Store
	SessionName  string
	Router       *mux.Router
}

type Flash struct {
	Type  string
	Title string
	Body  template.HTML
}

type Var map[string]any

type Ctx struct {
	Config
	Res       http.ResponseWriter
	Req       *http.Request
	path      map[string]string
	CSRFToken string
	CSRFTag   template.HTML
	Flash     []Flash
	Var       any
	user      *models.User
}

func (c *Ctx) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx) Path(k string) string {
	return c.path[k]
}

func (c *Ctx) Yield(v any) *Ctx {
	c.Var = v
	return c
}

func (c *Ctx) URL(route string, pairs ...string) *url.URL {
	r := c.Router.Get(route)
	if r == nil {
		panic(fmt.Errorf("unknown route '%s'", route))
	}
	u, err := r.URL(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	if u.Host == "" {
		u.Host = c.Req.Host
	}
	if u.Scheme == "" {
		u.Scheme = c.Req.URL.Scheme
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u
}

func (c *Ctx) URLPath(route string, pairs ...string) *url.URL {
	r := c.Router.Get(route)
	if r == nil {
		panic(fmt.Errorf("unknown route '%s'", route))
	}
	u, err := r.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

func (c *Ctx) Redirect(route string, pairs ...string) {
	http.Redirect(c.Res, c.Req, c.URLPath(route, pairs...).String(), http.StatusSeeOther)
}

func (c *Ctx) User() *models.User {
	return c.user
}

func (c *Ctx) SetUser(u *models.User) error {
	s, err := c.SessionStore.Get(c.Req, c.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}
	s.Values[userSessionKey] = u
	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}
	c.user = u
	return nil
}

func (c *Ctx) DeleteUser() error {
	s, err := c.SessionStore.Get(c.Req, c.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}
	delete(s.Values, userSessionKey)
	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}
	c.user = nil
	return nil
}

func (c *Ctx) PersistFlash(f Flash) error {
	s, err := c.SessionStore.Get(c.Req, c.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}
	s.AddFlash(f, flashSessionKey)
	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}
	return nil
}

func (c *Ctx) loadSession() error {
	s, err := c.SessionStore.Get(c.Req, c.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}

	if user := s.Values[userSessionKey]; user != nil {
		c.user = user.(*models.User)
	}

	for _, f := range s.Flashes(flashSessionKey) {
		c.Flash = append(c.Flash, f.(Flash))
	}

	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}

	return nil
}

// TODO register error handlers
func (c *Ctx) handleError(err error) {
	if err == models.ErrNotFound {
		err = &HTTPError{Code: http.StatusNotFound}
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		httpErr = &HTTPError{Code: http.StatusInternalServerError}
	}

	switch httpErr.Code {
	case http.StatusNotFound:
		// TODO use controller action directly
		c.Router.NotFoundHandler.ServeHTTP(c.Res, c.Req)
	case http.StatusUnauthorized:
		c.Redirect("login")
	default:
		http.Error(c.Res, http.StatusText(httpErr.Code), httpErr.Code)
	}
}
