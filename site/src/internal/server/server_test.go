package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"portfolio/content"
	"portfolio/internal/gallery"
	"portfolio/internal/images"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	store, err := gallery.LoadYAML(content.FS)
	if err != nil {
		t.Fatal(err)
	}
	return New(store, images.Picsum{}, slog.New(slog.DiscardHandler))
}

func TestRoutes(t *testing.T) {
	h := newTestServer(t)
	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/healthz", http.StatusOK, "ok"},
		{"/", http.StatusOK, "Robert Norwood"},
		{"/no-such-page", http.StatusNotFound, "Oops!"},
		{"/static/img/favicon.png", http.StatusOK, ""},
		{"/static/img/", http.StatusNotFound, ""},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("body does not contain %q", tc.wantBody)
			}
		})
	}
}
