package main

import (
	"log"
	"net/http"

	"portfolio/internal/web"
)

func main() {
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", web.NewRouter()))
}
