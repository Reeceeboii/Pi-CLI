package auth

import (
	"testing"
)

const (
	// Sample API key for test case usage.
	testKey = "c808f484a4e88cc32a9a8bfcce19169c77bcd9c5eec18d859e1bb4b318bf42bf"
)

// Calling init() in order to overwrite global variables for test purposes.
func init() {
	KeyringService = "test-service" // Overwrite KeyringService for test cases
	KeyringUsr = "test-key"         // Overwrite KeyringUsr for test cases
}

/*
  NOTE:
  Each test case is self-contained, meaning a key is stored at the beginning of each case and deleted before it ends.
  We do this because we cannot rely on Go to run its tests sequentially every time.
*/

// Tests for auth.APIKeyIsInKeyring()
func TestAPIKeyIsInKeyring(t *testing.T) {
	// Ensuring StoreAPIKeyInKeyring() can successfully store a key in the keyring.
	err := StoreAPIKeyInKeyring(testKey)
	if err != nil {
		t.Errorf("@TestAPIKeyIsInKeyring: auth.StoreAPIKeyInKeyring() failed to store API key: %s", err)
	}

	// Ensuring APIKeyIsInKeyring() can successfully find the stored key.
	if !APIKeyIsInKeyring() {
		t.Error("@TestAPIKeyIsInKeyring: auth.APIKeyIsInKeyring() failed to find key in keyring.")
	}

	// Ensuring DeleteAPIKeyFromKeyring() is able to successfully find and delete key
	if !DeleteAPIKeyFromKeyring() {
		t.Error("@TestRetrieveAPIKeyFromKeyring: auth.DeleteAPIKeyFromKeyring() did not find/delete key in keyring.")
	}

	// Ensuring APIKeyIsInKeyring() cannot find a key that should not exist.
	if APIKeyIsInKeyring() {
		t.Error("@TestAPIKeyIsInKeyring: auth.APIKeyIsInKeyring() found key in keyring after it should have been deleted.")
	}
}

// Tests for auth.RetrieveAPIKeyFromKeyring()
func TestRetrieveAPIKeyFromKeyring(t *testing.T) {
	// Ensuring StoreAPIKeyInKeyring() can successfully store a key in the keyring.
	err := StoreAPIKeyInKeyring(testKey)
	if err != nil {
		t.Errorf("@TestRetrieveAPIKeyFromKeyring: auth.StoreAPIKeyInKeyring() failed to store API key: %s", err)
	}

	// Ensuring RetrieveAPIKeyFromKeyring() can successfully find the right key in keyring.
	key := RetrieveAPIKeyFromKeyring()
	if key != testKey {
		t.Error("@TestRetrieveAPIKeyFromKeyring: auth.RetrieveAPIKeyFromKeyring() did not match provided test key.")
	}

	// Ensuring DeleteAPIKeyFromKeyring() is able to successfully find and delete key
	if !DeleteAPIKeyFromKeyring() {
		t.Error("@TestRetrieveAPIKeyFromKeyring: auth.DeleteAPIKeyFromKeyring() did not find/delete key in keyring.")
	}
}

// Tests for auth.DeleteAPIKeyFromKeyring()
func TestDeleteAPIKeyFromKeyring(t *testing.T) {
	// Ensuring StoreAPIKeyInKeyring() can successfully store a key in the keyring.
	err := StoreAPIKeyInKeyring(testKey)
	if err != nil {
		t.Errorf("@TestDeleteAPIKeyFromKeyring: auth.StoreAPIKeyInKeyring() failed to store API key: %s", err)
	}

	// Ensuring DeleteAPIKeyFromKeyring() is able to successfully find and delete key
	if !DeleteAPIKeyFromKeyring() {
		t.Error("@TestDeleteAPIKeyFromKeyring: auth.DeleteAPIKeyFromKeyring() did not find/delete key in keyring.")
	}

	// Ensuring DeleteAPIKeyFromKeyring() does not find or delete a key as expected when the key does not exist.
	if DeleteAPIKeyFromKeyring() {
		t.Error("@TestDeleteAPIKeyFromKeyring: auth.DeleteAPIKeyFromKeyring() found/deleted key from keyring when one should not exist, it should have been deleted in the previous assertion.")
	}
}
