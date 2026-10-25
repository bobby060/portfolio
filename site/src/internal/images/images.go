// Package images turns gallery images into URLs. It is the only package that
// knows where image bytes are served from.
package images

import (
	"fmt"
	"strings"

	"portfolio/internal/gallery"
)

// URLBuilder returns the URL of img resized to the given width in pixels,
// keeping img's aspect ratio. width must be > 0.
type URLBuilder interface {
	URL(img gallery.Image, width int) string
}

// Widths are the resized variants available for every image, smallest first.
// They match the variants the ingest pipeline will generate.
var Widths = []int{400, 800, 1280, 1920}

// Srcset returns an HTML srcset value listing img's variants, for example
// "https://…/400/267 400w, https://…/800/533 800w".
// Widths larger than the original are skipped, but the smallest is always included.
func Srcset(b URLBuilder, img gallery.Image) string {
	var parts []string
	for _, w := range Widths {
		if w > img.Width && len(parts) > 0 {
			break
		}
		parts = append(parts, fmt.Sprintf("%s %dw", b.URL(img, w), w))
	}
	return strings.Join(parts, ", ")
}

// Height returns img's height when scaled to width, rounded to the nearest pixel
// and at least 1.
func Height(img gallery.Image, width int) int {
	// Adding half the divisor before dividing turns Go's round-down division into round-to-nearest
	return max((width*img.Height+img.Width/2)/img.Width, 1)
}

// Largest returns the widest variant Srcset offers for img.
func Largest(img gallery.Image) int {
	best := Widths[0]
	for _, w := range Widths {
		if w <= img.Width {
			best = w
		}
	}
	return best
}
