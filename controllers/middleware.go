package controllers

import (
	"net/http"

	"github.com/ugent-library/dilliver/ulid"
)

func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			r.Header.Set("X-Request-ID", ulid.MustGenerate())
		}
		next.ServeHTTP(w, r)
	})
}
