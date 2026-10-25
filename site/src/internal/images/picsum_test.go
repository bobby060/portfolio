package images

import (
	"testing"

	"portfolio/internal/gallery"
)

func TestPicsumURL(t *testing.T) {
	tests := []struct {
		name  string
		img   gallery.Image
		width int
		want  string
	}{
		{"landscape rounds height to nearest", gallery.Image{SHA: "land", Width: 6000, Height: 4000}, 400,
			"https://picsum.photos/seed/land/400/267"},
		{"portrait", gallery.Image{SHA: "port", Width: 4000, Height: 6000}, 400,
			"https://picsum.photos/seed/port/400/600"},
		{"square", gallery.Image{SHA: "sq", Width: 5000, Height: 5000}, 800,
			"https://picsum.photos/seed/sq/800/800"},
		{"hero", gallery.Image{SHA: "hero", Width: 1920, Height: 1280}, 1280,
			"https://picsum.photos/seed/hero/1280/853"},
		{"extreme panorama keeps height >= 1", gallery.Image{SHA: "pano", Width: 10000, Height: 10}, 400,
			"https://picsum.photos/seed/pano/400/1"},
		{"seed is path-escaped", gallery.Image{SHA: "a b/c", Width: 100, Height: 100}, 100,
			"https://picsum.photos/seed/a%20b%2Fc/100/100"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := (Picsum{}).URL(tc.img, tc.width); got != tc.want {
				t.Errorf("URL() = %q, want %q", got, tc.want)
			}
		})
	}
}
