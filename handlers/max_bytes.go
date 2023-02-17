package handlers

import "net/http"

func MaxBytesHandler(size int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.MaxBytesHandler(next, size)
	}
}
