package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireOriginAuth(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	tests := []struct {
		name       string
		secrets    []string
		path       string
		header     string // "" sends no X-Origin-Auth header
		wantStatus int
	}{
		{"current secret", []string{"cur", "prev"}, "/", "cur", http.StatusOK},
		{"previous secret", []string{"cur", "prev"}, "/", "prev", http.StatusOK},
		{"wrong secret", []string{"cur", "prev"}, "/", "nope", http.StatusForbidden},
		{"prefix of secret", []string{"cur"}, "/", "cu", http.StatusForbidden},
		{"no header", []string{"cur"}, "/", "", http.StatusForbidden},
		{"no secrets locks everything", nil, "/", "anything", http.StatusForbidden},
		{"empty secret never matches", []string{""}, "/", "", http.StatusForbidden},
		{"healthz skips auth", []string{"cur"}, "/healthz", "", http.StatusOK},
		{"readyz skips auth", []string{"cur"}, "/readyz", "", http.StatusOK},
		{"probe prefix does not skip auth", []string{"cur"}, "/healthz/../albums", "", http.StatusForbidden},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.header != "" {
				req.Header.Set("X-Origin-Auth", tc.header)
			}
			rec := httptest.NewRecorder()
			RequireOriginAuth(tc.secrets, ok).ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}
