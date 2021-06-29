package update

import (
	"time"
)

// Represents a single binary as part of a larger release
type Asset struct {
	// Name of the asset
	Name string `json:"name"`
	// Size of the asset in bytes
	Size int64 `json:"size_bytes"`
	// URL location of the asset
	DownloadURL string `json:"download_url"`
}

// Represents a release that contains multiple Asset instances
type Release struct {
	// Unix timestamp of when this release was received from GitHub
	TimeChecked string `json:"time_checked"`
	// SemVer Release tag of the latest GitHub release, I.e 1.1.5
	RemoteTag string `json:"remote_tag"`
	// Slice of Asset instances belonging to this Release
	Assets []Asset `json:"assets"`
}

// NewRelease creates a Release instance with some defaults
func NewRelease() Release {
	return Release{
		TimeChecked: time.Now().Format(time.RFC3339),
		RemoteTag:   "",
		Assets:      []Asset{},
	}
}
