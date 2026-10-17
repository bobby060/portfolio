package web

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		Page("Home Page").Render(r.Context(), w)
	})

	mux.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		Content("Successfully retrieved dynamically over HTMX!").Render(r.Context(), w)
	})

	return mux
}
