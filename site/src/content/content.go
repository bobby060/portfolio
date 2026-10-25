// Package content holds album metadata as YAML, embedded in the binary
//
// (for now)
package content

import "embed"

// FS contains site.yaml and albums/<slug>/album.yaml
//
//go:embed site.yaml albums
var FS embed.FS
