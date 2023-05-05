package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/zaphttp"
)

// TODO set __Host- cookie prefix in production
const (
	rememberCookie = "deliver.remember"
	flashCookie    = "deliver.flash"
)

type Map = map[string]any

type Config struct {
	UserFunc     func(context.Context, string) (*models.User, error)
	Router       *mux.Router
	ErrorHandler func(http.ResponseWriter, *http.Request, *ctx.Ctx, error)
	Permissions  *models.Permissions
	Assets       mix.Manifest
	Hub          *htmx.Hub
}

type Handler func(http.ResponseWriter, *http.Request, *ctx.Ctx) error

// TODO add Ctx as request Context value in middleware?
func Wrapper(config Config) func(...Handler) http.Handler {
	return func(handlers ...Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &ctx.Ctx{
				Log:         zaphttp.Logger(r.Context()).Sugar(),
				Res:         w,
				Req:         r,
				CSRFToken:   csrf.Token(r),
				CSRFTag:     string(csrf.TemplateField(r)),
				Cookies:     crumb.Cookies(r),
				Permissions: config.Permissions,
				Router:      config.Router,
				PathVars:    mux.Vars(r),
				Assets:      config.Assets,
				Hub:         config.Hub,
			}
			if err := LoadSession(config.UserFunc, c); err != nil {
				config.ErrorHandler(w, r, c, err)
				return
			}
			for _, fn := range handlers {
				if err := fn(w, r, c); err != nil {
					config.ErrorHandler(w, r, c, err)
					return
				}
			}
		})
	}
}

func LoadSession(userFunc func(context.Context, string) (*models.User, error), c *ctx.Ctx) error {
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

func RequireUser(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	if c.User == nil {
		return httperror.Unauthorized
	}
	return nil
}

func RequireAdmin(w http.ResponseWriter, r *http.Request, c *ctx.Ctx) error {
	if c.User == nil {
		return httperror.Unauthorized
	}
	if !c.IsAdmin(c.User) {
		return httperror.Forbidden
	}
	return nil
}
