// Code generated by templ@(devel) DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import "libdb.so/onlygithub"

type HeadOpts struct {
	Title string
	Owner *onlygithub.User
}

func headTitle(title string, owner *onlygithub.User) string {
	if title == "" {
		return owner.Username
	}
	return title + " – " + owner.Username
}

func Head(opts HeadOpts) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (void)
		_, err = templBuffer.WriteString("<meta")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" name=\"viewport\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" content=\"width=device-width, initial-scale=1\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://necolas.github.io/normalize.css/8.0.1/normalize.css\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://cdn.jsdelivr.net/npm/tiny-markdown-editor@0.1.5/dist/tiny-mde.min.css\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://fonts.googleapis.com/icon?family=Material+Icons\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://fonts.googleapis.com/css2?family=Source+Sans+Pro:ital,wght@0,400;0,600;0,700;1,400;1,600;1,700&amp;display=swap\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;500;600;700&amp;display=swap\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://fonts.googleapis.com/css2?family=Nunito:wght@400;500;600;700;800&amp;display=swap\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://fonts.googleapis.com/css2?family=Lato:wght@400;700;900&amp;display=swap\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"/static/styles.css\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if opts.Owner != nil {
			// Element (standard)
			_, err = templBuffer.WriteString("<title>")
			if err != nil {
				return err
			}
			// StringExpression
			var var_2 string = headTitle(opts.Title, opts.Owner)
			_, err = templBuffer.WriteString(templ.EscapeString(var_2))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</title>")
			if err != nil {
				return err
			}
		} else if opts.Title != "" {
			// Element (standard)
			_, err = templBuffer.WriteString("<title>")
			if err != nil {
				return err
			}
			// StringExpression
			var var_3 string = opts.Title
			_, err = templBuffer.WriteString(templ.EscapeString(var_3))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</title>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

