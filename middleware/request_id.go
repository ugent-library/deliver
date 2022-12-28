package middleware

import (
	"net/http"
)

func SetRequestID(generator func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Request-ID") == "" {
				r.Header.Set("X-Request-ID", generator())
			}
			next.ServeHTTP(w, r)
		})
	}
}
