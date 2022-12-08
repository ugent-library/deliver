package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func Wrapper(c Config) func(func(http.ResponseWriter, *http.Request, Ctx)) http.HandlerFunc {
	return func(fn func(http.ResponseWriter, *http.Request, Ctx)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := Ctx{c}
			fn(w, r, ctx)
		}
	}
}

type Config struct {
	Router       *mux.Router
	SessionName  string
	SessionStore sessions.Store
}

type Ctx struct {
	Config
}

func (c Ctx) URL(route string, pairs ...string) *url.URL {
	r := c.Router.Get(route)
	if r == nil {
		panic(fmt.Errorf("route '%s' not found", route))
	}
	u, err := r.URL(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

func (c Ctx) URLPath(route string, pairs ...string) *url.URL {
	r := c.Router.Get(route)
	if r == nil {
		panic(fmt.Errorf("route '%s' not found", route))
	}
	u, err := r.URLPath(pairs...)
	if err != nil {
		panic(fmt.Errorf("can't reverse route '%s': %w", route, err))
	}
	return u
}

//////////////////////////

// type InnerCtx struct {
// 	Ctx
// 	Inner bool
// }

// func WrapInnerCtx(fn func(http.ResponseWriter, *http.Request, InnerCtx)) http.HandlerFunc {
// 	return WithCtx(w, r, func(w http.ResponseWriter, r *http.Request, ctx Ctx) {
// 		ictx := InnerCtx{
// 			Ctx:     ctx,
// 			Inner: true,
// 		}
// 		fn(w, r, ictx)
// 	})
// }
