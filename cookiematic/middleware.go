package cookiematic

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

func Cookies(r *http.Request) *Jar {
	if v := r.Context().Value(cookiesKey); v != nil {
		return v.(*Jar)
	}
	return nil
}

// futureproofing, not used yet
type options struct{}
type Option func(*options)

// TODO error handler
func Enable(opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			j := NewJar(w, r)

			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(nextFunc httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						j.Write()
						nextFunc(code)
					}
				},
				Write: func(nextFunc httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(b []byte) (int, error) {
						j.Write()
						return nextFunc(b)
					}
				},
			})
			c := context.WithValue(r.Context(), cookiesKey, j)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}
