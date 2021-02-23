package main

import (
	"github.com/buger/jsonparser"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"log"
	"net/http"
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
func storeAPIKeyInKeyring(key string) {
	if err := keyring.Set(service, usr, key); err != nil {
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

// does an key allow authentication? I.e., is is valid?
func validateAPIKey(key string) bool {
	/*
		To test the validity of the API key, we can attempt to enable the Pi-Hole.

		The response for a correct key:
				{
					"status": "enabled"
				}

		And the response for an incorrect key:
				[]

		Therefore we can simply perform a lookup for that "status" key. If it's there, the key is valid.

	*/

	url := piCLIData.FormattedAPIAddress + "?enable" + "&auth=" + key

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	parsedBody, _ := ioutil.ReadAll(res.Body)

	if _, err := jsonparser.GetString(parsedBody, "status"); err != nil {
		return false
	}
	return true
}
