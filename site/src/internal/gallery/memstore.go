package gallery

import (
	"slices"
)

// MemStore is an in-memory Store built from already-validated albums.
type MemStore struct {
	site   Site
	albums []Album        // newest first
	bySlug map[string]int // slug -> index into albums
	tags   []string       // unique, sorted
}

// Compile-time check that MemStore implements Store.
var _ Store = (*MemStore)(nil)

func newMemStore(site Site, albums []Album) *MemStore {
	sorted := slices.Clone(albums)
	slices.SortStableFunc(sorted, func(a, b Album) int { return b.Date.Compare(a.Date) })

	s := &MemStore{site: site, albums: sorted, bySlug: make(map[string]int, len(sorted))}
	seen := map[string]bool{}
	for i, a := range sorted {
		s.bySlug[a.Slug] = i
		for _, t := range a.Tags {
			if !seen[t] {
				seen[t] = true
				s.tags = append(s.tags, t)
			}
		}
	}
	slices.Sort(s.tags)
	return s
}

func (s *MemStore) Site() Site { return s.site }

func (s *MemStore) Album(slug string) (Album, bool) {
	i, ok := s.bySlug[slug]
	if !ok {
		return Album{}, false
	}
	return s.albums[i], true
}

func (s *MemStore) Featured() []Album {
	var out []Album
	for _, a := range s.albums {
		if a.Featured {
			out = append(out, a)
		}
	}
	return out
}

func (s *MemStore) Tags() []string { return slices.Clone(s.tags) }

// Albums returns one page of albums, newest first.
func (s *MemStore) Albums(q AlbumQuery) AlbumPage {
	var albums []Album

	if q.PerPage < 1 {
		panic("gallery: AlbumQuery.PerPage must be > 0")
	}

	if q.Tag == "" {
		albums = s.albums
	} else {
		for _, album := range s.albums {
			if slices.Contains(album.Tags, q.Tag) {
				albums = append(albums, album)
			}
		}
	}

	page := max(q.Page, 1)

	start := (page - 1) * q.PerPage

	if start >= len(albums) {
		return AlbumPage{Next: 0}
	}

	end := min(start+q.PerPage, len(albums))

	nextPage := 0
	if end < len(albums) {
		nextPage = page + 1
	}

	return AlbumPage{Albums: slices.Clone(albums[start:end]), Next: nextPage}
}
