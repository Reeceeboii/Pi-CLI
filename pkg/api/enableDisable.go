package api

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"log"
	"net/http"
)

// Enable the Pi-Hole
func EnablePiHole() {
	url := data.LivePiCLIData.FormattedAPIAddress + "?enable" + "&auth=" + data.LivePiCLIData.APIKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = network.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

// Disable the Pi-Hole
func DisablePiHole(timeout bool, time int64) {
	disable := "?disable"
	if timeout {
		disable += fmt.Sprintf("=%d", time)
	}
	url := data.LivePiCLIData.FormattedAPIAddress + disable + "&auth=" + data.LivePiCLIData.APIKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = network.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
