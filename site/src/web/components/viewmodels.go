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

// HomePage is everything the home page shows.
type HomePage struct {
	Title    string
	Tagline  string
	Links    []Link
	Hero     []Img
	Featured []AlbumCard
	Albums   []AlbumCard
}
