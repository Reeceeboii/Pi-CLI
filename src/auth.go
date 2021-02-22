package main

import (
	"github.com/zalando/go-keyring"
	"log"
)

const (
	service = "PiCLI"
	usr     = "api-key"
)

// retrieve the API key from the system keyring
func retrieveAPIKeyFromKeyring() string {
	APIKey, err := keyring.Get(service, usr)
	if err != nil {
		log.Fatal(err)
	}
	return APIKey
}

// store the API key in the system keyring
func storeAPIKeyInKeyring(key *string) {
	if err := keyring.Set(service, usr, *key); err != nil {
		log.Fatal(err)
	}
}

// delete the stored API key if it exists
func deleteAPIKeyFromKeyring() bool {
	if err := keyring.Delete(service, usr); err != nil {
		return false
	}
	return true
}

// is there an entry for the API key in the system keyring?
func APIKeyIsInKeyring() bool {
	if _, err := keyring.Get(service, usr); err != nil {
		return false
	}
	return true
}
