package utils

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type DateData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type LocationData struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}
type RelationData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type FullArtistData struct {
	Artist
	Locations []string
	Dates     []string
	Relations map[string][]string
}

type FullArtistDetail struct {
	Artist
	Locations      []string
	Dates          []string
	DatesLocations map[string][]string
}

type LocationsWrapper struct {
	Index []LocationData `json:"index"`
}

type DatesWrapper struct {
	Index []DateData `json:"index"`
}

type RelationsWrapper struct {
	Index []RelationData `json:"index"`
}

type ErrorData struct {
	Code    int
	Message string
}
