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

// TODO reduce type requirements
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
	Banner        string
}

func Get(r *http.Request) *Ctx {
	return r.Context().Value(ctxKey).(*Ctx)
}

func Set(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := New(config, w, r)

			r = r.WithContext(context.WithValue(r.Context(), ctxKey, c))

			if err := c.init(w, r, config.Repo.Users); err != nil {
				c.HandleError(w, r, err)
				return
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
	Repo          *repositories.Repo
	Storage       objectstore.Store
	MaxFileSize   int64
	Auth          *oidc.Auth
	host          string
	scheme        string
	errorHandlers map[int]http.HandlerFunc
	router        *ich.Mux
	assets        mix.Manifest
	Log           *zap.SugaredLogger // TODO use plain logger
	Hub           *htmx.Hub
	CSRFToken     string
	CSRFTag       string
	User          *models.User
	Flash         []Flash
	Banner        string
	*models.Permissions
}

func New(config Config, w http.ResponseWriter, r *http.Request) *Ctx {
	c := &Ctx{
		Repo:          config.Repo,
		Storage:       config.Storage,
		MaxFileSize:   config.MaxFileSize,
		Auth:          config.Auth,
		host:          r.Host,
		scheme:        r.URL.Scheme,
		errorHandlers: config.ErrorHandlers,
		Log:           zaphttp.Logger(r.Context()).Sugar(),
		CSRFToken:     csrf.Token(r),
		CSRFTag:       string(csrf.TemplateField(r)),
		router:        config.Router,
		assets:        config.Assets,
		Hub:           config.Hub,
		Banner:        config.Banner,
		Permissions:   config.Permissions,
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

	if h, ok := c.errorHandlers[httpErr.StatusCode]; ok {
		h(w, r)
		return
	}

	c.Log.Error(err)

	http.Error(w, http.StatusText(httpErr.StatusCode), httpErr.StatusCode)
}

func (c *Ctx) init(w http.ResponseWriter, r *http.Request, users *repositories.UsersRepo) error {
	// remember token cookie
	if cookie, _ := r.Cookie(RememberCookie); cookie != nil {
		user, err := users.GetByRememberToken(r.Context(), cookie.Value)
		if err != nil && err != models.ErrNotFound {
			return err
		}
		c.User = user
	}

	// flash cookies
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

	return nil
}

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	return c.router.PathTo(name, pairs...)
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	u := c.router.PathTo(name, pairs...)
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
