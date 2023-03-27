package turbo

import "net/http"

func FrameRequest(r *http.Request) bool {
	return r.Header.Get("Turbo-Frame") != ""
}
