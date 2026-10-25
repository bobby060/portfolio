package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"portfolio/web/components"
)

const (
	heroMax       = 5 // must match the five .hero-slide delays in web/static/css/input.css
	albumsPerPage = 6 // divisible by 2 and 3, so grid rows are full at every breakpoint

	// sizes tell the browser how wide each image displays, so it downloads the right variant
	sizesHero     = "100vw"
	sizesFeatured = "(min-width: 1280px) 67vw, (min-width: 1024px) 83vw, 100vw"
	sizesGrid     = "(min-width: 751px) 33vw, (min-width: 351px) 50vw, 100vw"
)

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

// home is "/": the unfiltered first page.
func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, http.StatusOK, components.Home(s.homePage("", 1)))
}

// albums is "/albums?tag=&page=": the home page filtered by tag, showing pages 1..page.
// It always returns a full page, so it works without JS and is safe to cache by URL.
func (s *Server) albums(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	s.render(w, r, http.StatusOK, components.Home(s.homePage(q.Get("tag"), parsePage(q.Get("page")))))
}

// albumsFragment is "/fragments/albums?tag=&page=": one batch of cards for infinite scroll.
// It has its own URL (not an HX-Request check on /albums) because the edge cache keys on URL alone.
func (s *Server) albumsFragment(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	grid := s.albumGrid(q.Get("tag"), parsePage(q.Get("page")), false)
	w.Header().Set("X-Robots-Tag", "noindex") // partial HTML; keep it out of search results
	s.render(w, r, http.StatusOK, components.AlbumGridItems(grid))
}

func (s *Server) homePage(tag string, page int) components.HomePage {
	site := s.store.Site()
	p := components.HomePage{Title: site.Title, Tagline: site.Tagline}

	for _, l := range site.Links {
		p.Links = append(p.Links, components.Link{
			Label: l.Label, Href: l.Href, External: strings.HasPrefix(l.Href, "http"),
		})
	}
	for _, h := range site.Hero[:min(len(site.Hero), heroMax)] {
		p.Hero = append(p.Hero, s.img(h, "", sizesHero)) // decorative: empty alt
	}
	for _, a := range s.store.Featured() {
		p.Featured = append(p.Featured, s.albumCard(a, sizesFeatured))
	}
	p.Tags = tagChips(s.store.Tags(), tag)
	p.Grid = s.albumGrid(tag, page, true)
	return p
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
