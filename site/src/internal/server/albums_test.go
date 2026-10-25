package server

import "testing"

func TestParsePage(t *testing.T) {
	tests := map[string]int{
		"": 1, "abc": 1, "0": 1, "-3": 1, "2": 2,
		"1000": 1000, "1001": 1000, "99999999999999999999": 1,
	}
	for in, want := range tests {
		if got := parsePage(in); got != want {
			t.Errorf("parsePage(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestAlbumsURL(t *testing.T) {
	tests := []struct {
		tag  string
		page int
		want string
	}{
		{"", 1, "/albums"},
		{"", 3, "/albums?page=3"},
		{"wildlife", 1, "/albums?tag=wildlife"},
		{"a b&c", 2, "/albums?page=2&tag=a+b%26c"},
	}
	for _, tc := range tests {
		if got := albumsURL("/albums", tc.tag, tc.page); got != tc.want {
			t.Errorf("albumsURL(%q, %d) = %q, want %q", tc.tag, tc.page, got, tc.want)
		}
	}
}
