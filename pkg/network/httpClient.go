package network

import (
	"net/http"
	"os"
	"time"
)

// Construct a http.Client with a 3 second timeout for use in API requests
var HttpClient = NewHTTPClient(time.Second * 3)

// Create a new http.Client with a given timeout duration
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

/*
	Generate a new http.Request (HTTP GET) with an accept gzip header, and,
	from the environment if it has been loaded:
		- A GitHub API key used to increase the number of hourly requests Pi-CLI
		is allowed to make. This is useful in development.
*/
func NewRequestWithGzipHeaders(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip")

	gitHubToken := os.Getenv("GITHUB_API_TOKEN")
	// add a GitHub API token to the request if one is provided via the environment
	if gitHubToken != "" {
		req.Header.Add("Authorization", "token "+gitHubToken)
	}

	return req, nil
}
