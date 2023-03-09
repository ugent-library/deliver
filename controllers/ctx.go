package controllers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/zaphttp"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

// TODO set __Host- cookie prefix in production
const (
	rememberCookie = "deliver.remember"
	flashCookie    = "deliver.flash"
)

type Map = map[string]any

type Ctx struct {
	Log     *zap.SugaredLogger // TODO use plain logger
	Req     *http.Request
	Res     http.ResponseWriter
	Cookies *crumb.CookieJar
	User    *models.User
	*models.Permissions
	Flash  []Flash
	router *mux.Router
	path   map[string]string
	render *render.Render
}

type Flash struct {
	Type         string
	Title        string
	Body         template.HTML
	DismissAfter time.Duration
}

type TemplateData struct {
	ctx       *Ctx
	CSRFToken string
	CSRFTag   template.HTML
	Data      any
}

func (t TemplateData) User() *models.User {
	return t.ctx.User
}

func (t TemplateData) Flash() []Flash {
	return t.ctx.Flash
}

func (t TemplateData) URLTo(name string, pairs ...string) *url.URL {
	return t.ctx.URLTo(name, pairs...)
}

func (t TemplateData) PathTo(name string, pairs ...string) *url.URL {
	return t.ctx.PathTo(name, pairs...)
}

func (t TemplateData) IsAdmin(user *models.User) bool {
	return t.ctx.IsAdmin(user)
}

func (t TemplateData) IsSpaceAdmin(user *models.User, space *models.Space) bool {
	return t.ctx.IsSpaceAdmin(user, space)
}

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type Config struct {
	UserFunc     func(context.Context, string) (*models.User, error)
	Router       *mux.Router
	ErrorHandler func(*Ctx, error)
	Permissions  *models.Permissions
	Render       *render.Render
}

// TODO add Ctx as request Context value in middleware?
func Wrapper(config Config) func(...func(*Ctx) error) http.Handler {
	return func(handlers ...func(*Ctx) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &Ctx{
				Log:         zaphttp.Logger(r.Context()).Sugar(),
				Res:         w,
				Req:         r,
				Cookies:     crumb.Cookies(r),
				Permissions: config.Permissions,
				router:      config.Router,
				path:        mux.Vars(r),
				render:      config.Render,
			}
			if err := LoadSession(config.UserFunc, c); err != nil {
				config.ErrorHandler(c, err)
				return
			}
			for _, fn := range handlers {
				if err := fn(c); err != nil {
					config.ErrorHandler(c, err)
					return
				}
			}
		})
	}
}

func (c *Ctx) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx) HTML(status int, layout, tmpl string, data any) error {
	return c.render.HTML(c.Res, status, tmpl, TemplateData{
		ctx:       c,
		CSRFToken: csrf.Token(c.Req),
		CSRFTag:   csrf.TemplateField(c.Req),
		Data:      data,
	}, render.HTMLOptions{
		Layout: layout,
	})
}

func (c *Ctx) Path(param string) string {
	return c.path[param]
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	route := c.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URL(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
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

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	route := c.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
	}
	return u
}

func (c *Ctx) RedirectTo(name string, pairs ...string) {
	route := c.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
	}
	http.Redirect(c.Res, c.Req, u.String(), http.StatusSeeOther)
}

func (c *Ctx) AddFlash(f Flash) {
	c.Cookies.Append(flashCookie, f, time.Now().Add(3*time.Minute))
}

func LoadSession(userFunc func(context.Context, string) (*models.User, error), c *Ctx) error {
	if token := c.Cookies.Get(rememberCookie); token != "" {
		user, err := userFunc(c.Context(), token)
		if err != nil && err != models.ErrNotFound {
			return err
		}
		c.User = user
	}

	c.Cookies.Unmarshal(flashCookie, &c.Flash)
	c.Cookies.Delete(flashCookie)

	return nil
}

func RequireUser(c *Ctx) error {
	if c.User == nil {
		return httperror.Unauthorized
	}
	return nil
}

func RequireAdmin(c *Ctx) error {
	if c.User == nil {
		return httperror.Unauthorized
	}
	if !c.IsAdmin(c.User) {
		return httperror.Forbidden
	}
	return nil
}
