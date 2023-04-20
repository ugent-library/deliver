package htmx

import (
	"encoding/json"
	"net/http"
)

func Request(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func AddTrigger[T string | map[string]any](w http.ResponseWriter, t T) {
	switch v := any(t).(type) {
	case string:
		w.Header().Set("HX-Trigger", v)
	case map[string]any:
		j, _ := json.Marshal(v)
		w.Header().Set("HX-Trigger", string(j))
	}
}

// type LocationConfig struct {
// 	Path    string            `json:"path,omitempty"`
// 	Source  string            `json:"source,omitempty"`
// 	Event   string            `json:"event,omitempty"`
// 	Handler string            `json:"handler,omitempty"`
// 	Target  string            `json:"target,omitempty"`
// 	Swap    string            `json:"swap,omitempty"`
// 	Values  map[string]string `json:"values,omitempty"`
// 	Headers map[string]string `json:"headers,omitempty"`
// }

// func Location[T string | LocationConfig](w http.ResponseWriter, t T) {
// 	switch v := any(t).(type) {
// 	case string:
// 		w.Header().Set("HX-Location", v)
// 	case LocationConfig:
// 		j, _ := json.Marshal(v)
// 		w.Header().Set("HX-Location", string(j))
// 	}
// }

// func PushURL(w http.ResponseWriter, v string) {
// 	w.Header().Set("HX-Push-Url", v)
// }

// func Redirect(w http.ResponseWriter, v string) {
// 	w.Header().Set("HX-Redirect", v)
// }

// func Refresh(w http.ResponseWriter) {
// 	w.Header().Set("HX-Refresh", "true")
// }

// func ReplaceURL(w http.ResponseWriter, v string) {
// 	w.Header().Set("HX-Replace-Url", v)
// }
