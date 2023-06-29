package ctx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/nics/ich"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
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
	RememberCookie    = "deliver.remember"
	FlashCookiePrefix = "deliver.flash."
)

type Config struct {
	Repo          *repositories.Repo
	Storage       objectstore.Store
	MaxFileSize   int64
	Auth          *oidc.Auth
	Router        *ich.Mux
	ErrorHandlers map[int]http.HandlerFunc
	Permissions   *models.Permissions
	Assets        mix.Manifest
	Hub           *htmx.Hub
	Env           string
}

func Get(r *http.Request) *Ctx {
	return r.Context().Value(ctxKey).(*Ctx)
}

func Set(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := New(config, w, r)

			r = r.WithContext(context.WithValue(r.Context(), ctxKey, c))

			// load user from remember token cookie
			if cookie, _ := r.Cookie(RememberCookie); cookie != nil {
				user, err := config.Repo.Users.GetByRememberToken(r.Context(), cookie.Value)
				if err != nil && err != models.ErrNotFound {
					c.HandleError(w, r, err)
					return
				}
				c.User = user
			}

			// load flash from cookies
			for _, cookie := range r.Cookies() {
				if !strings.HasPrefix(cookie.Name, FlashCookiePrefix) {
					continue
				}

				// delete after read
				http.SetCookie(w, &http.Cookie{
					Name:     cookie.Name,
					Value:    "",
					Expires:  time.Now(),
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})

				j, err := base64.URLEncoding.DecodeString(cookie.Value)
				if err != nil {
					continue
				}
				f := Flash{}
				if err = json.Unmarshal(j, &f); err == nil {
					c.Flash = append(c.Flash, f)
				}
			}

			next.ServeHTTP(w, r)
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
	Config
	host      string
	scheme    string
	Log       *zap.SugaredLogger
	CSRFToken string
	CSRFTag   string
	User      *models.User
	Flash     []Flash
	Env       string
}

func New(config Config, w http.ResponseWriter, r *http.Request) *Ctx {
	c := &Ctx{
		Config:    config,
		host:      r.Host,
		scheme:    r.URL.Scheme,
		Log:       zaphttp.Logger(r.Context()).Sugar(),
		CSRFToken: csrf.Token(r),
		CSRFTag:   string(csrf.TemplateField(r)),
		Env:       config.Env,
	}
	if c.scheme == "" {
		c.scheme = "http"
	}

	return c
}

func (c *Ctx) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == models.ErrNotFound {
		err = httperror.NotFound
	}

	var httpErr *httperror.Error
	if !errors.As(err, &httpErr) {
		httpErr = httperror.InternalServerError
	}

	if h, ok := c.ErrorHandlers[httpErr.StatusCode]; ok {
		h(w, r)
		return
	}

	c.Log.Error(err)

	http.Error(w, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
}

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	return c.Router.PathTo(name, pairs...)
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	u := c.Router.PathTo(name, pairs...)
	u.Host = c.host
	u.Scheme = c.scheme
	return u
}

func (c *Ctx) PersistFlash(w http.ResponseWriter, f Flash) {
	j, err := json.Marshal(f)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookiePrefix + ulid.Make().String(),
		Value:    base64.URLEncoding.EncodeToString(j),
		Expires:  time.Now().Add(3 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (c *Ctx) AssetPath(asset string) string {
	ap, err := c.Assets.AssetPath(asset)
	if err != nil {
		panic(err)
	}
	return ap
}

func (c *Ctx) WebSocketPath(channels ...string) string {
	h, err := c.Hub.EncryptChannelNames(channels)
	if err != nil {
		panic(err)
	}
	return "/ws?channels=" + url.QueryEscape(h)
}
