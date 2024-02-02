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
	"github.com/ugent-library/catbird"
	"github.com/ugent-library/crypt"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

// TODO set __Host- cookie prefix in production
const (
	RememberCookie    = "deliver.remember"
	FlashCookiePrefix = "deliver.flash."
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ctxKey = contextKey("ctx")

func Get(r *http.Request) *Ctx {
	return r.Context().Value(ctxKey).(*Ctx)
}

func Set(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &Ctx{
				Config:    config,
				host:      r.Host,
				scheme:    r.URL.Scheme,
				Log:       zaphttp.Logger(r.Context()).Sugar(),
				CSRFToken: csrf.Token(r),
				CSPNonce:  secure.CSPNonce(r.Context()),
			}
			if c.scheme == "" {
				if config.Env == "local" {
					c.scheme = "http"
				} else {
					c.scheme = "https"
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), ctxKey, c))

			// load user from remember token cookie
			u, err := getUser(r, config.Repo)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}
			c.User = u

			// load flash from cookies
			f, err := getFlash(r, w)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}
			c.Flash = f

			next.ServeHTTP(w, r)
		})
	}
}

type Config struct {
	*crypt.Crypt
	Env           string
	Repo          *repositories.Repo
	Storage       objectstores.Store
	MaxFileSize   int64
	Auth          *oidc.Auth
	Router        *ich.Mux
	ErrorHandlers map[int]http.HandlerFunc
	Permissions   *models.Permissions
	Assets        mix.Manifest
	Hub           *catbird.Hub
	Timezone      *time.Location
	CSRFName      string
}

type Flash struct {
	Type         string
	Title        string
	Body         string
	DismissAfter time.Duration
	AlwaysShow   bool
}

type Ctx struct {
	Config
	host      string
	scheme    string
	Log       *zap.SugaredLogger
	CSRFToken string
	CSPNonce  string
	User      *models.User
	Flash     []Flash
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

func getUser(r *http.Request, repo *repositories.Repo) (u *models.User, err error) {
	if cookie, _ := r.Cookie(RememberCookie); cookie != nil {
		u, err = repo.Users.GetByRememberToken(r.Context(), cookie.Value)
		if err == models.ErrNotFound {
			err = nil
		}
	}
	return
}

func getFlash(r *http.Request, w http.ResponseWriter) ([]Flash, error) {
	var flashes []Flash

	for _, cookie := range r.Cookies() {
		if !strings.HasPrefix(cookie.Name, FlashCookiePrefix) {
			continue
		}

		// delete cookie
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
			return nil, err
		}

		f := Flash{}
		if err = json.Unmarshal(j, &f); err != nil {
			return nil, err
		}
		flashes = append(flashes, f)
	}

	return flashes, nil
}

func (c *Ctx) AssetPath(asset string) string {
	ap, err := c.Assets.AssetPath(asset)
	if err != nil {
		panic(err)
	}
	return ap
}

func (c *Ctx) WebSocketPath(topics ...string) string {
	token, err := c.EncryptValue(topics)
	if err != nil {
		panic(err)
	}
	return c.PathTo("ws", "token", token).String()
}
