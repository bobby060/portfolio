package gallery

import (
	"slices"
	"testing"
	"time"
)

func date(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t
}

func slugs(albums []Album) []string {
	out := []string{}
	for _, a := range albums {
		out = append(out, a.Slug)
	}
	return out
}

// testStore builds a store directly, without YAML. Newest-first order is: b d c a e.
func testStore() *MemStore {
	return newMemStore(Site{}, []Album{
		{Slug: "a", Date: date("2021-01-01"), Tags: []string{"x"}},
		{Slug: "b", Date: date("2024-01-01"), Tags: []string{"x", "y"}, Featured: true},
		{Slug: "c", Date: date("2022-01-01")},
		{Slug: "d", Date: date("2023-01-01"), Tags: []string{"y"}},
		{Slug: "e", Date: date("2020-01-01"), Tags: []string{"x"}},
	})
}

func TestAlbums(t *testing.T) {
	tests := []struct {
		name     string
		q        AlbumQuery
		want     []string
		wantNext int
	}{
		{"first page", AlbumQuery{Page: 1, PerPage: 2}, []string{"b", "d"}, 2},
		{"middle page", AlbumQuery{Page: 2, PerPage: 2}, []string{"c", "a"}, 3},
		{"last partial page", AlbumQuery{Page: 3, PerPage: 2}, []string{"e"}, 0},
		{"past the end", AlbumQuery{Page: 4, PerPage: 2}, []string{}, 0},
		{"exact fit has no next", AlbumQuery{Page: 1, PerPage: 5}, []string{"b", "d", "c", "a", "e"}, 0},
		{"page 0 treated as 1", AlbumQuery{Page: 0, PerPage: 2}, []string{"b", "d"}, 2},
		{"negative page treated as 1", AlbumQuery{Page: -5, PerPage: 2}, []string{"b", "d"}, 2},
		{"tag filter", AlbumQuery{Tag: "x", Page: 1, PerPage: 10}, []string{"b", "a", "e"}, 0},
		{"tag filter paged", AlbumQuery{Tag: "x", Page: 1, PerPage: 2}, []string{"b", "a"}, 2},
		{"tag filter page 2", AlbumQuery{Tag: "x", Page: 2, PerPage: 2}, []string{"e"}, 0},
		{"unknown tag", AlbumQuery{Tag: "nope", Page: 1, PerPage: 2}, []string{}, 0},
	}
	s := testStore()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := s.Albums(tc.q)
			if !slices.Equal(slugs(got.Albums), tc.want) {
				t.Errorf("albums = %v, want %v", slugs(got.Albums), tc.want)
			}
			if got.Next != tc.wantNext {
				t.Errorf("next = %d, want %d", got.Next, tc.wantNext)
			}
		})
	}
}

func TestAlbumsPanicsOnBadPerPage(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic for PerPage=0")
		}
	}()
	testStore().Albums(AlbumQuery{Page: 1})
}

func TestStoreLookups(t *testing.T) {
	s := testStore()
	if got := slugs(s.Featured()); !slices.Equal(got, []string{"b"}) {
		t.Errorf("Featured() = %v, want [b]", got)
	}
	if got := s.Tags(); !slices.Equal(got, []string{"x", "y"}) {
		t.Errorf("Tags() = %v, want [x y]", got)
	}
	if _, ok := s.Album("c"); !ok {
		t.Error(`Album("c") not found`)
	}
	if _, ok := s.Album("zzz"); ok {
		t.Error(`Album("zzz") should not be found`)
	}
}
