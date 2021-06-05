package network

import (
	"net/http"
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
