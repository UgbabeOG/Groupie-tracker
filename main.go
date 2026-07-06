package main

import (
	"log"
	"net/http"
	"time"

	"acad.learn2earn.ng/git/gugbabe/groupie-tracker/utils"
)

func main() {

	// ------------------------------------
	// Initialize application cache
	// ------------------------------------

	// Load data from local JSON cache or live API.
	utils.InitCache()

	// Start background synchronization worker.
	utils.StartSyncWorker(24 * time.Hour)

	// ------------------------------------
	// Router configuration
	// ------------------------------------

	mux := http.NewServeMux()

	// Serve static assets.
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Application routes.
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/artist", ArtistDetailHandler)

	// ------------------------------------
	// HTTP server configuration
	// ------------------------------------

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}