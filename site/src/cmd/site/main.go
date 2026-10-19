package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"portfolio/content"
	"portfolio/internal/gallery"
	"portfolio/internal/images"
	"portfolio/internal/server"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("startup failed", "err", err)
		os.Exit(1)
	}
}

// run wires the app together. Swapping picsum for real images is a one-line change here.
func run(log *slog.Logger) error {
	var contentFS fs.FS = content.FS
	if dir := os.Getenv("CONTENT_DIR"); dir != "" {
		contentFS = os.DirFS(dir) // dev: YAML edits apply on restart, no rebuild
		log.Info("reading content from disk", "dir", dir)
	}

	store, err := gallery.LoadYAML(contentFS)
	if err != nil {
		return fmt.Errorf("load content: %w", err)
	}

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           server.New(store, images.Picsum{}, log),
		ReadHeaderTimeout: 10 * time.Second, // drop clients that never finish sending headers
	}
	log.Info("listening", "addr", srv.Addr)
	return srv.ListenAndServe()
}
