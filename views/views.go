package views

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func raw(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, s)
		return
	})
}
