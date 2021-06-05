package network

import (
	"fmt"
	"net/http"
)

// Plug the Pi-Hole address and port together to get a full URL
func GenerateAPIAddress(address string, port int) string {
	return fmt.Sprintf("http://%s:%d/admin/api.php", address, port)
}

/*
	Do the provided address & port actually point to a live Pi-Hole?
	Issue #16 @https://github.com/Reeceeboii/Pi-CLI/issues/16
*/
func ValidatePiHoleDetails(res *http.Response) bool {
	return res.StatusCode == http.StatusOK
}
