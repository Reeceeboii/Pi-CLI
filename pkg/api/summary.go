package api

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/constants"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var LiveSummary = NewSummary()

// Summary holds things that do not require authentication to retrieve
type Summary struct {
	// Total number of queries logged today
	QueriesToday string
	// Total number of queries blocked today
	BlockedToday string
	// Percentage of today's queries that have been blocked
	PercentBlockedToday string
	// How large is Pi-Hole's active blocklist?
	DomainsOnBlocklist string
	// Enabled vs. disabled
	Status string
	// Pi-Hole's current data privacy level
	PrivacyLevel string
	// Mapping between privacy level numbers and their meanings
	PrivacyLevelNumberMapping map[string]string
	// The total number of clients that the Pi-Hole has seen
	TotalClientsSeen string
}

func NewSummary() *Summary {
	return &Summary{
		QueriesToday:        "",
		BlockedToday:        "",
		PercentBlockedToday: "",
		DomainsOnBlocklist:  "",
		Status:              "",
		PrivacyLevel:        "",
		PrivacyLevelNumberMapping: map[string]string{
			"0": "Show Everything",
			"1": "Hide Domains",
			"2": "Hide Domains and Clients",
			"3": "Anonymous",
		},
		TotalClientsSeen: "",
	}
}

// Updates a Summary struct with up to date information
func (summary *Summary) Update(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	// create the URL for the summary data and send a request to it
	url := data.LivePiCLIData.FormattedAPIAddress + "?summary"
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
	// yoink out all the data from the response
	// pack it into the struct
	summary.QueriesToday, _ = jsonparser.GetString(parsedBody, constants.DNSQueriesTodayKey)
	summary.BlockedToday, _ = jsonparser.GetString(parsedBody, constants.AdsBlockedTodayKey)
	summary.PercentBlockedToday, _ = jsonparser.GetString(parsedBody, constants.PercentBlockedTodayKey)
	summary.DomainsOnBlocklist, _ = jsonparser.GetString(parsedBody, constants.DomainsOnBlockListKey)
	summary.Status, _ = jsonparser.GetString(parsedBody, constants.StatusKey)
	summary.PrivacyLevel, _ = jsonparser.GetString(parsedBody, constants.PrivacyLevelKey)
	summary.TotalClientsSeen, _ = jsonparser.GetString(parsedBody, constants.TotalClientsSeenKey)
}