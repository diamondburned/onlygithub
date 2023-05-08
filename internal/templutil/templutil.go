package templutil

import (
	"context"
	"html/template"
	"io"

	"github.com/a-h/templ"
)

// UnsafeHTML is a hack to get around the templ package's lack of support for
// inserting unsafe HTML.
func UnsafeHTML(html template.HTML) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		io.WriteString(w, string(html))
		return nil
	})
}
