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
	formDecoder   = form.NewDecoder()
	queryDecoder  = form.NewDecoder()
	headerDecoder = form.NewDecoder()
	formEncoder   = form.NewEncoder()
	queryEncoder  = form.NewEncoder()
	headerEncoder = form.NewEncoder()
)

func init() {
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	headerDecoder.SetTagName("header")
	headerDecoder.SetMode(form.ModeExplicit)
	formEncoder.SetTagName("form")
	formEncoder.SetMode(form.ModeExplicit)
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
	headerEncoder.SetTagName("query")
	headerEncoder.SetMode(form.ModeExplicit)
}

func Decode(r *http.Request, v any, flags ...Flag) error {
	if err := DecodeForm(r, v, flags...); err != nil {
		return err
	}
	if err := DecodeQuery(r, v, flags...); err != nil {
		return err
	}
	if err := DecodeHeader(r, v, flags...); err != nil {
		return err
	}
	return nil
}

func DecodeFormValues(values map[string][]string, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vacuum(values)
	}
	return formDecoder.Decode(v, values)
}

func DecodeForm(r *http.Request, v any, flags ...Flag) error {
	r.ParseForm()
	return DecodeFormValues(r.Form, v, flags...)
}

func DecodeQueryValues(values map[string][]string, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vacuum(values)
	}
	return queryDecoder.Decode(v, values)
}

func DecodeQuery(r *http.Request, v any, flags ...Flag) error {
	return DecodeQueryValues(r.URL.Query(), v, flags...)
}

func DecodeHeaderValues(values map[string][]string, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vacuum(values)
	}
	return headerDecoder.Decode(v, values)
}

func DecodeHeader(r *http.Request, v any, flags ...Flag) error {
	return DecodeHeaderValues(r.Header, v, flags...)
}

func EncodeForm(v any) (url.Values, error) {
	return formEncoder.Encode(v)
}

func EncodeQuery(v any) (url.Values, error) {
	return queryEncoder.Encode(v)
}

func EncodeHeader(v any) (url.Values, error) {
	return headerEncoder.Encode(v)
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
