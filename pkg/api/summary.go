package api

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var LiveSummary = NewSummary()

// Keys that can be used to index JSON responses from the Pi-Hole's API
const (
	DNSQueriesTodayKey     = "dns_queries_today"
	AdsBlockedTodayKey     = "ads_blocked_today"
	PercentBlockedTodayKey = "ads_percentage_today"
	DomainsOnBlockListKey  = "domains_being_blocked"
	StatusKey              = "status"
	PrivacyLevelKey        = "privacy_level"
	TotalClientsSeenKey    = "clients_ever_seen"
)

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

/*
Returns a new Summary instance with default values for all fields
*/
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
func (summary *Summary) Update(url string, key string, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	// create the URL for the summary data and send a request to it
	url += "?summary"
	if len(key) > 0 {
		url += "&auth=" + key
	}

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
	summary.QueriesToday, _ = jsonparser.GetString(parsedBody, DNSQueriesTodayKey)
	summary.BlockedToday, _ = jsonparser.GetString(parsedBody, AdsBlockedTodayKey)
	summary.PercentBlockedToday, _ = jsonparser.GetString(parsedBody, PercentBlockedTodayKey)
	summary.DomainsOnBlocklist, _ = jsonparser.GetString(parsedBody, DomainsOnBlockListKey)
	summary.Status, _ = jsonparser.GetString(parsedBody, StatusKey)
	summary.PrivacyLevel, _ = jsonparser.GetString(parsedBody, PrivacyLevelKey)
	summary.TotalClientsSeen, _ = jsonparser.GetString(parsedBody, TotalClientsSeenKey)
}
