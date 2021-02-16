package main

import (
	"log"
	"os"
	"time"
)

// constant values required by Pi-CLI
const (
	// port that the Pi-Hole API is defaulted to
	defaultPort = 80
	// the default refresh rate of the data in seconds
	defaultRefreshS = 1
	// the name of the configuration file
	configFileName = "config.json"
	// the starting setting for the number of queries that are included in the live log
	defaultAmountOfQueries = 10
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
	QueriesToday:        "",
	BlockedToday:        "",
	PercentBlockedToday: "",
	DomainsOnBlocklist:  "",
	Status:              "",
	PrivacyLevel:        "",
	TotalClientsSeen:    "",
}

var topItems = TopItems{
	TopQueries:       map[string]int{},
	TopAds:           map[string]int{},
	PrettyTopQueries: []string{},
	PrettyTopAds:     []string{},
}

var allQueries = AllQueries{
	Queries:              make([]query, defaultAmountOfQueries),
	AmountOfQueriesInLog: defaultAmountOfQueries,
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
