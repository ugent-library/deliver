package crumb

import (
	"context"
	"net/http"

	"github.com/felixge/httpsnoop"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var cookiesKey = contextKey("cookies")

func Cookies(r *http.Request) *CookieJar {
	if v := r.Context().Value(cookiesKey); v != nil {
		return v.(*CookieJar)
	}
	return nil
}

type options struct {
	errorHandler func(error)
}

type Option func(*options)

func WithErrorHandler(fn func(error)) Option {
	return func(opts *options) {
		opts.errorHandler = fn
	}
}

func Enable(opts ...Option) func(http.Handler) http.Handler {
	o := &options{
		errorHandler: func(error) {},
	}
	for _, opt := range opts {
		opt(o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jar := New(r.Cookies())
			written := false

			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(nextFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						if !written {
							if err := jar.Write(w); err != nil {
								o.errorHandler(err)
							}
							written = true
						}
						nextFunc(code)
					}
				},
				Write: func(nextFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						if !written {
							if err := jar.Write(w); err != nil {
								o.errorHandler(err)
							}
							written = true
						}
						return nextFunc(b)
					}
				},
			})

			c := context.WithValue(r.Context(), cookiesKey, jar)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
