// Package server wires HTTP routes to handlers and renders pages.
package server

import (
	"log/slog"
	"net/http"
	"strings"

	"portfolio/internal/gallery"
	"portfolio/internal/images"
	"portfolio/web/static"
)

// Server holds what handlers need. Everything arrives through New, so main picks
// the real implementations and tests can pass fakes.
type Server struct {
	store gallery.Store
	imgs  images.URLBuilder
	log   *slog.Logger
}

// New returns the site's HTTP handler.
func New(store gallery.Store, imgs images.URLBuilder, log *slog.Logger) http.Handler {
	s := &Server{store: store, imgs: imgs, log: log}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", noListing(http.FileServerFS(static.FS))))
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("/", s.notFound)
	return mux
}

// noListing hides directory listings, so /static/css/ can't be browsed.
func noListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
