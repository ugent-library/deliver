package controllers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ugent-library/dilliver/autosession"
	"github.com/ugent-library/dilliver/httperror"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/routes"
	"github.com/ugent-library/dilliver/zaphttp"
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
	Session autosession.Session
	User    *models.User
	Flash   []Flash
}

type Flash struct {
	Type  string
	Title string
	Body  template.HTML
}

type Renderer interface {
	Render(http.ResponseWriter, any) error
}

type ViewData struct {
	routes.URLHelpers
	CSRFToken string
	CSRFTag   template.HTML
	User      *models.User
	Flash     []Flash
	Data      any
}

func Wrapper(router *mux.Router, errorHandler func(*Ctx, error)) func(...func(*Ctx) error) http.Handler {
	return func(handlers ...func(*Ctx) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &Ctx{
				HandlerHelpers: routes.NewHandlerHelpers(router, w, r),
				URLHelpers:     routes.NewURLHelpers(router, r),
				Log:            zaphttp.Logger(r.Context()).Sugar(),
				Res:            w,
				Req:            r,
				Session:        autosession.Get(r),
			}
			if err := LoadSession(c); err != nil {
				errorHandler(c, err)
				return
			}
			for _, fn := range handlers {
				if err := fn(c); err != nil {
					errorHandler(c, err)
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
		URLHelpers: c.URLHelpers,
		CSRFToken:  csrf.Token(c.Req),
		CSRFTag:    csrf.TemplateField(c.Req),
		User:       c.User,
		Flash:      c.Flash,
		Data:       data,
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
