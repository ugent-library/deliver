package controllers

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ugent-library/deliver/autosession"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/routes"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/zaphttp"
	"go.uber.org/zap"
)

const (
	userKey  = "user"
	flashKey = "flash"

	infoFlash = "info"
)

type Map = map[string]any

type Ctx struct {
	routes.HandlerHelpers
	routes.URLHelpers
	Log     *zap.SugaredLogger // TODO use plain logger
	Req     *http.Request
	Res     http.ResponseWriter
	Session *autosession.Session
	User    *models.User
	*models.Permissions
	Flash []Flash
}

type Flash struct {
	Type         string
	Title        string
	Body         template.HTML
	DismissAfter time.Duration
}

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type ViewData struct {
	routes.URLHelpers
	CSRFToken string
	CSRFTag   template.HTML
	User      *models.User
	*models.Permissions
	Flash []Flash
	Data  any
}

type Config struct {
	Router       *mux.Router
	ErrorHandler func(*Ctx, error)
	Permissions  *models.Permissions
}

// TODO add Ctx as request Context value in middleware?
func Wrapper(config Config) func(...func(*Ctx) error) http.Handler {
	return func(handlers ...func(*Ctx) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &Ctx{
				HandlerHelpers: routes.NewHandlerHelpers(config.Router, w, r),
				URLHelpers:     routes.NewURLHelpers(config.Router, r),
				Log:            zaphttp.Logger(r.Context()).Sugar(),
				Res:            w,
				Req:            r,
				Session:        autosession.Get(r),
				Permissions:    config.Permissions,
			}
			if err := LoadSession(c); err != nil {
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

func (c *Ctx) Render(r Renderer, data any) error {
	return r.Render(c.Res, ViewData{
		URLHelpers:  c.URLHelpers,
		CSRFToken:   csrf.Token(c.Req),
		CSRFTag:     csrf.TemplateField(c.Req),
		User:        c.User,
		Permissions: c.Permissions,
		Flash:       c.Flash,
		Data:        data,
	})
}

func LoadSession(c *Ctx) error {
	if val := c.Session.Get(userKey); val != nil {
		c.User = val.(*models.User)
	}
	if vals := c.Session.Pop(flashKey); vals != nil {
		flash := make([]Flash, len(vals.([]any)))
		for i, v := range vals.([]any) {
			flash[i] = v.(Flash)
		}
		c.Flash = flash
	}
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
