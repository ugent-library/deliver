package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type Flag int

const (
	Vacuum Flag = iota
)

// TODO add Flash methods instead of Pop, Append?
type Session interface {
	Get(string) any
	Pop(string) any
	Set(string, any)
	Append(string, any)
	Delete(string)
	Clear()
	Save() error
}

type gorillaSession struct {
	session *sessions.Session
	changed bool
	req     *http.Request
	res     http.ResponseWriter
}

func (s *gorillaSession) Get(k string) any {
	return s.session.Values[k]
}

func (s *gorillaSession) Pop(k string) any {
	if v, ok := s.session.Values[k]; ok {
		delete(s.session.Values, k)
		s.changed = true
		return v
	}
	return nil
}

func (s *gorillaSession) Set(k string, v any) {
	s.session.Values[k] = v
	s.changed = true
}

func (s *gorillaSession) Append(k string, v any) {
	if vals, ok := s.session.Values[k]; ok {
		s.session.Values[k] = append(vals.([]any), v)
	}
	s.session.Values[k] = []any{v}
	s.changed = true
}

func (s *gorillaSession) Delete(k string) {
	delete(s.session.Values, k)
	s.changed = true
}

func (s *gorillaSession) Clear() {
	for k := range s.session.Values {
		delete(s.session.Values, k)
	}
	s.changed = true
}

func (s *gorillaSession) Save() error {
	if s.changed {
		return s.session.Save(s.req, s.res)
	}
	return nil
}

// TODO Var constructor function that allows type inference?
// TODO make Routes interface
// TODO make a strongly typed session using mapstructure?
// TODO request ID and request scoped logger
// TODO pass error in httperror
type Config[V any] struct {
	Log          *zap.SugaredLogger
	SessionStore sessions.Store
	SessionName  string
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

	formDecoder := form.NewDecoder()
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
	queryDecoder := form.NewDecoder()
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)

	return func(w http.ResponseWriter, r *http.Request) {
		c := &Ctx[V]{
			Log:          config.Log,
			Res:          w,
			Req:          r,
			path:         mux.Vars(r),
			CSRFToken:    csrf.Token(r),
			CSRFTag:      csrf.TemplateField(r),
			router:       config.Router,
			formDecoder:  formDecoder,
			queryDecoder: queryDecoder,
		}
		session, err := config.SessionStore.Get(r, config.SessionName)
		if err != nil {
			config.ErrorHandler(c, err)
			return
		}
		c.Session = &gorillaSession{
			req:     r,
			res:     w,
			session: session,
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
	Log          *zap.SugaredLogger
	Req          *http.Request
	Res          http.ResponseWriter
	Session      Session
	path         map[string]string
	CSRFToken    string
	CSRFTag      template.HTML
	Var          V
	router       *mux.Router
	formDecoder  *form.Decoder
	queryDecoder *form.Decoder
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

func (c *Ctx[V]) Redirect(route string, pairs ...string) {
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

// TODO return httperror, embed err in httperror
func (c *Ctx[V]) Bind(v any, flags ...Flag) error {
	m := c.Req.Method
	if m == http.MethodGet || m == http.MethodDelete || m == http.MethodHead {
		return c.BindQuery(v, flags...)
	}
	return c.BindForm(v, flags...)
}

func (c *Ctx[V]) BindQuery(v any, flags ...Flag) error {
	vals := c.Req.URL.Query()
	if hasFlag(flags, Vacuum) {
		vacuum(vals)
	}
	return c.queryDecoder.Decode(v, vals)
}

func (c *Ctx[V]) BindForm(v any, flags ...Flag) error {
	c.Req.ParseForm()
	vals := c.Req.Form
	if hasFlag(flags, Vacuum) {
		vacuum(vals)
	}
	return c.formDecoder.Decode(v, vals)
}

func vacuum(values url.Values) {
	for key, vals := range values {
		var tmp []string
		for _, val := range vals {
			val = strings.TrimSpace(val)
			if val != "" {
				tmp = append(tmp, val)
			}
		}
		if len(tmp) > 0 {
			values[key] = tmp
		} else {
			delete(values, key)
		}
	}
}

func hasFlag(flags []Flag, flag Flag) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}
