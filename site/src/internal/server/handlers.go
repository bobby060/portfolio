package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"portfolio/internal/gallery"
	"portfolio/web/components"
)

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

const (
	heroMax       = 5 // must match the five .hero-slide delays in web/static/css/input.css
	albumsPerPage = 8

	// sizes tell the browser how wide each image displays, so it downloads the right variant
	sizesHero     = "100vw"
	sizesFeatured = "(min-width: 1280px) 67vw, (min-width: 1024px) 83vw, 100vw"
	sizesGrid     = "(min-width: 751px) 33vw, (min-width: 351px) 50vw, 100vw"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	site := s.store.Site()
	page := components.HomePage{Title: site.Title, Tagline: site.Tagline}

	for _, l := range site.Links {
		page.Links = append(page.Links, components.Link{
			Label: l.Label, Href: l.Href, External: strings.HasPrefix(l.Href, "http"),
		})
	}
	for _, h := range site.Hero[:min(len(site.Hero), heroMax)] {
		page.Hero = append(page.Hero, s.img(h, "", sizesHero)) // decorative: empty alt
	}
	for _, a := range s.store.Featured() {
		page.Featured = append(page.Featured, s.albumCard(a, sizesFeatured))
	}
	for _, a := range s.store.Albums(gallery.AlbumQuery{Page: 1, PerPage: albumsPerPage}).Albums {
		page.Albums = append(page.Albums, s.albumCard(a, sizesGrid))
	}

	s.render(w, r, http.StatusOK, components.Home(page))
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, http.StatusNotFound, components.NotFound())
}

// render writes c with the given status. It renders into a buffer first, so a
// template error becomes a clean 500 instead of a half-written page sent as 200.
func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	var buf bytes.Buffer
	if err := c.Render(r.Context(), &buf); err != nil {
		s.log.Error("render failed", "path", r.URL.Path, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}
