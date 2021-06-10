package auth

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/constants"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"log"
	"net/http"
)

// Retrieve the API key from the system keyring
func RetrieveAPIKeyFromKeyring() string {
	APIKey, err := keyring.Get(constants.KeyringService, constants.KeyringUsr)
	if err != nil {
		log.Fatal(err)
	}
	return APIKey
}

/*
	Store the API key in the system keyring. Returns an error if this action failed.
*/
func StoreAPIKeyInKeyring(key string) error {
	if err := keyring.Set(constants.KeyringService, constants.KeyringUsr, key); err != nil {
		return err
	}
	return nil
}

// Delete the stored API key if it exists
func DeleteAPIKeyFromKeyring() bool {
	if err := keyring.Delete(constants.KeyringService, constants.KeyringUsr); err != nil {
		return false
	}
	return true
}

// Is there an entry for the API key in the system keyring?
func APIKeyIsInKeyring() bool {
	if _, err := keyring.Get(constants.KeyringService, constants.KeyringUsr); err != nil {
		return false
	}
	return true
}

// Does an key allow authentication? I.e., is is valid?
func ValidateAPIKey(key string) bool {
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

	url := data.LivePiCLIData.FormattedAPIAddress + "?enable" + "&auth=" + key

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := network.HttpClient.Do(req)
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
