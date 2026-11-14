package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"portfolio/content"
	"portfolio/internal/config"
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
	cfg, err := config.FromEnv(os.Getenv)
	if err != nil {
		return err
	}

	var contentFS fs.FS = content.FS
	if cfg.ContentDir != "" {
		contentFS = os.DirFS(cfg.ContentDir) // dev: YAML edits apply on restart, no rebuild
		log.Info("reading content from disk", "dir", cfg.ContentDir)
	}

	store, err := gallery.LoadYAML(contentFS)
	if err != nil {
		return fmt.Errorf("load content: %w", err)
	}

	var h http.Handler = server.New(store, images.Picsum{}, log)
	if cfg.DevMode {
		log.Warn("DEV_MODE=true: origin auth is OFF")
	} else {
		h = server.RequireOriginAuth(cfg.OriginSecrets, h)
	}
	h = server.SecurityHeaders(h) // outermost, so 403s get the headers too

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second, // drop clients that never finish sending headers
	}
	// ...rest unchanged

	log.Info("listening", "addr", srv.Addr)
	return srv.ListenAndServe()
}
