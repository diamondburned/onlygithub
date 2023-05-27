// Code generated by templ@(devel) DO NOT EDIT.

package index

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import "libdb.so/onlygithub/internal/templutil"
import "libdb.so/onlygithub/frontend/components"
import "libdb.so/onlygithub/frontend/layouts"
import "libdb.so/onlygithub"
import "net/http"
import "time"

type indexOpts struct {
	Me *onlygithub.User // optional
	Posts []onlygithub.Post
	OwnerAvatarURL string
}

func index(r *http.Request, site *onlygithub.SiteConfig, owner *onlygithub.User, opts indexOpts) templ.Component {
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
		// DocType
		_, err = templBuffer.WriteString(`<!doctype html>`)
		if err != nil {
			return err
		}
		// TemplElement
		var_2 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<header>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"user-banner\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// If
			if site.BannerAsset != nil {
				// Element (void)
				_, err = templBuffer.WriteString("<img")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" src=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString("/images/" + site.BannerAsset.String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" alt=\"banner\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"user-info\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Element (void)
			_, err = templBuffer.WriteString("<img")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"user-avatar\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" src=")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(opts.OwnerAvatarURL))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" alt=\"avatar\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<h1")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"user-name\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// If
			if owner.RealName != "" {
				// Element (standard)
				_, err = templBuffer.WriteString("<span")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"real-name\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_3 string = owner.RealName
				_, err = templBuffer.WriteString(templ.EscapeString(var_3))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				// Whitespace (normalised)
				_, err = templBuffer.WriteString(` `)
				if err != nil {
					return err
				}
				// Element (void)
				_, err = templBuffer.WriteString("<br>")
				if err != nil {
					return err
				}
				// Whitespace (normalised)
				_, err = templBuffer.WriteString(` `)
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<small")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"username\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_4 string = owner.Username
				_, err = templBuffer.WriteString(templ.EscapeString(var_4))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</small>")
				if err != nil {
					return err
				}
			} else {
				// Element (standard)
				_, err = templBuffer.WriteString("<span")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"username\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_5 string = owner.Username
				_, err = templBuffer.WriteString(templ.EscapeString(var_5))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</h1>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<p")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"site-description\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// TemplElement
			err = templutil.UnsafeHTML(site.Description).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</header>")
			if err != nil {
				return err
			}
			// Whitespace (normalised)
			_, err = templBuffer.WriteString(` `)
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<nav>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"left\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<ul>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<li>")
			if err != nil {
				return err
			}
			// TemplElement
			err = navButton(r, "/", "Home").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<li>")
			if err != nil {
				return err
			}
			// TemplElement
			err = navButton(r, "/about", "About").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<li>")
			if err != nil {
				return err
			}
			// TemplElement
			err = navButton(r, "/membership", "Membership").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</li>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</ul>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"right\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"padding\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"actions\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// If
			if opts.Me != nil && opts.Me.ID == owner.ID {
				// Element (standard)
				_, err = templBuffer.WriteString("<a")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" href=\"/create\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" role=\"button\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// TemplElement
				err = components.Icon("create", components.InlineIcon).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a>")
				if err != nil {
					return err
				}
			}
			// If
			if opts.Me != nil {
				// Element (standard)
				_, err = templBuffer.WriteString("<a")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" href=\"/settings\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" role=\"button\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" title=\"Settings\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// TemplElement
				err = components.Icon("settings", components.InlineIcon).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a>")
				if err != nil {
					return err
				}
			} else {
				// Element (standard)
				_, err = templBuffer.WriteString("<a")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" href=\"/login\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" role=\"button\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// TemplElement
				err = components.Icon("login", components.InlineIcon).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"padding\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</nav>")
			if err != nil {
				return err
			}
			// Whitespace (normalised)
			_, err = templBuffer.WriteString(` `)
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<section")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" id=\"posts\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<h2>")
			if err != nil {
				return err
			}
			// Text
			var_6 := `Latest Posts`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2>")
			if err != nil {
				return err
			}
			// For
			for _, post := range opts.Posts {
				// Element (standard)
				// Element CSS
				var var_7 = []any{
						"post",
						templ.KV("concealed", post.IsConcealed),
						templ.KV("has-images", len(post.Images) > 0),
					}
				err = templ.RenderCSSItems(ctx, templBuffer, var_7...)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<article")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" id=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString("post-" + post.ID.String()))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" class=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_7).String()))
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
				// If
				if post.IsConcealed {
					// Element (standard)
					_, err = templBuffer.WriteString("<div")
					if err != nil {
						return err
					}
					// Element Attributes
					_, err = templBuffer.WriteString(" class=\"concealed-overlay\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(">")
					if err != nil {
						return err
					}
					// TemplElement
					err = components.Icon("paid", components.LargeIcon).Render(ctx, templBuffer)
					if err != nil {
						return err
					}
					// Element (standard)
					_, err = templBuffer.WriteString("<span>")
					if err != nil {
						return err
					}
					// StringExpression
					var var_8 string = "Post is only available to " + visibilityMap[post.Visibility]
					_, err = templBuffer.WriteString(templ.EscapeString(var_8))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</span>")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</div>")
					if err != nil {
						return err
					}
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"cover\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// For
				for _, image := range post.Images {
					// Element (void)
					_, err = templBuffer.WriteString("<img")
					if err != nil {
						return err
					}
					// Element Attributes
					if !image.ID.IsZero() {
						// Element Attributes
						_, err = templBuffer.WriteString(" src=")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(templ.EscapeString("/images/" + image.ID.String()))
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(" alt=\"post image\"")
						if err != nil {
							return err
						}
					} else {
						// Element Attributes
						_, err = templBuffer.WriteString(" src=")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(templ.EscapeString(image.PreviewURL))
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString("\"")
						if err != nil {
							return err
						}
						_, err = templBuffer.WriteString(" alt=\"blurred post image\"")
						if err != nil {
							return err
						}
					}
					_, err = templBuffer.WriteString(">")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"body\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"author\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// Element (void)
				_, err = templBuffer.WriteString("<img")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"avatar\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" src=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(opts.OwnerAvatarURL))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" alt=\"avatar\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<div>")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<span")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"username\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_9 string = owner.Username
				_, err = templBuffer.WriteString(templ.EscapeString(var_9))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<time")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"relative\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" datetime=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(post.ID.Time().Format(time.RFC3339)))
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
				// StringExpression
				var var_10 string = post.ID.Time().Format("January _2, 2006 at 03:04pm")
				_, err = templBuffer.WriteString(templ.EscapeString(var_10))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</time>")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"content markdown markdown-unsafe\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_11 string = post.Markdown
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</article>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</section>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = layouts.Main("index", "", site, owner).Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func navButton(r *http.Request, dst templ.SafeURL, name string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_12 := templ.GetChildren(ctx)
		if var_12 == nil {
			var_12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_13 templ.SafeURL = dst
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_13)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		if r.URL.Path == string(dst) {
			_, err = templBuffer.WriteString(" data-active")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_14 string = name
		_, err = templBuffer.WriteString(templ.EscapeString(var_14))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

