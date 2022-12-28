package bind

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
)

type Flag int

const (
	Vacuum Flag = iota
)

var (
	queryDecoder = form.NewDecoder()
	formDecoder  = form.NewDecoder()
)

func init() {
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
}

func Request(r *http.Request, v any, flags ...Flag) error {
	if r.Method == http.MethodGet || r.Method == http.MethodDelete || r.Method == http.MethodHead {
		return Query(r, v, flags...)
	}
	return Form(r, v, flags...)
}

func Query(r *http.Request, v any, flags ...Flag) error {
	vals := r.URL.Query()
	if hasFlag(flags, Vacuum) {
		vacuum(vals)
	}
	return queryDecoder.Decode(v, vals)
}

func Form(r *http.Request, v any, flags ...Flag) error {
	r.ParseForm()
	vals := r.Form
	if hasFlag(flags, Vacuum) {
		vacuum(vals)
	}
	return formDecoder.Decode(v, vals)
}

func vacuum(values url.Values) {
	for key, vals := range values {
		var tmp []string
		for _, val := range vals {
			val = strings.TrimSpace(val)
			if val != "" {
				tmp = append(tmp, val)
			}
		}
		if len(tmp) > 0 {
			values[key] = tmp
		} else {
			delete(values, key)
		}
	}
}

func hasFlag(flags []Flag, flag Flag) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}
