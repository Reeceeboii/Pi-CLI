package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Settings contains the current configuration options being used by the program
type Settings struct {
	PiHoleAddress string `json:"pi_hole_address"`
	PiHolePort    int    `json:"pi_hole_port"`
	RefreshS      int    `json:"refresh_s"`
}

// checks for the existence of a config file
func configFileExists() bool {
	_, err := os.Stat(configFileName)
	return !os.IsNotExist(err)
}

// Attempts to create a settings instance from a config file
func (settings *Settings) loadFromFile() {
	if byteArr, err := ioutil.ReadFile(configFileName); err != nil {
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
	if err = ioutil.WriteFile(configFileName, byteArr, 0644); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Saved configuration to " + configFileName)
	}
}
