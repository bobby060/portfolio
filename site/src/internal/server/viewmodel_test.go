package server

import (
	"strings"
	"testing"
	"time"

	"portfolio/internal/gallery"
	"portfolio/internal/images"
)

func TestFormatDate(t *testing.T) {
	if got := formatDate(time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC)); got != "7/4/2023" {
		t.Errorf("formatDate = %q, want 7/4/2023", got)
	}
}

func TestAlbumCard(t *testing.T) {
	s := &Server{imgs: images.Picsum{}}
	a := gallery.Album{
		Title: "Killer Whales",
		Slug:  "killer-whales",
		Date:  time.Date(2023, 7, 14, 0, 0, 0, 0, time.UTC),
		Cover: "c2",
		Images: []gallery.Image{
			{SHA: "c1", Width: 600, Height: 400},
			{SHA: "c2", Width: 400, Height: 600},
		},
	}
	c := s.albumCard(a, "100vw")

	if c.Href != "/albums/killer-whales" {
		t.Errorf("Href = %q", c.Href)
	}
	if c.Date != "7/14/2023" {
		t.Errorf("Date = %q", c.Date)
	}
	if c.Img.W != 400 || c.Img.H != 600 || !strings.Contains(c.Img.Src, "/seed/c2/") {
		t.Errorf("card should use the cover image c2, got %+v", c.Img)
	}
	if c.Img.Alt != "Killer Whales" {
		t.Errorf("Alt = %q", c.Img.Alt)
	}
}
