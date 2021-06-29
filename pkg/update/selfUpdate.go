package update

import (
	"compress/gzip"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// The URL that can be accessed to access Pi-CLI release information
	GitHubReleaseURL = "https://api.github.com/repos/Reeceeboii/Pi-CLI/releases"
)

// Set in the build via make - passed in via the main package
var Version = "Undefined"
var GitHash = "Undefined"

func GetLatestGitHubRelease(client *http.Client) Release {
	req, err := network.NewRequestWithGzipHeaders(GitHubReleaseURL)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var parsedBody []byte

	/*
		Not all of the GitHub API endpoints support gzip content-encoding. However, if it
		is available we want to accept it to speed up the request. We do however need to actually
		check if the request *was* gzip encoded or not. If we try decompress a non gzip encoded request
		we won't be having a good day
	*/
	if res.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gz.Close()
		parsedBody, _ = ioutil.ReadAll(gz)
	} else {
		parsedBody, _ = ioutil.ReadAll(res.Body)
	}

	release := NewRelease()

	// extract the tag of the latest release
	remoteTag, _ := jsonparser.GetString(parsedBody, "[0]", "tag_name")
	release.RemoteTag = remoteTag

	_, _ = jsonparser.ArrayEach(parsedBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, _ := jsonparser.GetString(value, "name")
		size, _ := jsonparser.GetInt(value, "size")
		downloadURL, _ := jsonparser.GetString(value, "browser_download_url")

		// use the extracted data to create a new asset, and append it to the assets slice
		release.Assets = append(release.Assets, Asset{
			Name:        name,
			Size:        size,
			DownloadURL: downloadURL,
		})
	}, "[0]", "assets")

	return release
}

// SetVersion sets the update.Version variable to the value passed into the binary via the linker
func SetVersion(version string) {
	Version = version
}

// SetGitHash sets the update.GitHash variable to the value passed into the binary via the linker
func SetGitHash(hash string) {
	GitHash = hash
}
