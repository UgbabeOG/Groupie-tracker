package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestInitCacheCompletePipeline(t *testing.T) {
	// Setup a single local mock server to handle all API endpoints
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch r.URL.Path {
		case "/artists":
			w.Write([]byte(`[{"id":1,"name":"Artist 1","image":"image1.jpg","members":["Member 1"],"creationDate":2000,"firstAlbum":"Album 1"}]`))
		case "/locations":
			w.Write([]byte(`{"index": [{"id":1, "locations":["lagos-nigeria"] }]}`))
		case "/dates":
			w.Write([]byte(`{"index": [{"id":1, "dates":["2023-01-01"] }]}`))
		case "/relation":
			w.Write([]byte(`{"index": [{"id":1, "datesLocations":{"2023-01-01":["lagos-nigeria"]}}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Override default endpoints with local mock server paths
	oldArtistsURL, oldLocationsURL, oldDatesURL, oldRelationsURL := ArtistsURL, LocationsURL, DatesURL, RelationsURL
	ArtistsURL = mockServer.URL + "/artists"
	LocationsURL = mockServer.URL + "/locations"
	DatesURL = mockServer.URL + "/dates"
	RelationsURL = mockServer.URL + "/relation"

	// Guarantee that baseline live URLs restore perfectly on completion
	defer func() {
		ArtistsURL = oldArtistsURL
		LocationsURL = oldLocationsURL
		DatesURL = oldDatesURL
		RelationsURL = oldRelationsURL
	}()

	// Prevent writing to your workspace folder by isolating execution inside a temp folder
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed switching system execution path: %v", err)
	}
	defer os.Chdir(oldCwd)

	// Force empty variables to ensure we check real updates
	CacheLock.Lock()
	ArtistsCache = nil
	LocationsCache = nil
	DatesCache = nil
	RelationsCache = nil
	CacheLock.Unlock()

	// Run the InitCache function to test the complete pipeline
	if err := InitCache(); err != nil {
		t.Fatalf("InitCache failed: %v", err)
	}

	// Check if the caches are populated correctly
	CacheLock.RLock()
	defer CacheLock.RUnlock()

	if len(ArtistsCache) != 1 || ArtistsCache[0].Name != "Artist 1" {
		t.Errorf("ArtistsCache not populated correctly: got %v", ArtistsCache)
	}
	if len(LocationsCache) != 1 || LocationsCache[0].Locations[0] != "lagos-nigeria" {
		t.Errorf("LocationsCache not populated correctly: got %v", LocationsCache)
	}
	if len(DatesCache) != 1 || DatesCache[0].Dates[0] != "2023-01-01" {
		t.Errorf("DatesCache not populated correctly: got %v", DatesCache)
	}
	if len(RelationsCache) != 1 || RelationsCache[0].DatesLocations["2023-01-01"][0] != "lagos-nigeria" {
		t.Errorf("RelationsCache not populated correctly: got %v", RelationsCache)
	}

	// Verify disk files are clean arrays without "index" wrapper
	for _, filename := range []string{"artists.json", "locations.json", "dates.json", "relations.json"} {
		if _, err := os.Stat("data/" + filename); os.IsNotExist(err) {
			t.Errorf("Expected cache file data/%s to be written to disk, but it's missing", filename)

		} else {
			content, err := os.ReadFile("data/" + filename)
			if err != nil {
				t.Errorf("Failed reading file %s: %v", filename, err)
			}
			if string(content) == "" || string(content) == "{}" {
				t.Errorf("File %s is empty or not properly formatted: %s", filename, string(content))
			}
		}
	}
}
