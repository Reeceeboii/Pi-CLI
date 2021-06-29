package data

import (
	"encoding/json"
	"github.com/Reeceeboii/Pi-CLI/pkg/logger"
	"github.com/Reeceeboii/Pi-CLI/pkg/update"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"sync"
)

// Store PiCLI settings
var PICLISettings = NewSettings()

// Constant values required by Pi-CLI
const (
	// Port that the Pi-Hole API is defaulted to
	DefaultPort = 80
	// The default refresh rate of the data in seconds
	DefaultRefreshS = 1
	// The name of the configuration file
	ConfigFileName = ".piclirc"
	// The default setting for automatic update checks
	DefaultAutoCheckForUpdates = true
)

// Settings contains the current configuration options being used by Pi-CLI
type Settings struct {
	// The Pi-Hole's address
	PiHoleAddress string `json:"pi_hole_address"`
	// The port the Pi-Hole exposes that can be used for HTTP/S traffic
	PiHolePort int `json:"pi_hole_port"`
	// The number of seconds to wait between each data refresh
	RefreshS int `json:"refresh_s"`
	// API key used to authenticate with the Pi-Hole instance
	APIKey string `json:"api_key"`
	// Has the user chosen to automatically check for updates
	AutoCheckForUpdates bool `json:"auto_check_for_updates"`
	// Caches a response from the GitHub API release endpoint for Pi-CLI
	LatestRemoteRelease update.Release `json:"latest_remote_release"`
}

/*
	The location of the config file (or at least where it should be), along with a
	sync.Once instance. This allows the GetConfigFileLocation function to be called multiple times,
	while only the first call will do actual work. Subsequent calls will simply return a cached value
	as this location is not expected to change during  runtime.
*/
var (
	configFileLocation     string
	configFileLocationOnce sync.Once
)

// Return a new Settings instance
func NewSettings() *Settings {
	return &Settings{
		PiHoleAddress:       "",
		PiHolePort:          DefaultPort,
		RefreshS:            DefaultRefreshS,
		APIKey:              "",
		AutoCheckForUpdates: DefaultAutoCheckForUpdates,
		LatestRemoteRelease: update.Release{},
	}
}

// Checks for the existence of a config file
func ConfigFileExists(configFileLocation string) bool {
	_, err := os.Stat(configFileLocation)
	return !os.IsNotExist(err)
}

// Attempts to create a settings instance from a config file
func (settings *Settings) LoadFromFile(configFileLocation string) {
	if byteArr, err := ioutil.ReadFile(configFileLocation); err != nil {
		logger.LivePiCLILogger.LogError(err)
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(byteArr, settings); err != nil {
			logger.LivePiCLILogger.LogError(err)
			log.Fatal(err)
		}
	}
}

// SaveToFile saves the current settings to a config file
func (settings *Settings) SaveToFile() error {
	byteArr, err := json.MarshalIndent(settings, "", "\t")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(GetConfigFileLocation(), byteArr, 0644); err != nil {
		return err
	}
	return nil
}

// Is API key stored in the config file? If not, off to the system keyring you go!
func (settings *Settings) APIKeyIsInFile() bool {
	return settings.APIKey != ""
}

// DeleteConfigFile deletes the config file if it exists
func DeleteConfigFile(configFileLocation string) bool {
	// first, check if the file actually exists
	if !ConfigFileExists(configFileLocation) {
		return false
	}
	if err := os.Remove(configFileLocation); err != nil {
		return false
	}
	logger.LivePiCLILogger.LogInformation("Config file at " + GetConfigFileLocation() + " has been deleted!")
	return true
}

// GetConfigFileLocation gets the (expected) path to the config file
func GetConfigFileLocation() string {
	configFileLocationOnce.Do(func() {
		usr, err := user.Current()
		if err != nil {
			logger.LivePiCLILogger.LogError(err)
			log.Fatal(err)
		}

		/*
			Set user's home directory plus the config file name. If on Windows, make sure path is returned
			with backslashes as the directory separators rather than forward slashes
		*/
		if runtime.GOOS == "windows" {
			configFileLocation = strings.ReplaceAll(path.Join(usr.HomeDir, ConfigFileName), "/", "\\")
		} else {
			configFileLocation = path.Join(usr.HomeDir, ConfigFileName)
		}
	})
	return configFileLocation
}
