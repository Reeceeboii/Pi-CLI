package network

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
)

// Plug the Pi-Hole address and port together to get a full URL
func GenerateAPIAddress(address string, port int) string {
	return fmt.Sprintf("http://%s:%d/admin/api.php", address, port)
}

// IsAlive will, given an IP and port, return true or false denoting if the address is alive
func IsAlive(address string) bool {
	color.Yellow("Validating " + address)
	req, err := http.NewRequest("GET", address, nil)

	if err != nil {
		color.Red("Failed to generate HTTP GET in IsAlive()")
		log.Fatal(err)
	}

	_, err = HttpClient.Do(req)
	if err != nil {
		color.Red("Address not reachable!")
		return false
	}
	return true
}

/*
	Do the provided address & port actually point to a live Pi-Hole?
	Issue #16 @https://github.com/Reeceeboii/Pi-CLI/issues/16
*/
func ValidatePiHoleDetails(res *http.Response) bool {
	return res.StatusCode == http.StatusOK
}
