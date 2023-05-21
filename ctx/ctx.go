package ctx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/nics/ich"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/zaphttp"
	"go.uber.org/zap"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ctxKey = contextKey("ctx")

// TODO set __Host- cookie prefix in production
const (
	RememberCookie = "deliver.remember"
	FlashCookie    = "deliver.flash"
)

// TODO reduce type requirements
type Config struct {
	GetUserByRememberToken func(context.Context, string) (*models.User, error)
	Router                 *ich.Mux
	ErrorHandlers          map[int]http.HandlerFunc
	Permissions            *models.Permissions
	Assets                 mix.Manifest
	Hub                    *htmx.Hub
	Banner                 string
}

func Get(ctx context.Context) *Ctx {
	return ctx.Value(ctxKey).(*Ctx)
}

func Set(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := New(config, w, r)
			if err := c.loadSession(config.GetUserByRememberToken); err != nil {
				c.HandleError(err)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type Flash struct {
	Type         string
	Title        string
	Body         string
	DismissAfter time.Duration
}

type Ctx struct {
	w             http.ResponseWriter
	r             *http.Request
	errorHandlers map[int]http.HandlerFunc
	router        *ich.Mux
	assets        mix.Manifest
	Log           *zap.SugaredLogger // TODO use plain logger
	Hub           *htmx.Hub
	CSRFToken     string
	CSRFTag       string
	Cookies       *crumb.CookieJar
	User          *models.User
	Flash         []Flash
	Banner        string
	*models.Permissions
}

func New(config Config, w http.ResponseWriter, r *http.Request) *Ctx {
	c := &Ctx{
		w:             w,
		r:             r,
		errorHandlers: config.ErrorHandlers,
		Log:           zaphttp.Logger(r.Context()).Sugar(),
		CSRFToken:     csrf.Token(r),
		CSRFTag:       string(csrf.TemplateField(r)),
		Cookies:       crumb.Cookies(r),
		router:        config.Router,
		assets:        config.Assets,
		Hub:           config.Hub,
		Banner:        config.Banner,
		Permissions:   config.Permissions,
	}

	return c
}

func (c *Ctx) HandleError(err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	if h, ok := c.errorHandlers[httpErr.StatusCode]; ok {
		h(c.w, c.r)
		return
	}

	c.Log.Error(err)
	http.Error(c.w, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
}

func (c *Ctx) loadSession(userSource func(context.Context, string) (*models.User, error)) error {
	if token := c.Cookies.Get(RememberCookie); token != "" {
		user, err := userSource(c.r.Context(), token)
		if err != nil && err != models.ErrNotFound {
			return err
		}
		c.User = user
	}

	c.Cookies.Unmarshal(FlashCookie, &c.Flash)
	c.Cookies.Delete(FlashCookie)

	return nil
}

func (c *Ctx) PathParam(param string) string {
	return chi.URLParamFromCtx(c.r.Context(), param)
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	u := c.router.PathTo(name, pairs...)
	u.Host = c.r.Host
	u.Scheme = c.r.URL.Scheme
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u
}

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	return c.router.PathTo(name, pairs...)
}

func (c *Ctx) PersistFlash(f Flash) {
	c.Cookies.Append(FlashCookie, f, time.Now().Add(3*time.Minute))
}

func (c *Ctx) AssetPath(asset string) string {
	ap, err := c.assets.AssetPath(asset)
	if err != nil {
		panic(err)
	}
	return ap
}

func (c *Ctx) WebSocketPath(channels ...string) string {
	h, err := c.Hub.EncryptChannelNames(channels)
	if err != nil {
		c.Log.Error(err)
		return ""
	}
	return "/ws?channels=" + url.QueryEscape(h)
}
