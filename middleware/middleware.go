package middleware

import "net/http"

func Apply(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func If(cond bool, mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	if cond {
		return mw
	}
	return func(h http.Handler) http.Handler {
		return h
	}
}
