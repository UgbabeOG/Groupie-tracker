package main

import (
	"log"
	"net/http"
	"strconv"

	"acad.learn2earn.ng/git/gugbabe/groupie-tracker/utils"
)

// HomeHandler serves the root path "/". It reads the active cache array of
// artists and feeds it into the main grid template (index.html).
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.CacheLock.Lock()
	defer utils.CacheLock.Unlock()

	artists := utils.ArtistsCache

	if err := utils.RenderTemplate(w, "index.html", artists); err != nil {
		log.Printf("[handler] ERROR rendering index: %v", err)
		utils.RenderError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// ArtistDetailHandler serves the "/artist" path. It collects the "id" query
// parameter, compiles a unified FullArtistDetail object from the memory cache
// blocks, and passes it to the profile view (artistDetail.html).
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	utils.CacheLock.Lock()
	defer utils.CacheLock.Unlock()
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.RenderError(w, http.StatusBadRequest, "Bad Request: missing 'id' parameter")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RenderError(w, http.StatusBadRequest, "Bad Request: 'id' must be a number")
		return
	}
	detail := utils.FullArtistDetail{}

	for _, artist := range utils.ArtistsCache {
		if artist.ID == id {
			detail.ID = artist.ID
			detail.Name = artist.Name
			detail.Image = artist.Image
			detail.Members = artist.Members
			detail.CreationDate = artist.CreationDate
			detail.FirstAlbum = artist.FirstAlbum
			break
		}
	}

	if err := utils.RenderTemplate(w, "artistDetail.html", detail); err != nil {
		log.Printf("[handler] ERROR rendering artist detail: %v", err)
		utils.RenderError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
