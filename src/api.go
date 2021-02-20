package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 3,
}

// keys that can be used to index JSON responses from the Pi-Hole's API
const (
	// Summary
	DNSQueriesTodayKey     = "dns_queries_today"
	AdsBlockedTodayKey     = "ads_blocked_today"
	PercentBlockedTodayKey = "ads_percentage_today"
	DomainsOnBlockListKey  = "domains_being_blocked"
	StatusKey              = "status"
	PrivacyLevelKey        = "privacy_level"
	TotalClientsSeenKey    = "clients_ever_seen"
	// TopItems
	TopQueriesTodayKey = "top_queries"
	TopAdsTodayKey     = "top_ads"
	// GetAllQueries
	AllQueryDataKey = "data"
)

// Summary holds things that do not require authentication to retrieve
type Summary struct {
	QueriesToday              string
	BlockedToday              string
	PercentBlockedToday       string
	DomainsOnBlocklist        string
	Status                    string
	PrivacyLevel              string
	PrivacyLevelNumberMapping map[string]string
	TotalClientsSeen          string
}

// TopItems stores top permitted domains and top blocked domains (requires authentication to retrieve)
type TopItems struct {
	TopQueries       map[string]int
	TopAds           map[string]int
	PrettyTopQueries []string
	PrettyTopAds     []string
}

// holds information about a single query logged by Pi-Hole
type Query struct {
	UnixTime     string
	QueryType    string
	Domain       string
	OriginClient string
	ForwardedTo  string
}

// holds a slice of query structs
type AllQueries struct {
	Queries              []Query
	AmountOfQueriesInLog int
	Table                []string
}

type domainOccurrencePair struct {
	domain     string
	occurrence int
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
	// pack it into the struct
	summary.QueriesToday, _ = jsonparser.GetString(parsedBody, DNSQueriesTodayKey)
	summary.BlockedToday, _ = jsonparser.GetString(parsedBody, AdsBlockedTodayKey)
	summary.PercentBlockedToday, _ = jsonparser.GetString(parsedBody, PercentBlockedTodayKey)
	summary.DomainsOnBlocklist, _ = jsonparser.GetString(parsedBody, DomainsOnBlockListKey)
	summary.Status, _ = jsonparser.GetString(parsedBody, StatusKey)
	summary.PrivacyLevel, _ = jsonparser.GetString(parsedBody, PrivacyLevelKey)
	summary.TotalClientsSeen, _ = jsonparser.GetString(parsedBody, TotalClientsSeenKey)
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

	topItems.prettyConvert()
}

// convert maps of domain:hits to nice lists that can be displayed
func (topItems *TopItems) prettyConvert() {
	var sortedTopQueries []domainOccurrencePair
	var sortedTopAds []domainOccurrencePair
	topItems.PrettyTopQueries = []string{}
	topItems.PrettyTopAds = []string{}

	for key, value := range topItems.TopQueries {
		sortedTopQueries = append(sortedTopQueries, domainOccurrencePair{
			domain:     key,
			occurrence: value,
		})
	}

	for key, value := range topItems.TopAds {
		sortedTopAds = append(sortedTopAds, domainOccurrencePair{
			domain:     key,
			occurrence: value,
		})
	}

	// sort ads and domains by occurrence
	sort.SliceStable(sortedTopQueries[:], func(i, j int) bool {
		return sortedTopQueries[i].occurrence > sortedTopQueries[j].occurrence
	})
	sort.SliceStable(sortedTopAds[:], func(i, j int) bool {
		return sortedTopAds[i].occurrence > sortedTopAds[j].occurrence
	})

	for _, domain := range sortedTopQueries {
		listEntry := fmt.Sprintf("%d hits | %s", domain.occurrence, domain.domain)
		topItems.PrettyTopQueries = append(topItems.PrettyTopQueries, listEntry)
	}
	for _, domain := range sortedTopAds {
		listEntry := fmt.Sprintf("%d hits | %s", domain.occurrence, domain.domain)
		topItems.PrettyTopAds = append(topItems.PrettyTopAds, listEntry)
	}
}

// convert slice of queries to a formatted multidimensional slice
func (allQueries *AllQueries) convertToTable() {
	table := make([]string, allQueries.AmountOfQueriesInLog)

	for i, q := range allQueries.Queries {
		iTime, _ := strconv.ParseInt(q.UnixTime, 10, 64)
		parsedTime := time.Unix(iTime, 0)
		entry := fmt.Sprintf("%d [%s] Query type %s from %s to %s forwarded to %s",
			(allQueries.AmountOfQueriesInLog)-i,
			parsedTime.Format("15:04:05"),
			q.QueryType,
			q.OriginClient,
			q.Domain,
			q.ForwardedTo,
		)
		table[(allQueries.AmountOfQueriesInLog-1)-i] = entry
	}
	allQueries.Table = table
}

func (allQueries *AllQueries) update() {
	url := piCLIData.FormattedAPIAddress + "?getAllQueries=" + strconv.Itoa(allQueries.AmountOfQueriesInLog) + "&auth=" + piCLIData.APIKey
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
	// for every index in the parsed body's data array, pull out the required fields
	// I tried to use ArrayEach here but couldn't seem to get it to work the way I wanted
	// there has to be a nicer way to do this. This approach is absolute garbage
	for iter := 0; iter < allQueries.AmountOfQueriesInLog; iter++ {
		queryArray, _, _, _ := jsonparser.Get(parsedBody, AllQueryDataKey, fmt.Sprintf("[%d]", iter))
		unixTime, _ := jsonparser.GetString(queryArray, "[0]")
		queryType, _ := jsonparser.GetString(queryArray, "[1]")
		domain, _ := jsonparser.GetString(queryArray, "[2]")
		originClient, _ := jsonparser.GetString(queryArray, "[3]")
		forwardedTo, _ := jsonparser.GetString(queryArray, "[10]")
		allQueries.Queries[iter] = Query{
			UnixTime:     unixTime,
			QueryType:    queryType,
			Domain:       domain,
			OriginClient: originClient,
			ForwardedTo:  forwardedTo,
		}
	}

	allQueries.convertToTable()
}

// plug the Pi-Hole address and port together to get a full URL
func generateAPIAddress() string {
	return fmt.Sprintf("http://%s:%d/admin/api.php", settings.PiHoleAddress, settings.PiHolePort)
}

// do the provided address & port actually point to a Pi-Hole?
func validatePiHoleDetails(res *http.Response) bool {
	return res.StatusCode == 200 && res.Header.Get("X-Pi-hole") != ""
}
