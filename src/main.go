package main

import (
	"log"
	"os"
	"time"
)

// stores the data needed by Pi-CLI during runtime
type PiCLIData struct {
	Settings            *Settings
	FormattedAPIAddress string
	APIKey              string
	LastUpdated         time.Time
}

var piCLIData = PiCLIData{}

var summary = Summary{
	QueriesToday:        "0",
	BlockedToday:        "0",
	PercentBlockedToday: "0.0",
	DomainsOnBlocklist:  "0",
}

var topItems = TopItems{
	TopQueries:       map[string]int{},
	TopAds:           map[string]int{},
	PrettyTopQueries: []string{},
	PrettyTopAds:     []string{},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
