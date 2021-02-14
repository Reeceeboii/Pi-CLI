package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

// BasicData holds 'basic' data: things that do not require authentication to retrieve
type BasicData struct {
	QueriesToday        int64
	BlockedToday        int64
	PercentBlockedToday float32
	DomainsOnBlocklist  int64
}

func (basicData *BasicData) update() {
	req, err := http.NewRequest("GET", piCLIData.FormattedAPIAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	parsedBody, _ := ioutil.ReadAll(res.Body)
	queriesToday, _ := jsonparser.GetInt(parsedBody, "dns_queries_today")
	blockedToday, _ := jsonparser.GetInt(parsedBody, "ads_blocked_today")
	percentageToday, _ := jsonparser.GetFloat(parsedBody, "ads_percentage_today")
	domainsOnBlocklist, _ := jsonparser.GetInt(parsedBody, "domains_being_blocked")
	basicData.QueriesToday = queriesToday
	basicData.BlockedToday = blockedToday
	basicData.PercentBlockedToday = float32(percentageToday)
	basicData.DomainsOnBlocklist = domainsOnBlocklist
}

func generateAPIAddress() string {
	return fmt.Sprintf("http://%s:%d/admin/api.php", settings.PiHoleAddress, settings.PiHolePort)
}
