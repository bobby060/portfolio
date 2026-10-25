package server

import (
	"net/url"
	"time"

	"portfolio/internal/gallery"
	"portfolio/internal/images"
	"portfolio/web/components"
)

// formatDate matches the old site's M/D/YYYY format, with no zero padding.
func formatDate(t time.Time) string { return t.Format("1/2/2006") }

func (s *Server) img(img gallery.Image, alt, sizes string) components.Img {
	return components.Img{
		Src:    s.imgs.URL(img, 800),
		Srcset: images.Srcset(s.imgs, img),
		Sizes:  sizes,
		W:      img.Width,
		H:      img.Height,
		Alt:    alt,
	}
}

func (s *Server) albumCard(a gallery.Album, sizes string) components.AlbumCard {
	return components.AlbumCard{
		Title: a.Title,
		Href:  "/albums/" + url.PathEscape(a.Slug),
		Desc:  a.Description,
		Date:  formatDate(a.Date),
		Img:   s.img(a.CoverImage(), a.Title, sizes),
	}
}
