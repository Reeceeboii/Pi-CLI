package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

/*
	This file stores commands that can be ran 'one-off', i.e. without needing to boot the live UI
*/

/*
	Extracts a quick summary of the previous 24/hr of data from the Pi-Hole.
*/
func RunSummaryCommand(*cli.Context) error {
	InitialisePICLI()
	api.LiveSummary.Update(nil)
	fmt.Printf("Summary @ %s\n", time.Now().Format(time.Stamp))
	fmt.Println()

	if api.LiveSummary.Status == "enabled" {
		fmt.Printf("Pi-Hole status: %s\n", color.GreenString(strings.Title(api.LiveSummary.Status)))
	} else {
		fmt.Printf("Pi-Hole status: %s\n", color.RedString(strings.Title(api.LiveSummary.Status)))
	}

	fmt.Println()
	fmt.Printf("Queries /24hr: %s\n", api.LiveSummary.QueriesToday)
	fmt.Printf("Blocked /24hr: %s\n", api.LiveSummary.BlockedToday)
	fmt.Printf("Percent blocked: %s%%\n", api.LiveSummary.PercentBlockedToday)
	fmt.Printf("Domains on blocklist: %s\n", api.LiveSummary.DomainsOnBlocklist)
	fmt.Printf("Privacy level: %s - %s\n",
		api.LiveSummary.PrivacyLevel,
		api.LiveSummary.PrivacyLevelNumberMapping[api.LiveSummary.PrivacyLevel],
	)
	fmt.Printf("Total clients seen: %s\n", api.LiveSummary.TotalClientsSeen)
	fmt.Println()
	return nil
}

/*
	Extract the current top 10 permitted domains that have been forwarded to the upstream DNS resolver
*/
func RunTopTenForwardedCommand(*cli.Context) error {
	InitialisePICLI()

	api.LiveTopItems.Update(nil)
	fmt.Printf("Top queries as of @ %s\n\n", time.Now().Format(time.Stamp))
	for _, q := range api.LiveTopItems.PrettyTopQueries {
		fmt.Println(q)
	}

	return nil
}

/*
	Extract the current top 10 blocked domains that the FTL has filtered out and not forwarded
	to the upstream DNS resolver
*/
func RunTopTenBlockedCommand(*cli.Context) error {
	InitialisePICLI()

	api.LiveTopItems.Update(nil)
	fmt.Printf("Top blocked domains as of @ %s\n\n", time.Now().Format(time.Stamp))
	for _, q := range api.LiveTopItems.PrettyTopAds {
		fmt.Println(q)
	}

	return nil
}

func RunLatestQueriesCommand(c *cli.Context) error {
	queryAmount := c.Int("limit")
	if queryAmount == 0 {
		queryAmount = 10
	}

	if queryAmount < 1 {
		fmt.Println("Please enter a number of queries >= 1")
		return nil
	}

	InitialisePICLI()

	api.LiveAllQueries.AmountOfQueriesInLog = queryAmount
	api.LiveAllQueries.Queries = make([]api.Query, api.LiveAllQueries.AmountOfQueriesInLog)
	api.LiveAllQueries.Update(nil)

	for _, query := range api.LiveAllQueries.Table {
		fmt.Println(query)
	}

	return nil
}
