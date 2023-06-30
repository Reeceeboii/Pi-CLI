package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	// Sample API key for test case usage.
	testKey = "c808f484a4e88cc32a9a8bfcce19169c77bcd9c5eec18d859e1bb4b318bf42bf"
)

// Tests for api.Summary.Update() with an API key
func TestUpdateWithApiKey(t *testing.T) {
	summary := NewSummary()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure URL is formatted with the correct query string.
		if !strings.Contains(r.URL.RequestURI(), "/api.php?summary&auth="+testKey) {
			t.Error("@TestUpdateWithApiKey: api.Summary.Update() did not request the expected Pi Hole api endpoint with expected API key.")
		}
	}))
	defer mockServer.Close()
	url := mockServer.URL + "/api.php"

	summary.Update(url, testKey, nil)
}

// Tests for api.Summary.Update() without an API key
func TestUpdateWithoutAPIKey(t *testing.T) {
	summary := NewSummary()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure URL does not contain &auth as part of its query string.
		if strings.Contains(r.URL.RequestURI(), "/api.php?summary&auth=") {
			t.Error("@TestUpdateWithoutAPIKey: api.Summary.Update() did not send the expected query string when calling with an API key: /api.php?summary&auth=")
		}
		// Ensure URL is formatted with the correct query string.
		if !strings.Contains(r.URL.RequestURI(), "/api.php?summary") {
			t.Error("@TestUpdateWithoutAPIKey: api.Summary.Update() did not send the expected query string when calling with an empty API key: /api.php?summary")
		}
	}))
	defer mockServer.Close()
	url := mockServer.URL + "/api.php"
	summary.Update(url, "", nil)
}
