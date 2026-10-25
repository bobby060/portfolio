package components

// View models: plain data the templates render. Handlers in internal/server build
// them, so templates never format dates, build URLs, or see gallery types.

// Img is everything an <img> needs.
type Img struct {
	Src    string // fallback URL for browsers that ignore srcset
	Srcset string
	Sizes  string // how wide the image displays, so the browser picks a variant
	W, H   int    // intrinsic size: reserves space before load (no layout shift)
	Alt    string
}

// AlbumCard is one album tile, used in the grid and the featured carousel.
type AlbumCard struct {
	Title string
	Href  string
	Desc  string
	Date  string // already formatted, e.g. "7/14/2023"
	Img   Img
}

// Link is a button or link from site.yaml.
type Link struct {
	Label    string
	Href     string
	External bool // opens in a new tab
}

// HomePage is everything the home page (and the filtered /albums page) shows.
type HomePage struct {
	Title    string
	Tagline  string
	Links    []Link
	Hero     []Img
	Featured []AlbumCard
	Tags     []TagChip
	Grid     AlbumGrid
}

// TagChip is one tag filter button.
type TagChip struct {
	Name   string
	Href   string // full-page URL: works without JS, and is what the address bar shows
	HX     string // URL htmx fetches; the #albums section of the response replaces ours
	Active bool
}

// AlbumGrid is a batch of album cards and where to get the next batch.
type AlbumGrid struct {
	Albums   []AlbumCard
	MoreHref string // full-page link to the next page; "" when there are no more
	MoreHX   string // fragment URL htmx fetches for infinite scroll
}

// AlbumPage is everything an album page shows.
type AlbumPage struct {
	Title     string
	Date      string
	Desc      string
	HeaderURL string // cover image, used as a CSS background for the parallax header
	HeaderW   int    // cover aspect ratio, which sets the header height
	HeaderH   int
	Photos    []Photo
}

// Photo is one tile in the album's masonry grid.
type Photo struct {
	Img          Img    // grid thumbnail
	Full         string // largest variant: the lightbox image, and the no-JS link target
	FullW, FullH int    // Full's pixel size; PhotoSwipe needs it before the image loads
}
