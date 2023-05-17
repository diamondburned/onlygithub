// Code generated by templ@(devel) DO NOT EDIT.

package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import "libdb.so/onlygithub"
import "libdb.so/onlygithub/frontend/components"

type MainOpts struct {
	Me *onlygithub.User // current user, optional
}

func Main(id, title string, site *onlygithub.SiteConfig, owner *onlygithub.User) templ.Component {
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
		// TemplElement
		err = components.Head(components.HeadOpts{Title: title, Owner: owner}).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<body")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"main\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" id=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(id))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Children
		err = var_1.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</body>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

