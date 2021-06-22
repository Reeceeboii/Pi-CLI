package api

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/constants"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// instance of AllQueries used at runtime
var LiveAllQueries = NewAllQueries()

// Holds information about a single query logged by Pi-Hole
type Query struct {
	// UNIX timestamp of when the query was logged
	UnixTime string
	// The type of query
	QueryType string
	// The domain the query was sent to
	Domain string
	// The client that sent the query
	OriginClient string
	// Where the query was forwarded to
	ForwardedTo string
}

// Holds a slice of query structs
type AllQueries struct {
	// Slice of Query structs
	Queries []Query
	// The amount of queries being stored in the log
	AmountOfQueriesInLog int
	// The queries stored in a format able to be displayed as a table
	Table []string
}

// Make a new AllQueries instance
func NewAllQueries() *AllQueries {
	return &AllQueries{
		Queries:              make([]Query, constants.DefaultAmountOfQueries),
		AmountOfQueriesInLog: constants.DefaultAmountOfQueries,
		Table:                []string{},
	}
}

// Updates the all queries list with up to date information from the Pi-Hole
func (allQueries *AllQueries) Update(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	queryAmount := strconv.Itoa(allQueries.AmountOfQueriesInLog)
	url := data.LivePiCLIData.FormattedAPIAddress +
		"?getAllQueries=" +
		queryAmount +
		"&auth=" +
		data.LivePiCLIData.APIKey

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

	/*
		For every index in the parsed body's data array, pull out the required fields.
		I tried to use ArrayEach here but couldn't seem to get it to work the way I wanted.
		There has to be a nicer way to do this. This approach is absolute garbage.
	*/
	for iter := 0; iter < allQueries.AmountOfQueriesInLog; iter++ {
		queryArray, _, _, _ := jsonparser.Get(parsedBody, constants.AllQueryDataKey, fmt.Sprintf("[%d]", iter))
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

// Convert slice of queries to a formatted multidimensional slice
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
