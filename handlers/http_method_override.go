package handlers

import "net/http"

const (
	HTTPMethodOverrideFormKey = "_method"
	HTTPMethodOverrideHeader  = "X-HTTP-Method-Override"
)

/*
	Copied from gorilla's HTTPMethodOverrideHandler,
	but it inverses the internal logic: read headers
	first, then form.
	Reading the form first triggers a body read,
	which is not what we want
*/

func HTTPMethodOverrideHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			om := r.Header.Get(HTTPMethodOverrideHeader)
			if om == "" {
				om = r.FormValue(HTTPMethodOverrideFormKey)
			}
			if om == "PUT" || om == "PATCH" || om == "DELETE" {
				r.Method = om
			}
		}
		next.ServeHTTP(w, r)
	})
}
