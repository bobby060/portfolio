// Package gallery holds the portfolio's domain model and the store that serves it.
// It knows nothing about HTTP or HTML.
package gallery

import "time"

// Image is one photo. SHA is the content hash that identifies its bytes in blob storage.
type Image struct {
	SHA     string    `yaml:"sha"`
	Width   int       `yaml:"width"`
	Height  int       `yaml:"height"`
	TakenAt time.Time `yaml:"taken_at"`
	Caption string    `yaml:"caption"`
}

// Album is one photo album, as described by albums/<dir>/album.yaml.
type Album struct {
	Title       string    `yaml:"title"`
	Slug        string    `yaml:"slug"` // defaults to the directory name
	Date        time.Time `yaml:"date"`
	Description string    `yaml:"description"`
	Tags        []string  `yaml:"tags"`
	Featured    bool      `yaml:"featured"`
	Cover       string    `yaml:"cover"`  // image SHA; defaults to the first image
	Images      []Image   `yaml:"images"` // display order
}

// CoverImage returns the album's cover image.
// Albums returned by a Store always have a valid cover.
func (a Album) CoverImage() Image {
	for _, img := range a.Images {
		if img.SHA == a.Cover {
			return img
		}
	}
	return a.Images[0]
}

// Link is a labelled hyperlink.
type Link struct {
	Label string `yaml:"label"`
	Href  string `yaml:"href"`
}

// Site is the site-wide content from site.yaml.
type Site struct {
	Title   string  `yaml:"title"`
	Tagline string  `yaml:"tagline"`
	Links   []Link  `yaml:"links"`
	Hero    []Image `yaml:"hero"`
}

// Store is read-only access to the portfolio's content.
// Phase 3 replaces the YAML-backed implementation behind this same interface.
type Store interface {
	Site() Site
	Albums(q AlbumQuery) AlbumPage   // newest first, optionally filtered by tag
	Album(slug string) (Album, bool) // false if slug is unknown
	Featured() []Album               // featured albums, newest first
	Tags() []string                  // every tag in use, unique and sorted
}

// AlbumQuery selects one page of albums.
type AlbumQuery struct {
	Tag     string // "" means all albums
	Page    int    // 1-based; values < 1 are treated as 1
	PerPage int    // must be > 0
}

// AlbumPage is one page of query results.
type AlbumPage struct {
	Albums []Album
	Next   int // page number to request next; 0 means this is the last page
}
