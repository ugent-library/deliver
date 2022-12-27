package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

const (
	Info = "info"

	flashSessionKey = "flash"
	userSessionKey  = "user"
)

type Unused struct{}

// TODO constructor function that allows type inference?
// TODO don't pass whole Config to ctx
// TODO make Router and Session Interfaces
// TODO view package: DefaultConfig object
type Config[U, V, F any] struct {
	Log          *zap.SugaredLogger
	SessionStore sessions.Store
	SessionName  string
	Router       *mux.Router
	ErrorHandler func(*Ctx[U, V, F], error)
}

func (c Config[U, V, F]) Wrap(handlers ...func(*Ctx[U, V, F]) error) http.HandlerFunc {
	if c.ErrorHandler == nil {
		c.ErrorHandler = func(c *Ctx[U, V, F], err error) {
			http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx[U, V, F]{
			config:    c,
			Res:       w,
			Req:       r,
			path:      mux.Vars(r),
			CSRFToken: csrf.Token(r),
			CSRFTag:   csrf.TemplateField(r),
		}
		if err := ctx.loadSession(); err != nil {
			c.ErrorHandler(ctx, err)
			return
		}
		for _, fn := range handlers {
			if err := fn(ctx); err != nil {
				c.ErrorHandler(ctx, err)
				return
			}
		}
	}
}

type Ctx[U, V, F any] struct {
	config    Config[U, V, F]
	Res       http.ResponseWriter
	Req       *http.Request
	path      map[string]string
	CSRFToken string
	CSRFTag   template.HTML
	Flash     []F
	user      *U
	Var       V
}

func (c *Ctx[U, V, F]) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx[U, V, F]) Path(k string) string {
	return c.path[k]
}

func (c *Ctx[U, V, F]) URL(route string, pairs ...string) *url.URL {
	r := c.config.Router.Get(route)
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

func (c *Ctx[U, V, F]) URLPath(route string, pairs ...string) *url.URL {
	r := c.config.Router.Get(route)
	if r == nil {
		panic(fmt.Errorf("unknown route '%s'", route))
	}
	u, err := r.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

func (c *Ctx[U, V, F]) ExecuteHandler(route string) error {
	r := c.config.Router.Get(route)
	if r == nil {
		return fmt.Errorf("unknown route '%s'", route)
	}
	r.GetHandler().ServeHTTP(c.Res, c.Req)
	return nil
}

func (c *Ctx[U, V, F]) Redirect(route string, pairs ...string) {
	http.Redirect(c.Res, c.Req, c.URLPath(route, pairs...).String(), http.StatusSeeOther)
}

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type RenderData[U, V, F any] struct {
	*Ctx[U, V, F]
	Data any
}

func (c *Ctx[U, V, F]) Render(r Renderer, data any) error {
	return r.Render(c.Res, RenderData[U, V, F]{c, data})
}

func (c *Ctx[U, V, F]) User() *U {
	return c.user
}

func (c *Ctx[U, V, F]) SetUser(u *U) error {
	s, err := c.config.SessionStore.Get(c.Req, c.config.SessionName)
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

func (c *Ctx[U, V, F]) DeleteUser() error {
	s, err := c.config.SessionStore.Get(c.Req, c.config.SessionName)
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

func (c *Ctx[U, V, F]) PersistFlash(f F) error {
	s, err := c.config.SessionStore.Get(c.Req, c.config.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}
	s.AddFlash(f, flashSessionKey)
	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}
	return nil
}

func (c *Ctx[U, V, F]) loadSession() error {
	s, err := c.config.SessionStore.Get(c.Req, c.config.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}

	if user := s.Values[userSessionKey]; user != nil {
		c.user = user.(*U)
	}

	for _, f := range s.Flashes(flashSessionKey) {
		c.Flash = append(c.Flash, f.(F))
	}

	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}

	return nil
}
