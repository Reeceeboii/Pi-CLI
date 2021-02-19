package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"strings"
	"time"
)

// update data so it can be displayed
// this function makes calls to the Pi-Hole's API
func updateData() {
	summary.update()
	topItems.update()
	allQueries.update()
	topItems.prettyConvert()
	piCLIData.LastUpdated = time.Now()
}

// given a value representing the current privacy level, return the level name
// https://docs.pi-hole.net/ftldns/privacylevels/
func getPrivacyLevel(level *string) string {
	return summary.PrivacyLevelNumberMapping[*level]
}

// create and start the UI rendering
func startUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	piHoleInfo := widgets.NewList()
	piHoleInfo.Border = false

	totalQueries := widgets.NewParagraph()
	totalQueries.Title = "Queries /24hr"
	totalQueries.TitleStyle.Fg = ui.ColorGreen
	totalQueries.BorderStyle.Fg = ui.ColorGreen

	queriesBlocked := widgets.NewParagraph()
	queriesBlocked.Title = "Blocked /24hr"
	queriesBlocked.TitleStyle.Fg = ui.ColorBlue
	queriesBlocked.BorderStyle.Fg = ui.ColorBlue

	percentBlocked := widgets.NewParagraph()
	percentBlocked.Title = "Percent Blocked"
	percentBlocked.TitleStyle.Fg = ui.ColorYellow
	percentBlocked.BorderStyle.Fg = ui.ColorYellow

	domainsOnBlocklist := widgets.NewParagraph()
	domainsOnBlocklist.Title = "Blocklist Size"
	domainsOnBlocklist.TitleStyle.Fg = ui.ColorRed
	domainsOnBlocklist.BorderStyle.Fg = ui.ColorRed

	topQueries := widgets.NewList()
	topQueries.Title = "Top 10 Permitted Domains"
	topQueries.Rows = topItems.PrettyTopQueries

	topAds := widgets.NewList()
	topAds.Title = "Top 10 Blocked Domains"
	topAds.Rows = topItems.PrettyTopAds

	queryLog := widgets.NewTable()
	queryLog.Title = fmt.Sprintf("Latest %d queries", allQueries.AmountOfQueriesInLog)
	queryLog.Rows = [][]string{
		{"Time", "Type", "Domain", "Client", "Upstream"},
		{"1", "test", "www.google.com", "test", "test"},
	}

	grid := ui.NewGrid()
	w, h := ui.TerminalDimensions()
	grid.SetRect(0, 0, w, h)

	grid.Set(
		ui.NewRow(.2,
			ui.NewCol(.2,
				ui.NewRow(.5, totalQueries),
				ui.NewRow(.5, percentBlocked),
			),
			ui.NewCol(.2,
				ui.NewRow(.5, queriesBlocked),
				ui.NewRow(.5, domainsOnBlocklist),
			),
			ui.NewCol(.6,
				ui.NewRow(1, piHoleInfo),
			),
		),
		ui.NewRow(.4,
			ui.NewCol(.5, topQueries),
			ui.NewCol(.5, topAds),
		),
		ui.NewRow(.4,
			ui.NewCol(1, queryLog),
		),
	)

	draw := func() {
		// 4 top summary boxes
		totalQueries.Text = summary.QueriesToday
		queriesBlocked.Text = summary.BlockedToday
		percentBlocked.Text = summary.PercentBlockedToday + "%"
		domainsOnBlocklist.Text = summary.DomainsOnBlocklist

		// domain lists
		topQueries.Rows = topItems.PrettyTopQueries
		topAds.Rows = topItems.PrettyTopAds

		// status text
		formattedTime := piCLIData.LastUpdated.Format("15:04:05")

		piHoleInfo.Rows = []string{
			fmt.Sprintf("Pi-Hole Status: %s", strings.Title(summary.Status)),
			fmt.Sprintf("Data last updated: %s", formattedTime),
			fmt.Sprintf("Privacy Level: %s", getPrivacyLevel(&summary.PrivacyLevel)),
			fmt.Sprintf("Total Clients Seen: %s", summary.TotalClientsSeen),
		}

		// update the grid given the current terminal dimensions
		w, h := ui.TerminalDimensions()
		grid.SetRect(0, 0, w, h)

		ui.Render(grid)
	}

	uiEvents := ui.PollEvents()
	tickerDuration := time.Duration(piCLIData.Settings.RefreshS)
	ticker := time.NewTicker(time.Second * tickerDuration).C

	updateData()
	draw()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			updateData()
			draw()
		}
	}
}
