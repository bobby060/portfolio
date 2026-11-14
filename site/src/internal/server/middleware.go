package server

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// probePaths skip origin auth: ACA's health probes call the container directly, without the header.
var probePaths = map[string]bool{"/healthz": true, "/readyz": true}

// RequireOriginAuth rejects requests whose X-Origin-Auth header matches none of secrets.
// With no secrets it rejects everything, so a wiring mistake locks the site instead of opening it.
func RequireOriginAuth(secrets []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if probePaths[r.URL.Path] || validOrigin(r.Header.Get("X-Origin-Auth"), secrets) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "forbidden", http.StatusForbidden)
	})
}

// validOrigin compares in constant time and checks every secret, so response
// timing reveals neither how much of a guess matched nor which slot matched.
func validOrigin(got string, secrets []string) bool {
	if len(got) == 0 {
		return false
	}
	ok := 0
	for _, s := range secrets {
		ok |= subtle.ConstantTimeCompare([]byte(got), []byte(s))
	}
	return ok == 1
}

// csp is joined once at startup; each entry is one directive.
var csp = strings.Join([]string{
	"default-src 'self'",
	"base-uri 'self'",
	"frame-ancestors 'none'",
	"object-src 'none'",
	"img-src 'self' https://picsum.photos https://fastly.picsum.photos", // picsum redirects to fastly; drop both in phase 2
	"style-src 'self' 'unsafe-inline'",
}, "; ")

// SecurityHeaders sets the browser hardening headers on every response (DESIGN.md §3.5).
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("String-Transport-Security", "max-age=31546000; includeSubDomains")
		h.Add("X-Content-Type-Options", "nosniff")
		h.Add("Refferrer-Policy", "strict-origin-when-cross-origin")
		h.Add("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		h.Add("Content-Security-Policy", csp)
		next.ServeHTTP(w, r)
	})
}
