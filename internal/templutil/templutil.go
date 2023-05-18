package templutil

import (
	"context"
	"html/template"
	"io"
	"strings"

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

// Literal prints the given (possibly indented) text as regular dedented text.
func Literal(text string) string {
	lines := strings.Split(text, "\n")
	var indent int
	for _, line := range lines {
		if i := countIndent(line); i > 0 && (indent == 0 || i < indent) {
			indent = i
		}
	}

	for i := range lines {
		for j := 0; j < indent; j++ {
			lines[i] = strings.TrimPrefix(lines[i], "\t")
		}
	}

	text = strings.Join(lines, "\n")
	return text
}

func countIndent(line string) int {
	for i, r := range line {
		if r != '\t' {
			return i
		}
	}
	return 0
}
