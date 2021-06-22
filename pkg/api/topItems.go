package api

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

// Instance of TopItems used at runtime
var LiveTopItems = NewTopItems()

// Keys that can be used to index JSON responses from the Pi-Hole's API
const (
	TopQueriesTodayKey = "top_queries"
	TopAdsTodayKey     = "top_ads"
)

// TopItems stores top permitted domains and top blocked domains (requires authentication to retrieve)
type TopItems struct {
	// Mapping of top DNS queried domains and their occurrences
	TopQueries map[string]int
	// Mapping of top blocked DNS domains (ads and/or tracking) and their occurrences
	TopAds map[string]int
	// Pretty list version of TopQueries
	PrettyTopQueries []string
	// Pretty list version of TopAds
	PrettyTopAds []string
}

// A single domain and the number of times it occurs
type domainOccurrencePair struct {
	// The domain
	domain string
	// The number of times it has occurred
	occurrence int
}

// Create a new TopItems instance
func NewTopItems() *TopItems {
	return &TopItems{
		TopQueries:       map[string]int{},
		TopAds:           map[string]int{},
		PrettyTopQueries: []string{},
		PrettyTopAds:     []string{},
	}
}

// Updates a TopItems struct with up to date information
func (topItems *TopItems) Update(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	url := data.LivePiCLIData.FormattedAPIAddress + "?topItems" + "&auth=" + data.LivePiCLIData.APIKey
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

// Convert maps of domain:hits to nice lists that can be displayed
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
