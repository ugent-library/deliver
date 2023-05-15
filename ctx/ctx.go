package ctx

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/nics/ich"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/mix"
	"go.uber.org/zap"
)

// TODO set __Host- cookie prefix in production
const (
	rememberCookie = "deliver.remember"
	flashCookie    = "deliver.flash"
)

type Ctx struct {
	Log       *zap.SugaredLogger // TODO use plain logger
	Req       *http.Request
	Res       http.ResponseWriter
	CSRFToken string
	CSRFTag   string
	Cookies   *crumb.CookieJar
	User      *models.User
	*models.Permissions
	Flash    []Flash
	Router   *ich.Mux
	PathVars map[string]string
	Assets   mix.Manifest
	Hub      *htmx.Hub
	Banner   string
}

type Flash struct {
	Type         string
	Title        string
	Body         string
	DismissAfter time.Duration
}

func (c *Ctx) Context() context.Context {
	return c.Req.Context()
}

func (c *Ctx) HTML(status int, body string) error {
	if hdr := c.Res.Header(); hdr.Get("Content-Type") == "" {
		hdr.Set("Content-Type", "text/html; charset=utf-8")
	}
	c.Res.WriteHeader(status)
	_, err := c.Res.Write([]byte(body))
	return err
}

func (c *Ctx) Path(param string) string {
	return c.PathVars[param]
}

func (c *Ctx) URLTo(name string, pairs ...string) *url.URL {
	u := c.Router.PathTo(name, pairs...)
	u.Host = c.Req.Host
	u.Scheme = c.Req.URL.Scheme
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u
}

func (c *Ctx) PathTo(name string, pairs ...string) *url.URL {
	u := c.Router.PathTo(name, pairs...)
	return u
}

func (c *Ctx) RedirectTo(name string, pairs ...string) {
	u := c.Router.PathTo(name, pairs...)
	http.Redirect(c.Res, c.Req, u.String(), http.StatusSeeOther)
}

func (c *Ctx) AssetPath(asset string) string {
	ap, err := c.Assets.AssetPath(asset)
	if err != nil {
		panic(err)
	}
	return ap
}

func (c *Ctx) AddFlash(f Flash) {
	c.Cookies.Append(flashCookie, f, time.Now().Add(3*time.Minute))
}

func (c *Ctx) WebSocketPath(channels ...string) string {
	h, err := c.Hub.EncryptChannelNames(channels)
	if err != nil {
		c.Log.Error(err)
		return ""
	}
	return "/ws?channels=" + url.QueryEscape(h)
}
