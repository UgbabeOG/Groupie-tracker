package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	port := ":5000"
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Server is active") })

	log.Printf("Server is running on port %v", port)
	http.ListenAndServe(port, mux)
}
