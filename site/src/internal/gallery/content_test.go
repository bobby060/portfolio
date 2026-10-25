// External test package: this test uses gallery only through its public API,
// the same way the server will.
package gallery_test

import (
	"testing"

	"portfolio/content"
	"portfolio/internal/gallery"
)

func TestRealContentLoads(t *testing.T) {
	s, err := gallery.LoadYAML(content.FS)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Site().Hero) == 0 {
		t.Error("site.yaml has no hero images")
	}
	if len(s.Featured()) == 0 {
		t.Error("no featured albums")
	}
}
