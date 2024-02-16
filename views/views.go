//go:generate go get -u github.com/a-h/templ/cmd/templ
//go:generate templ generate
package views

import (
	"context"
	"strings"

	"github.com/a-h/templ"
)

func String(c templ.Component) string {
	b := strings.Builder{}
	// TODO handle context and error
	c.Render(context.TODO(), &b)
	return b.String()
}

type SelectOption struct {
	Value string
	Label string
}
