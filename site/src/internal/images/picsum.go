package images

import (
	"fmt"
	"net/url"

	"portfolio/internal/gallery"
)

// Picsum serves placeholder photos from picsum.photos, seeded by each image's SHA
// so the same image always gets the same picture. For mock content and dev only.
type Picsum struct{}

var _ URLBuilder = Picsum{}

func (Picsum) URL(img gallery.Image, width int) string {
	// Integer rounding to nearest: adding half the divisor before dividing
	// turns Go's round-down division into round-half-up.
	height := (width*img.Height + img.Width/2) / img.Width
	height = max(height, 1) // picsum rejects a zero dimension
	return fmt.Sprintf("https://picsum.photos/seed/%s/%d/%d", url.PathEscape(img.SHA), width, height)
}
