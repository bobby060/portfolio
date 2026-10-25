package server

import (
	"net/url"
	"strconv"

	"portfolio/internal/gallery"
	"portfolio/web/components"
)

// maxPage caps ?page= so a huge value can't overflow page*albumsPerPage.
const maxPage = 1000

// parsePage reads ?page=, which is user input: anything invalid becomes 1.
func parsePage(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 1
	}
	return min(n, maxPage)
}

// albumsURL builds path?page=&tag=, leaving out defaults so each view has one canonical URL.
func albumsURL(path, tag string, page int) string {
	v := url.Values{}
	if tag != "" {
		v.Set("tag", tag)
	}
	if page > 1 {
		v.Set("page", strconv.Itoa(page))
	}
	if len(v) == 0 {
		return path
	}
	return path + "?" + v.Encode()
}

// tagChips lists every tag. While a tag is active, only that chip shows, and it clears the filter.
func tagChips(all []string, active string) []components.TagChip {
	if active != "" {
		return []components.TagChip{{Name: active, Href: "/#albums", HX: "/", Active: true}}
	}
	chips := make([]components.TagChip, 0, len(all))
	for _, t := range all {
		u := albumsURL("/albums", t, 1)
		chips = append(chips, components.TagChip{Name: t, Href: u + "#albums", HX: u})
	}
	return chips
}

// albumGrid returns album cards for one page of the tag filter.
// With cumulative=true it returns pages 1..page instead, so a full-page load of
// /albums?page=3 shows everything the visitor had scrolled through.
func (s *Server) albumGrid(tag string, page int, cumulative bool) components.AlbumGrid {
	q := gallery.AlbumQuery{Tag: tag, Page: page, PerPage: albumsPerPage}
	if cumulative {
		q = gallery.AlbumQuery{Tag: tag, Page: 1, PerPage: page * albumsPerPage}
	}
	res := s.store.Albums(q)

	var g components.AlbumGrid
	for _, a := range res.Albums {
		g.Albums = append(g.Albums, s.albumCard(a, sizesGrid))
	}
	if res.Next != 0 {
		g.MoreHref = albumsURL("/albums", tag, page+1) + "#albums"
		g.MoreHX = albumsURL("/fragments/albums", tag, page+1)
	}
	return g
}
