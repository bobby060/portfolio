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
	return fmt.Sprintf("https://picsum.photos/seed/%s/%d/%d", url.PathEscape(img.SHA), width, Height(img, width))
}
