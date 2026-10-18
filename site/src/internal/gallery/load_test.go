package gallery

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

const siteYAML = `
title: Test Site
tagline: hello
hero:
  - { sha: h1, width: 1920, height: 1280 }
`

func mapFS(files map[string]string) fstest.MapFS {
	fsys := fstest.MapFS{}
	for p, data := range files {
		fsys[p] = &fstest.MapFile{Data: []byte(data)}
	}
	return fsys
}

func TestLoadYAML(t *testing.T) {
	fsys := mapFS(map[string]string{
		"site.yaml": siteYAML,
		"albums/old/album.yaml": `
title: Old
date: 2020-01-01
tags: [b, a]
images:
  - { sha: o1, width: 600, height: 400 }
  - { sha: o2, width: 400, height: 600 }
`,
		"albums/new-dir/album.yaml": `
title: New
slug: custom
date: 2024-06-01
featured: true
cover: n2
tags: [a, c]
images:
  - { sha: n1, width: 600, height: 400 }
  - sha: n2
    width: 500
    height: 500
    taken_at: 2024-06-01T09:12:00Z
    caption: hi
`,
	})

	s, err := LoadYAML(fsys)
	if err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	if got := s.Site().Title; got != "Test Site" {
		t.Errorf("site title = %q", got)
	}

	old, ok := s.Album("old")
	if !ok {
		t.Fatal("slug should default to the directory name")
	}
	if old.Cover != "o1" {
		t.Errorf("omitted cover should default to first image, got %q", old.Cover)
	}

	if _, ok := s.Album("new-dir"); ok {
		t.Error("explicit slug should replace the directory name")
	}
	nw, ok := s.Album("custom")
	if !ok {
		t.Fatal(`album "custom" not found`)
	}
	if nw.Cover != "n2" {
		t.Errorf("cover = %q, want n2", nw.Cover)
	}
	if img := nw.Images[1]; img.Caption != "hi" || img.TakenAt.IsZero() {
		t.Errorf("optional image fields not decoded: %+v", img)
	}

	if got := slugs(s.Featured()); !slices.Equal(got, []string{"custom"}) {
		t.Errorf("Featured() = %v", got)
	}
	if got := s.Tags(); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("Tags() = %v", got)
	}
}

func TestLoadYAMLErrors(t *testing.T) {
	const img = "images:\n  - { sha: x, width: 1, height: 1 }\n"
	const ok = "title: T\ndate: 2024-01-01\n" + img
	const a = "albums/a/album.yaml"

	tests := []struct {
		name     string
		files    map[string]string
		wantPath string // the error must name the offending file
		wantMsg  string
	}{
		{"missing site.yaml", map[string]string{a: ok}, "site.yaml", ""},
		{"duplicate slug", map[string]string{"site.yaml": siteYAML, a: ok, "albums/b/album.yaml": "slug: a\n" + ok},
			"albums/b/album.yaml", "duplicate slug"},
		{"unknown cover", map[string]string{"site.yaml": siteYAML, a: "cover: nope\n" + ok}, a, `cover "nope" not found`},
		{"zero width", map[string]string{"site.yaml": siteYAML, a: "title: T\ndate: 2024-01-01\nimages:\n  - { sha: x, width: 0, height: 1 }\n"},
			a, "width and height must be > 0"},
		{"missing title", map[string]string{"site.yaml": siteYAML, a: "date: 2024-01-01\n" + img}, a, "title is required"},
		{"missing date", map[string]string{"site.yaml": siteYAML, a: "title: T\n" + img}, a, "date is required"},
		{"impossible date", map[string]string{"site.yaml": siteYAML, a: "title: T\ndate: 2026-04-32\n" + img}, a, ""},
		{"typo in field name", map[string]string{"site.yaml": siteYAML, a: "featued: true\n" + ok}, a, "featued"},
		{"no images", map[string]string{"site.yaml": siteYAML, a: "title: T\ndate: 2024-01-01\n"}, a, "at least one image"},
		{"missing sha", map[string]string{"site.yaml": siteYAML, a: "title: T\ndate: 2024-01-01\nimages:\n  - { width: 1, height: 1 }\n"},
			a, "sha is required"},
		{"duplicate image", map[string]string{"site.yaml": siteYAML, a: ok + "  - { sha: x, width: 1, height: 1 }\n"},
			a, "duplicate image"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadYAML(mapFS(tc.files))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantPath) {
				t.Errorf("error %q does not name file %q", err, tc.wantPath)
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Errorf("error %q does not contain %q", err, tc.wantMsg)
			}
		})
	}
}
