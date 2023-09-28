//go:generate go get -u github.com/a-h/templ/cmd/templ
//go:generate templ generate
package views

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// TODO eliminate need for this
func rawHTML(text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		io.WriteString(w, text)
		return nil
	})
}

func String(c templ.Component) string {
	b := strings.Builder{}
	// TODO handle context and error
	c.Render(context.TODO(), &b)
	return b.String()
}
