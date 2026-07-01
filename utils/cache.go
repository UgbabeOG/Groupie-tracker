package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// apiResource struct defines the structure for each API resource.
type apiResource struct {
	url      string
	filename string
	target   any
}

// API URLs for each resource
var (
	ArtistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	DatesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsURL = "https://groupietrackers.herokuapp.com/api/relation"
)

// Cache variables to hold the in-memory representations of the API data
var (
	ArtistsCache   []Artist
	LocationsCache []LocationData
	DatesCache     []DateData
	RelationsCache []RelationData
	CacheLock      sync.RWMutex
)

// getResourcesConfig returns a slice of apiResource structs, each containing the URL, filename, and target variable for the respective API resources.
func getResourcesConfig() []apiResource {
	return []apiResource{
		{url: ArtistsURL, filename: "artists.json", target: &ArtistsCache},
		{url: LocationsURL, filename: "locations.json", target: &LocationsCache},
		{url: DatesURL, filename: "dates.json", target: &DatesCache},
		{url: RelationsURL, filename: "relations.json", target: &RelationsCache},
	}
}

// InitCache initializes the local cache by ensuring the data directory exists, downloading any missing files, and loading all data into memory.
func InitCache() error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("Failed to create data folder: %w", err)
	}

	for _, resource := range getResourcesConfig() {
		if _, err := os.Stat(filepath.Join("data", resource.filename)); os.IsNotExist(err) {
			log.Printf("Missing storage track: %s. Performing initial catch up...\n", resource.filename)
			if err := downloadToFile(resource.url, resource.filename); err != nil {
				return err
			}
		}
	}
	return loadAllFromDiskToMemory()
}

// loadAllFromDiskToMemory loads the JSON data from local files into memory for all configured API resources.
func loadAllFromDiskToMemory() error {
	CacheLock.Lock()
	defer CacheLock.Unlock()

	for _, resource := range getResourcesConfig() {
		data, err := os.ReadFile(filepath.Join("data", resource.filename))
		if err != nil {
			return fmt.Errorf("Failed reading cache target %s: %w", resource.filename, err)
		}

		if err := json.Unmarshal(data, resource.target); err != nil {
			return fmt.Errorf("Failed unmarshaling %s: %w", resource.filename, err)
		}
	}
	return nil
}

// downloadToFile downloads the JSON data from the specified URL and saves it to a local file with the given filename.
func downloadToFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Network exception on %s: %w", filename, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with bad status code %d for %s", response.StatusCode, err)
	}

	finalPath := filepath.Join("data", filename)
	tmpPath := finalPath + ".tmp"
	var processedJSON []byte

	if filename == "artists.json" {
		var raw json.RawMessage

		if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
			return fmt.Errorf("Parsing payload error on %s: %w", filename, err)
		}

		processedJSON, err = json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return fmt.Errorf("Failed to marshal artists data: %w", err)
		}
	} else {
		var wrapper struct {
			Index json.RawMessage `json:"index"`
		}
		if err := json.NewDecoder(response.Body).Decode(&wrapper); err != nil {
			return fmt.Errorf("Parsing wrapper error on %s: %w", filename, err)
		}

		processedJSON, err = json.MarshalIndent(wrapper.Index, "", "  ")
		if err != nil {
			return fmt.Errorf("Failed to marshal %s data: %w", filename, err)
		}
	}

	if err := os.WriteFile(tmpPath, processedJSON, 0644); err != nil {
		return fmt.Errorf("Failed temporary track save: %w", err)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.ReadDir(tmpPath)
		return fmt.Errorf("Failed atomic engine rename swap: %w", err)
	}
	return nil
}

// StartSyncWorker starts a background worker that periodically synchronizes the local cache with the API data at the specified interval.
func StartSyncWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			log.Println("Starting cache synchronization from API...")

			failed := false
			for _, resource := range getResourcesConfig() {
				if err := downloadToFile(resource.url, resource.filename); err != nil {
					log.Printf("Sync Worker Error downloading %s: %v\n", resource.filename, err)
					failed = true
					break
				}
			}

			if failed {
				continue
			}

			if err := loadAllFromDiskToMemory(); err != nil {
				log.Printf("Sync Worker Error (Memory Swap Fault): %v\n", err)
			} else {
				log.Println("Background memory cache refreshed safely from updated files.")
			}
		}
	}()
}
