package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

// keys that can be used to index JSON responses from the Pi-Hole's API
const (
	// Summary
	DNSQueriesTodayKey     = "dns_queries_today"
	AdsBlockedTodayKey     = "ads_blocked_today"
	PercentBlockedTodayKey = "ads_percentage_today"
	DomainsOnBlockListKey  = "domains_being_blocked"
	// TopItems
	TopQueriesTodayKey = "top_queries"
	TopAdsTodayKey     = "top_ads"
)

// Summary holds things that do not require authentication to retrieve
type Summary struct {
	QueriesToday        string
	BlockedToday        string
	PercentBlockedToday string
	DomainsOnBlocklist  string
}

// TopItems stores top permitted domains and top blocked domains (requires authentication to retrieve)
type TopItems struct {
	TopQueries       map[string]int
	TopAds           map[string]int
	PrettyTopQueries []string
	PrettyTopAds     []string
}

// updates a Summary struct with up to date information
func (summary *Summary) update() {
	url := piCLIData.FormattedAPIAddress + "?summary"
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
	// yoink out all the data from the response
	queriesToday, _ := jsonparser.GetString(parsedBody, DNSQueriesTodayKey)
	blockedToday, _ := jsonparser.GetString(parsedBody, AdsBlockedTodayKey)
	percentageToday, _ := jsonparser.GetString(parsedBody, PercentBlockedTodayKey)
	domainsOnBlocklist, _ := jsonparser.GetString(parsedBody, DomainsOnBlockListKey)
	// pack it into the struct
	summary.QueriesToday = queriesToday
	summary.BlockedToday = blockedToday
	summary.PercentBlockedToday = percentageToday
	summary.DomainsOnBlocklist = domainsOnBlocklist
}

// updates a TopItems struct with up to date information
func (topItems *TopItems) update() {
	url := piCLIData.FormattedAPIAddress + "?topItems" + "&auth=" + piCLIData.APIKey
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

	// parse the top queries response
	_ = jsonparser.ObjectEach(parsedBody, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		topItems.TopQueries[string(key)], _ = strconv.Atoi(string(value))
		return nil
	}, TopQueriesTodayKey)

	// and the same for the top ad networks
	_ = jsonparser.ObjectEach(parsedBody, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		topItems.TopAds[string(key)], _ = strconv.Atoi(string(value))
		return nil
	}, TopAdsTodayKey)
}

// convert maps of domain:hits to a nice lists that can be displayed
func (topItems *TopItems) prettyConvert() {
	topItems.PrettyTopQueries = []string{}
	topItems.PrettyTopAds = []string{}

	for key, value := range topItems.TopQueries {
		topItems.PrettyTopQueries = append(topItems.PrettyTopQueries, fmt.Sprintf("[%d] %s", value, key))
	}

	for key, value := range topItems.TopAds {
		topItems.PrettyTopAds = append(topItems.PrettyTopAds, fmt.Sprintf("[%d] %s", value, key))
	}
}

// plug the Pi-Hole address and port together to get a full URL
func generateAPIAddress() string {
	return fmt.Sprintf("http://%s:%d/admin/api.php", settings.PiHoleAddress, settings.PiHolePort)
}
