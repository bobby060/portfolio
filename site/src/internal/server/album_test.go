package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAlbumPage(t *testing.T) {
	h := newTestServer(t)
	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}

	rec := get("/albums/jinju-lights")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	for _, want := range []string{
		"<h1",
		"Jinju Lights",
		`alt="Lanterns on the Namgang River"`, // caption used as alt text
		`alt="Jinju Lights, photo 2"`,         // fallback alt when there's no caption
		`data-pswp-width="1920"`,              // 6000px original -> largest 1920 variant
		`data-pswp-height="1280"`,
		`src="/static/js/gallery.js"`,
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("album page does not contain %q", want)
		}
	}

	for _, path := range []string{"/albums/nope", "/albums/jinju-lights/extra"} {
		if rec := get(path); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, want 404", path, rec.Code)
		}
	}
}
