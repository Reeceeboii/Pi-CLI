package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

// Settings contains the current configuration options being used by the program
type Settings struct {
	PiHoleAddress string `json:"pi_hole_address"`
	PiHolePort    int    `json:"pi_hole_port"`
	RefreshS      int    `json:"refresh_s"`
}

// generate the location of the config file (or at least where it should be)
var configFileLocation = getConfigFileLocation()

// checks for the existence of a config file
func configFileExists() bool {
	_, err := os.Stat(configFileLocation)
	return !os.IsNotExist(err)
}

// Attempts to create a settings instance from a config file
func (settings *Settings) loadFromFile() {
	if byteArr, err := ioutil.ReadFile(configFileLocation); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(byteArr, settings); err != nil {
			log.Fatal(err)
		}
	}
}

// Saves the current settings to a config file
func (settings *Settings) saveToFile() {
	byteArr, err := json.MarshalIndent(settings, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile(configFileLocation, byteArr, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Saved configuration to " + configFileLocation)
	}
}

// delete the config file if it exists
func deleteConfigFile() bool {
	// first, check if the file actually exists
	if !configFileExists() {
		return false
	}
	if err := os.Remove(configFileLocation); err != nil {
		return false
	}
	return true
}

// return the path to the config file
func getConfigFileLocation() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	// return user's home directory plus the config file name
	return path.Join(usr.HomeDir, configFileName)
}
