package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"acad.learn2earn.ng/git/gugbabe/groupie-tracker/utils"
)

func TestHomeHandler_RendersArtistList(t *testing.T) {
    utils.CacheLock.Lock()
    utils.ArtistsCache = []utils.Artist{{
        ID:           1,
        Name:         "Test Band",
        Image:        "https://example.com/test.png",
        Members:      []string{"Alice", "Bob"},
        CreationDate: 1990,
        FirstAlbum:   "First Record",
    }}
    utils.CacheLock.Unlock()

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rr := httptest.NewRecorder()

    HomeHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
    }

    body := rr.Body.String()
    if !strings.Contains(body, "Test Band") {
        t.Errorf("response body does not contain artist name: %s", body)
    }
    if !strings.Contains(body, "/artist?id=1") {
        t.Errorf("response body does not contain expected artist link: %s", body)
    }
}

func TestArtistDetailHandler_MissingIDReturnsBadRequest(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/artist", nil)
    rr := httptest.NewRecorder()

    ArtistDetailHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
    }
    if !strings.Contains(rr.Body.String(), "missing &#39;id&#39; parameter") {
        t.Errorf("expected missing id message, got %q", rr.Body.String())
    }
}

func TestArtistDetailHandler_InvalidIDReturnsBadRequest(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/artist?id=bad", nil)
    rr := httptest.NewRecorder()

    ArtistDetailHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
    }
    if !strings.Contains(rr.Body.String(), "&#39;id&#39; must be a number") {
        t.Errorf("expected invalid id message, got %q", rr.Body.String())
    }
}

func TestArtistDetailHandler_RendersArtistDetail(t *testing.T) {
    utils.CacheLock.Lock()
    utils.ArtistsCache = []utils.Artist{{
        ID:           42,
        Name:         "Example Band",
        Image:        "https://example.com/example.png",
        Members:      []string{"One", "Two"},
        CreationDate: 1985,
        FirstAlbum:   "Debut",
    }}
    utils.CacheLock.Unlock()

    req := httptest.NewRequest(http.MethodGet, "/artist?id=42", nil)
    rr := httptest.NewRecorder()

    ArtistDetailHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
    }

    body := rr.Body.String()
    if !strings.Contains(body, "Example Band") {
        t.Errorf("expected response to contain artist name, got %q", body)
    }
    if !strings.Contains(body, "Founded 1985") {
        t.Errorf("expected response to contain creation date, got %q", body)
    }
    if !strings.Contains(body, "First Album: Debut") {
        t.Errorf("expected response to contain first album, got %q", body)
    }
}
