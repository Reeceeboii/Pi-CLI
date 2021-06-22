package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

/*
	This file stores commands that can be ran 'one-off', i.e. without needing to boot the live UI
*/

func RunSummaryCommand(c *cli.Context) error {
	InitialisePICLI()
	api.LiveSummary.Update(nil)
	fmt.Printf("Summary @ %s\n", time.Now().Format(time.Stamp))
	fmt.Println()
	fmt.Printf("Pi-Hole status: %s\n", strings.Title(api.LiveSummary.Status))
	fmt.Println()
	fmt.Printf("Queries /24hr: %s\n", api.LiveSummary.QueriesToday)
	fmt.Printf("Blocked /24hr: %s\n", api.LiveSummary.BlockedToday)
	fmt.Printf("Percent blocked: %s%s\n", api.LiveSummary.PercentBlockedToday, "%")
	fmt.Printf("Domains on blocklist: %s\n", api.LiveSummary.DomainsOnBlocklist)
	fmt.Printf("Privacy level: %s - %s\n",
		api.LiveSummary.PrivacyLevel,
		api.LiveSummary.PrivacyLevelNumberMapping[api.LiveSummary.PrivacyLevel],
	)
	fmt.Printf("Total clients seen: %s\n", api.LiveSummary.TotalClientsSeen)
	fmt.Println()
	return nil
}
