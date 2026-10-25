package server

import (
	"fmt"
	"net/http"

	"portfolio/internal/gallery"
	"portfolio/internal/images"
	"portfolio/web/components"
)

// Photo column width: 3 columns inside max-w-7xl (1280px), else a share of the viewport.
const sizesPhoto = "(min-width: 1280px) 427px, (min-width: 751px) 33vw, (min-width: 351px) 50vw, 100vw"

func (s *Server) album(w http.ResponseWriter, r *http.Request) {
	a, ok := s.store.Album(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	s.render(w, r, http.StatusOK, components.Album(s.albumPage(a)))
}

func (s *Server) albumPage(a gallery.Album) components.AlbumPage {
	cover := a.CoverImage()
	p := components.AlbumPage{
		Title:     a.Title,
		Date:      formatDate(a.Date),
		Desc:      a.Description,
		HeaderURL: s.imgs.URL(cover, images.Largest(cover)),
		HeaderW:   cover.Width,
		HeaderH:   cover.Height,
	}
	for i, img := range a.Images {
		alt := img.Caption
		if alt == "" {
			alt = fmt.Sprintf("%s, photo %d", a.Title, i+1)
		}
		full := images.Largest(img)
		p.Photos = append(p.Photos, components.Photo{
			Img:   s.img(img, alt, sizesPhoto),
			Full:  s.imgs.URL(img, full),
			FullW: full,
			FullH: images.Height(img, full),
		})
	}
	return p
}
