package main

import (
	"log"
	"net/http"
	"strconv"

	"acad.learn2earn.ng/git/gugbabe/groupie-tracker/utils"
)

// HomeHandler displays the homepage with all artists and their related data.
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	// Only allow the root path.
	if r.URL.Path != "/" {
		utils.RenderError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// Read cache safely.
	utils.CacheLock.RLock()
	defer utils.CacheLock.RUnlock()

	var artists []utils.FullArtistData

	// Build a complete artist object for every artist.
	for _, artist := range utils.ArtistsCache {

		full := utils.FullArtistData{
			Artist: artist,
		}

		// Attach locations.
		for _, location := range utils.LocationsCache {
			if location.ID == artist.ID {
				full.Locations = location.Locations
				break
			}
		}

		// Attach concert dates.
		for _, date := range utils.DatesCache {
			if date.ID == artist.ID {
				full.Dates = date.Dates
				break
			}
		}

		// Attach relations.
		for _, relation := range utils.RelationsCache {
			if relation.ID == artist.ID {
				full.Relations = relation.DatesLocations
				break
			}
		}

		artists = append(artists, full)
	}

	// Render homepage.
	if err := utils.RenderTemplate(w, "index.html", artists); err != nil {
		log.Printf("[HomeHandler] %v", err)
		utils.RenderError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// ArtistDetailHandler displays a single artist with all related information.
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {

	// Read artist ID.
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.RenderError(w, http.StatusBadRequest, "Missing artist ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RenderError(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	utils.CacheLock.RLock()
	defer utils.CacheLock.RUnlock()

	var detail utils.FullArtistDetail
	found := false

	// Find artist.
	for _, artist := range utils.ArtistsCache {
		if artist.ID == id {
			detail.Artist = artist
			found = true
			break
		}
	}

	if !found {
		utils.RenderError(w, http.StatusNotFound, "Artist Not Found")
		return
	}

	// Attach locations.
	for _, location := range utils.LocationsCache {
		if location.ID == id {
			detail.Locations = location.Locations
			break
		}
	}

	// Attach dates.
	for _, date := range utils.DatesCache {
		if date.ID == id {
			detail.Dates = date.Dates
			break
		}
	}

	// Attach relations.
	for _, relation := range utils.RelationsCache {
		if relation.ID == id {
			detail.DatesLocations = relation.DatesLocations
			break
		}
	}

	// Render detail page.
	if err := utils.RenderTemplate(w, "artistDetail.html", detail); err != nil {
		log.Printf("[ArtistDetailHandler] %v", err)
		utils.RenderError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
