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

// TODO make flash a generic type
// TODO constructor function that allows type inference?
// TODO don't pass whole wrapper to ctx
type Wrapper[U, V any] struct {
	Log          *zap.SugaredLogger
	SessionStore sessions.Store
	SessionName  string
	Router       *mux.Router
	ErrorHandler func(*Ctx[U, V], error)
}

func (c Wrapper[U, V]) Wrap(handlers ...func(*Ctx[U, V]) error) http.HandlerFunc {
	if c.ErrorHandler == nil {
		c.ErrorHandler = func(c *Ctx[U, V], err error) {
			http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx[U, V]{
			Wrapper:   c,
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

type Flash struct {
	Type  string
	Title string
	Body  template.HTML
}

type Ctx[U, V any] struct {
	Wrapper[U, V]
	Res       http.ResponseWriter
	Req       *http.Request
	path      map[string]string
	CSRFToken string
	CSRFTag   template.HTML
	Flash     []Flash
	user      *U
	Var       V
}

func (c *Ctx[U, V]) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx[U, V]) Path(k string) string {
	return c.path[k]
}

func (c *Ctx[U, V]) URL(route string, pairs ...string) *url.URL {
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

func (c *Ctx[U, V]) URLPath(route string, pairs ...string) *url.URL {
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

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type RenderData[U, V any] struct {
	*Ctx[U, V]
	Data any
}

func (c *Ctx[U, V]) Render(r Renderer, data any) error {
	return r.Render(c.Res, RenderData[U, V]{c, data})
}

func (c *Ctx[U, V]) Redirect(route string, pairs ...string) {
	http.Redirect(c.Res, c.Req, c.URLPath(route, pairs...).String(), http.StatusSeeOther)
}

func (c *Ctx[U, V]) User() *U {
	return c.user
}

func (c *Ctx[U, V]) SetUser(u *U) error {
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

func (c *Ctx[U, V]) DeleteUser() error {
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

func (c *Ctx[U, V]) PersistFlash(f Flash) error {
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

func (c *Ctx[U, V]) loadSession() error {
	s, err := c.SessionStore.Get(c.Req, c.SessionName)
	if err != nil {
		return fmt.Errorf("couldn't get session data: %w", err)
	}

	if user := s.Values[userSessionKey]; user != nil {
		c.user = user.(*U)
	}

	for _, f := range s.Flashes(flashSessionKey) {
		c.Flash = append(c.Flash, f.(Flash))
	}

	if err := s.Save(c.Req, c.Res); err != nil {
		return fmt.Errorf("couldn't save session data: %w", err)
	}

	return nil
}
