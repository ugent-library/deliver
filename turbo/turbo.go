package turbo

import (
	"net/http"
	"strings"
)

const ContentType = "text/vnd.turbo-stream.html; charset=utf-8"

func Request(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Accept"), "text/vnd.turbo-stream.html")
}

// func FrameRequestID(r *http.Request) string {
// 	return r.Header.Get("Turbo-Frame")
// }

// func FrameRequest(r *http.Request) bool {
// 	return FrameRequestID(r) != ""
// }
