package gallery

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"

	yaml "go.yaml.in/yaml/v3"
)

// LoadYAML reads site.yaml and albums/*/album.yaml from fsys, validates them,
// and returns a ready-to-query store. Every error names the offending file.
func LoadYAML(fsys fs.FS) (*MemStore, error) {
	var site Site
	err := decodeStrict(fsys, "site.yaml", &site)
	if err != nil {
		return nil, err
	}

	var albums []Album

	var albumsMap = make(map[string]string)

	albumFiles, err := fs.Glob(fsys, "albums/*/album.yaml")

	if err != nil {
		return nil, err
	}

	for _, albumFile := range albumFiles {
		var album Album
		err := decodeStrict(fsys, albumFile, &album)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", albumFile, err)
		}
		if album.Slug == "" {
			album.Slug = path.Base(path.Dir(albumFile))
		}

		err = validateAlbum(&album)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", albumFile, err)
		}

		other, exists := albumsMap[album.Slug]
		if exists {
			return nil, fmt.Errorf("%s: duplicate slug %q (also used by %s)",
				albumFile, album.Slug, other)
		}
		albumsMap[album.Slug] = albumFile

		albums = append(albums, album)

	}

	return newMemStore(site, albums), nil
}

func validateAlbum(album *Album) error {
	if album.Title == "" {
		return errors.New("title is required")
	}

	if album.Date.IsZero() {
		return errors.New("date is required")
	}

	if len(album.Images) == 0 {
		return errors.New("at least one image")
	}

	if album.Cover == "" {
		album.Cover = album.Images[0].SHA
	}

	validCover := false
	imagesMap := make(map[string]bool)

	for _, i := range album.Images {

		if i.SHA == album.Cover {
			validCover = true
		}

		if i.Height <= 0 || i.Width <= 0 {
			return errors.New("width and height must be > 0")

		}

		if i.SHA == "" {
			return errors.New("sha is required")
		}

		_, exists := imagesMap[i.SHA]
		if exists {
			return fmt.Errorf("duplicate image %q", i.SHA)
		}
		imagesMap[i.SHA] = true
	}

	if !validCover {
		return fmt.Errorf("cover %q not found in images", album.Cover)
	}

	return nil
}

func decodeStrict(fsys fs.FS, name string, out any) error {
	data, err := fs.ReadFile(fsys, name)

	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	decoder.KnownFields(true)
	if err := decoder.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("empty file")
		}
		return err
	}
	return nil

}
