package handler

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
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"go.uber.org/zap"
)

const (
	Info = "info"

	flashSessionKey = "flash"
	userSessionKey  = "user"
)

type Empty struct{}

// TODO turn wrapper into an object?
// TODO instead of Config, make the argument a constructor function? (this also allows type inference)
// TODO make flash a generic type
func Wrapper[U, V any](c Config) func(fns ...func(*Ctx[U, V]) error) http.HandlerFunc {
	return func(fns ...func(*Ctx[U, V]) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := &Ctx[U, V]{
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
			for _, fn := range fns {
				if err := fn(ctx); err != nil {
					ctx.handleError(err)
					return
				}
			}
		}
	}
}

type Config struct {
	Log          *zap.SugaredLogger
	SessionStore sessions.Store
	SessionName  string
	Router       *mux.Router
}

// func Wrap[U any](c Config, fn func(*Ctx[U,V]) error) *Handler[U] {
// 	return &Handler[U]{c, fn}
// }

// type Handler[U any] struct {
// 	c  Config
// 	fn func(*Ctx[U,V]) error
// }

// func (h *Handler[U]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	ctx := &Ctx[U]{
// 		Config:    h.c,
// 		Res:       w,
// 		Req:       r,
// 		path:      mux.Vars(r),
// 		CSRFToken: csrf.Token(r),
// 		CSRFTag:   csrf.TemplateField(r),
// 	}
// 	if err := ctx.loadSession(); err != nil {
// 		ctx.handleError(err)
// 		return
// 	}
// 	if err := h.fn(ctx); err != nil {
// 		ctx.handleError(err)
// 	}
// }

type Flash struct {
	Type  string
	Title string
	Body  template.HTML
}

type Ctx[U, V any] struct {
	Config
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

// TODO register error handlers
func (c *Ctx[U, V]) handleError(err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
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
