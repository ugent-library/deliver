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
