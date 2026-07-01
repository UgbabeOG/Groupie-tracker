package main

import (
	"log"
	"net/http"
	"time"

	"acad.learn2earn.ng/git/gugbabe/groupie-tracker/utils"
)

func main() {
	// Bootstrap: load all API entities into memory from disk or live API.
	utils.InitCache()

	// Launch a background goroutine that periodically refreshes the cache.
	utils.StartSyncWorker(24 * time.Hour)

	mux := http.NewServeMux()

	// Mount the static file server for CSS, JS, and other client assets.
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Bind route handlers.
	mux.HandleFunc("/{$}", HomeHandler)
	mux.HandleFunc("/artist", ArtistDetailHandler)

	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
