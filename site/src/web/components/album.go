package components

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// cssURLEscaper makes a URL safe inside an unquoted CSS url(...). Parens, quotes,
// whitespace and backslashes would otherwise end the URL and allow CSS injection.
// The URL is unquoted because templ HTML-escapes style values, which mangles quotes.
var cssURLEscaper = strings.NewReplacer(
	"(", "%28", ")", "%29", `"`, "%22", "'", "%27",
	`\`, "%5C", " ", "%20", "\t", "%09", "\n", "", "\r", "",
)

// headerStyle sets the album header's background image, and a height that follows
// the cover's aspect ratio, clamped to 200–400px like the old site.
//
// templ.SafeCSS bypasses templ's CSS sanitizer, which would reject url(...).
// That's safe because we escape u ourselves and w and h are integers.
func headerStyle(u string, w, h int) templ.SafeCSS {
	return templ.SafeCSS(fmt.Sprintf(
		`background-image:url(%s);height:clamp(200px,calc(100vw * %d / %d),400px)`,
		cssURLEscaper.Replace(u), h, w,
	))
}
