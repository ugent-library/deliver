package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/autosession"
	"go.uber.org/zap"
)

// TODO Var constructor function that allows type inference?
// TODO make Routes interface
// TODO get request scoped logger from context
// TODO embed original error in httperror
type Config[V any] struct {
	Log          *zap.SugaredLogger
	Router       *mux.Router
	Before       []func(*Ctx[V]) error
	ErrorHandler func(*Ctx[V], error)
}

func (config Config[V]) Wrap(handlers ...func(*Ctx[V]) error) http.HandlerFunc {
	if config.ErrorHandler == nil {
		config.ErrorHandler = func(c *Ctx[V], err error) {
			http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		c := &Ctx[V]{
			Log:       config.Log,
			Res:       w,
			Req:       r,
			Session:   autosession.Get(r.Context()),
			path:      mux.Vars(r),
			CSRFToken: csrf.Token(r),
			CSRFTag:   csrf.TemplateField(r),
			router:    config.Router,
		}
		for _, fn := range config.Before {
			if err := fn(c); err != nil {
				config.ErrorHandler(c, err)
				return
			}
		}
		for _, fn := range handlers {
			if err := fn(c); err != nil {
				config.ErrorHandler(c, err)
				return
			}
		}
	}
}

type Ctx[V any] struct {
	Log       *zap.SugaredLogger
	Req       *http.Request
	Res       http.ResponseWriter
	Session   autosession.Session
	path      map[string]string
	CSRFToken string
	CSRFTag   template.HTML
	Var       V
	router    *mux.Router
}

func (c *Ctx[V]) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx[V]) Path(k string) string {
	return c.path[k]
}

func (c *Ctx[V]) URL(route string, pairs ...string) *url.URL {
	r := c.router.Get(route)
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

func (c *Ctx[V]) URLPath(route string, pairs ...string) *url.URL {
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

func (c *Ctx[V]) ExecuteHandler(route string) error {
	r := c.router.Get(route)
	if r == nil {
		return fmt.Errorf("unknown route '%s'", route)
	}
	r.GetHandler().ServeHTTP(c.Res, c.Req)
	return nil
}

func (c *Ctx[V]) RedirectTo(route string, pairs ...string) {
	http.Redirect(c.Res, c.Req, c.URLPath(route, pairs...).String(), http.StatusSeeOther)
}

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type renderData[V any] struct {
	*Ctx[V]
	Data any
}

func (c *Ctx[V]) Render(r Renderer, data any) error {
	return r.Render(c.Res, renderData[V]{c, data})
}
