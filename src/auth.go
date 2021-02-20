package main

import (
	"github.com/zalando/go-keyring"
	"log"
)

const (
	service = "PiCLI"
	user    = service
)

// retrieve the API key from the system keyring
func retrieveAPIKey() string {
	APIKey, err := keyring.Get(service, user)
	if err != nil {
		log.Fatal(err)
	}
	return APIKey
}

// store the API key in the system keyring
func storeAPIKey(key *string) {
	if err := keyring.Set(service, user, *key); err != nil {
		log.Fatal(err)
	}
}

// delete the stored API key if it exists
func deleteAPIKey() bool {
	if err := keyring.Delete(service, user); err != nil {
		return false
	}
	return true
}

// does the API key exist?
func APIKeyExists() bool {
	if _, err := keyring.Get(service, user); err != nil {
		return false
	}
	return true
}
