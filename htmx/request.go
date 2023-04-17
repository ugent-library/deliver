package htmx

import "net/http"

func Request(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
