package middleware

import (
	"net/http"
)

func Recover(fn func(any)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// don't recover from http.ErrAbortHandler (see http.ErrAbortHandler docs)
					if err == http.ErrAbortHandler {
						panic(err)
					}

					fn(err)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
