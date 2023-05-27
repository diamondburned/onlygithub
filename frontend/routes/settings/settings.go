package settings

import (
	"context"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"

	_ "embed"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

type Services struct {
	Config onlygithub.ConfigService
	Images onlygithub.ImageService
}

type handler struct{ Services }

func New(s Services) http.Handler {
	h := &handler{s}

	r := chi.NewRouter()
	r.Use(frontend.LoggedInOnly)

	r.Get("/", h.get)
	r.Group(func(r chi.Router) {
		r.Use(frontend.ParseMultipartForm)
		r.Post("/site", h.saveSite)
	})

	return r
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)
	session := frontend.SessionFromRequest(r)

	cfg, err := h.Config.UserConfig(ctx, session.Me.ID)
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	data := settingsData{
		Me:       session.Me,
		MeConfig: cfg,
	}

	if session.Me.IsOwner {
		siteCfg, err := h.Config.SiteConfig(ctx)
		if err != nil {
			layouts.RenderError(w, r, err)
			return
		}

		data.SiteConfig = siteCfg
	}

	settings(r, site, owner, data).Render(ctx, w)
}

func (h handler) saveSite(w http.ResponseWriter, r *http.Request) {
	session := frontend.SessionFromRequest(r)
	if !session.Me.IsOwner {
		layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
		return
	}

	var form struct {
		onlygithub.Socials
		Description        string                `form:"description"`
		About              string                `form:"about"`
		CustomCSS          string                `form:"custom-css"`
		AllowComments      bool                  `form:"allow-comments"`
		AllowDMs           bool                  `form:"allow-dms"`
		AllowReactions     bool                  `form:"allow-reactions"`
		HomepageVisibility onlygithub.Visibility `form:"homepage-visibility"`
		Banner             *multipart.FileHeader `form:"-"`
		Avatar             *multipart.FileHeader `form:"-"`
	}

	if err := hrt.URLDecoder.Decode(r, &form); err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	form.Banner = getMultipartFile(r, "banner")
	form.Avatar = getMultipartFile(r, "avatar")

	site := frontend.SiteFromRequest(r)

	// These ones we can apply directly, since we always put them in the page.
	site.Socials = form.Socials
	site.Description = template.HTML(form.Description)
	site.About = template.HTML(form.About)
	site.CustomCSS = form.CustomCSS
	site.AllowComments = form.AllowComments
	site.AllowDMs = form.AllowDMs
	site.AllowReactions = form.AllowReactions
	site.HomepageVisibility = form.HomepageVisibility

	log.Printf("Homepage visibility: %v", form.HomepageVisibility)

	if form.Avatar != nil {
		if err := h.replaceAsset(r.Context(), &site.AvatarAsset, site.HomepageVisibility, form.Avatar); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to replace avatar"))
			return
		}
	} else if site.AvatarAsset != nil {
		if err := h.Images.SetImageVisibility(r.Context(), *site.AvatarAsset, site.HomepageVisibility); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to set avatar visibility"))
			return
		}
	}

	if form.Banner != nil {
		if err := h.replaceAsset(r.Context(), &site.BannerAsset, site.HomepageVisibility, form.Banner); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to replace banner"))
			return
		}
	} else if site.BannerAsset != nil {
		if err := h.Images.SetImageVisibility(r.Context(), *site.BannerAsset, site.HomepageVisibility); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to set banner visibility"))
			return
		}
	}

	if err := h.Config.SetSiteConfig(r.Context(), site); err != nil {
		layouts.RenderError(w, r, errors.Wrap(err, "failed to save site config"))
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (h handler) replaceAsset(ctx context.Context, oldAssetID **onlygithub.ID, visibility onlygithub.Visibility, header *multipart.FileHeader) error {
	// Delete the old banner if it exists.
	if *oldAssetID != nil {
		if err := h.Images.DeleteImage(ctx, **oldAssetID); err != nil {
			log := hclog.FromContext(ctx)
			log.Warn("failed to delete old asset", "err", err)
		}
	}

	f, err := header.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open asset form file")
	}
	defer f.Close()

	req := onlygithub.UploadImageRequest{
		Filename:    header.Filename,
		Visibility:  visibility, // same as
		MinimumCost: 0,          // public so always 0
	}

	a, err := h.Images.UploadImage(ctx, req, f)
	if err != nil {
		return errors.Wrap(err, "failed to upload asset")
	}

	*oldAssetID = &a.ID
	return nil
}

func getMultipartFile(r *http.Request, name string) *multipart.FileHeader {
	if values := r.MultipartForm.File[name]; len(values) > 0 {
		return values[0]
	}
	return nil
}
