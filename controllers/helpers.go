package controllers

import (
	"io"
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

func detectContentType(f io.ReadSeeker) (string, error) {
	b := make([]byte, 512)
	if _, err := f.Read(b); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(b)

	// rewind
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return contentType, nil
}
