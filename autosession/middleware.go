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

func Get(r *http.Request) Session {
	if l := r.Context().Value(sessionKey); l != nil {
		return l.(Session)
	}
	return nil
}

func Enable(fn SessionFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := fn(w, r)
			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(nextFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						// TODO handle error
						s.Save()
						nextFunc(code)
					}
				},
				Write: func(nextFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						// TODO handle error
						s.Save()
						return nextFunc(b)
					}
				},
			})
			c := context.WithValue(r.Context(), sessionKey, s)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
