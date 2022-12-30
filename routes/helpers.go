package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type URLHelpers struct {
	router *mux.Router
	r      *http.Request
}

func NewURLHelpers(router *mux.Router, r *http.Request) URLHelpers {
	return URLHelpers{
		router: router,
		r:      r,
	}
}

func (h URLHelpers) URLTo(name string, pairs ...string) *url.URL {
	route := h.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URL(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
	}
	if u.Host == "" {
		u.Host = h.r.Host
	}
	if u.Scheme == "" {
		u.Scheme = h.r.URL.Scheme
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u
}

func (h URLHelpers) PathTo(name string, pairs ...string) *url.URL {
	route := h.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
	}
	return u
}

type HandlerHelpers struct {
	router *mux.Router
	path   map[string]string
	w      http.ResponseWriter
	r      *http.Request
}

func NewHandlerHelpers(router *mux.Router, w http.ResponseWriter, r *http.Request) HandlerHelpers {
	return HandlerHelpers{
		router: router,
		path:   mux.Vars(r),
		w:      w,
		r:      r,
	}
}

func (h HandlerHelpers) Path(param string) string {
	return h.path[param]
}

func (h HandlerHelpers) ExecuteHandler(name string) {
	route := h.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	route.GetHandler().ServeHTTP(h.w, h.r)
}

func (h HandlerHelpers) RedirectTo(name string, pairs ...string) {
	route := h.router.Get(name)
	if route == nil {
		panic(fmt.Errorf("unknown route '%s'", name))
	}
	u, err := route.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", name, err))
	}
	http.Redirect(h.w, h.r, u.String(), http.StatusSeeOther)
}
