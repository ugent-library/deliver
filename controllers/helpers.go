package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

var formDecoder = schema.NewDecoder()

func init() {
	formDecoder.IgnoreUnknownKeys(true)
	formDecoder.SetAliasTag("form")
}

func bindForm(r *http.Request, b any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return formDecoder.Decode(b, r.Form)
}

func yield[C any, V any](c C, v V) Yield[C, V] {
	return Yield[C, V]{c, v}
}

type Yield[C any, V any] struct {
	Ctx C
	Var V
}
