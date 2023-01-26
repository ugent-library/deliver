package autosession

import (
	"context"
	"net/http"

	"github.com/felixge/httpsnoop"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var sessionKey = contextKey("session")

func Get(r *http.Request) *Session {
	if l := r.Context().Value(sessionKey); l != nil {
		return l.(*Session)
	}
	return nil
}

type Option func(*options)

type options struct {
	errorHandler func(error)
}

func WithErrorHandler(fn func(error)) Option {
	return func(opts *options) {
		opts.errorHandler = fn
	}
}

// TODO enable multiple sessions? maybe using sessions.Registry?
// TODO stop on error?
func Enable(provider SessionProvider, opts ...Option) func(http.Handler) http.Handler {
	o := &options{
		errorHandler: func(error) {},
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := provider(w, r)
			if err != nil {
				o.errorHandler(err)
			}

			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(nextFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						if err := s.Save(r.Context()); err != nil {
							o.errorHandler(err)
						}
						nextFunc(code)
					}
				},
				Write: func(nextFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						if err := s.Save(r.Context()); err != nil {
							o.errorHandler(err)
						}
						return nextFunc(b)
					}
				},
			})
			c := context.WithValue(r.Context(), sessionKey, s)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
