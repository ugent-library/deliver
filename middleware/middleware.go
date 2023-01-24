// Pakage middleware contains some generic middlewares and helper functions to
// make composing middleware more readable.
package middleware

import "net/http"

// Apply wraps a handler with middlewares. Middleware is applied in the order it
// is given.
//
//	handler = middleware.Apply(handler,
//	  middleware1,
//	  middleware2,
//	  middleware3,
//	)
//
// is equivalant to:
//
//	handler = middleware1(middleware2(middleware3(handler)))
func Apply(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// If only applies a middleware if cond is true.
//
//	handler = middleware.Apply(handler,
//	  middleware1,
//	  middleware.If(isProduction, middleware2),
//	  middleware3,
//	)
func If(cond bool, mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	if cond {
		return mw
	}
	return func(h http.Handler) http.Handler {
		return h
	}
}
