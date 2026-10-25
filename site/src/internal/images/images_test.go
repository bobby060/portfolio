package images

import (
	"fmt"
	"testing"

	"portfolio/internal/gallery"
)

// fakeBuilder is a test double: predictable URLs, so Srcset is tested on its own.
type fakeBuilder struct{}

func (fakeBuilder) URL(img gallery.Image, w int) string { return fmt.Sprintf("/%s/%d", img.SHA, w) }

func TestSrcset(t *testing.T) {
	tests := []struct {
		name string
		img  gallery.Image
		want string
	}{
		{"large original gets every width", gallery.Image{SHA: "x", Width: 6000, Height: 4000},
			"/x/400 400w, /x/800 800w, /x/1280 1280w, /x/1920 1920w"},
		{"exact match is included", gallery.Image{SHA: "x", Width: 1920, Height: 1280},
			"/x/400 400w, /x/800 800w, /x/1280 1280w, /x/1920 1920w"},
		{"skips widths larger than original", gallery.Image{SHA: "x", Width: 1000, Height: 500},
			"/x/400 400w, /x/800 800w"},
		{"tiny original still gets smallest", gallery.Image{SHA: "x", Width: 200, Height: 100},
			"/x/400 400w"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Srcset(fakeBuilder{}, tc.img); got != tc.want {
				t.Errorf("Srcset() =\n  %q\nwant\n  %q", got, tc.want)
			}
		})
	}
}

func TestLargest(t *testing.T) {
	tests := []struct {
		width, want int
	}{
		{6000, 1920},
		{1920, 1920},
		{1000, 800},
		{200, 400}, // smaller than every variant: the smallest, matching Srcset
	}
	for _, tc := range tests {
		if got := Largest(gallery.Image{Width: tc.width, Height: 100}); got != tc.want {
			t.Errorf("Largest(width %d) = %d, want %d", tc.width, got, tc.want)
		}
	}
}
