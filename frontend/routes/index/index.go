package index

import (
	"net/http"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

var visibilityMap = map[onlygithub.Visibility]string{
	onlygithub.NotVisible:        "the owner",
	onlygithub.VisibleToSponsors: "sponsors",
	onlygithub.VisibleToPrivate:  "signed in users",
	onlygithub.VisibleToPublic:   "everyone",
}

type Services struct {
	Posts onlygithub.PostService
}

type handler struct{ Services }

func New(s Services) http.Handler {
	h := &handler{s}

	r := chi.NewRouter()
	r.Get("/", h.get)

	return r
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)

	var form struct {
		Before onlygithub.ID `form:"before"`
	}

	if err := hrt.URLDecoder.Decode(r, &form); err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	var opts indexOpts

	if site.AvatarAsset != nil {
		opts.OwnerAvatarURL = "/images/" + site.AvatarAsset.String()
	} else {
		opts.OwnerAvatarURL = owner.AvatarURL
	}

	session := frontend.SessionFromRequest(r)
	opts.Me = session.Me

	switch site.HomepageVisibility {
	case onlygithub.NotVisible:
		if opts.Me == nil || opts.Me.ID != owner.ID {
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToSponsors:
		if opts.Me == nil || opts.Me.Sponsorship == nil {
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToPrivate:
		if opts.Me == nil {
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToPublic:
		// Nothing.
	}

	posts, err := h.Posts.Posts(r.Context(), opts.Me, form.Before)
	if err != nil {
		layouts.RenderError(w, r, errors.Wrap(err, "failed to get posts"))
		return
	}

	opts.Posts = posts

	index(r, site, owner, opts).Render(r.Context(), w)
}

func kickToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
