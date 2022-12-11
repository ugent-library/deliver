package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

const (
	Info = "info"

	flashSessionKey = "flash"
)

func Wrapper(c Config) func(func(http.ResponseWriter, *http.Request, Ctx)) http.HandlerFunc {
	return func(fn func(http.ResponseWriter, *http.Request, Ctx)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := Ctx{
				Log:          c.Log,
				router:       c.Router,
				sessionStore: c.SessionStore,
				sessionName:  c.SessionName,
			}
			if err := ctx.loadSession(w, r); err != nil {
				// TODO handle error gracefully
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			fn(w, r, ctx)
		}
	}
}

type Config struct {
	Log          *zap.SugaredLogger
	Router       *mux.Router
	SessionStore sessions.Store
	SessionName  string
}

type Flash struct {
	Type string
	Body template.HTML
}

type Var map[string]any

type Ctx struct {
	Log          *zap.SugaredLogger
	router       *mux.Router
	sessionStore sessions.Store
	sessionName  string
	Flash        []Flash
	Var          Var
}

func (c Ctx) Yield(v Var) Ctx {
	c.Var = v
	return c
}

func (c Ctx) URL(route string, pairs ...string) *url.URL {
	r := c.router.Get(route)
	if r == nil {
		panic(fmt.Errorf("unknown route '%s'", route))
	}
	u, err := r.URL(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

func (c Ctx) URLPath(route string, pairs ...string) *url.URL {
	r := c.router.Get(route)
	if r == nil {
		panic(fmt.Errorf("unknown route '%s'", route))
	}
	u, err := r.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

func (c Ctx) PersistFlash(w http.ResponseWriter, r *http.Request, f Flash) error {
	s, err := c.sessionStore.Get(r, c.sessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}

	s.AddFlash(f, flashSessionKey)

	if err := s.Save(r, w); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}
	return nil
}

func (c *Ctx) loadSession(w http.ResponseWriter, r *http.Request) error {
	s, err := c.sessionStore.Get(r, c.sessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}

	flashes := s.Flashes(flashSessionKey)

	if err := s.Save(r, w); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}

	for _, f := range flashes {
		c.Flash = append(c.Flash, f.(Flash))
	}

	return nil
}
