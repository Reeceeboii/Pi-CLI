package data

import (
	"encoding/json"
	"github.com/Reeceeboii/Pi-CLI/pkg/constants"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
)

// Store PiCLI settings
var PICLISettings = NewSettings()

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
}

// Generate the location of the config file (or at least where it should be)
var configFileLocation = GetConfigFileLocation()

// Checks for the existence of a config file
func ConfigFileExists() bool {
	_, err := os.Stat(configFileLocation)
	return !os.IsNotExist(err)
}

// Return a new Settings instance
func NewSettings() *Settings {
	return &Settings{
		PiHoleAddress: "",
		PiHolePort:    constants.DefaultPort,
		RefreshS:      constants.DefaultRefreshS,
		APIKey:        "",
	}
}

// Attempts to create a settings instance from a config file
func (settings *Settings) LoadFromFile() {
	if byteArr, err := ioutil.ReadFile(configFileLocation); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(byteArr, settings); err != nil {
			log.Fatal(err)
		}
	}
}

// Saves the current settings to a config file
func (settings *Settings) SaveToFile() error {
	byteArr, err := json.MarshalIndent(settings, "", "\t")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(configFileLocation, byteArr, 0644); err != nil {
		return err
	}
	return nil
}

// Is API key stored in the config file? If not, off to the system keyring you go!
func (settings *Settings) APIKeyIsInFile() bool {
	return settings.APIKey != ""
}

// Delete the config file if it exists
func DeleteConfigFile() bool {
	// first, check if the file actually exists
	if !ConfigFileExists() {
		return false
	}
	if err := os.Remove(configFileLocation); err != nil {
		return false
	}
	return true
}

// Return the path to the config file
func GetConfigFileLocation() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	/*
		Return user's home directory plus the config file name. If on Windows, make sure path is returned
		with backslashes as the directory separators rather than forward slashes
	*/
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path.Join(usr.HomeDir, constants.ConfigFileName), "/", "\\")
	}
	return path.Join(usr.HomeDir, constants.ConfigFileName)
}
