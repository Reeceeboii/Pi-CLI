package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"strings"
	"time"
)

const (
	Quarter = 1.0 / 4
	Tenth   = 1.0 / 10
)

// update data so it can be displayed
// this function makes calls to the Pi-Hole's API
func updateData() {
	summary.update()
	topItems.update()
	topItems.prettyConvert()
	piCLIData.LastUpdated = time.Now()
}

// create and start the UI rendering
func startUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

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

	bottomText := widgets.NewParagraph()
	bottomText.Border = false

	grid := ui.NewGrid()
	w, h := ui.TerminalDimensions()
	grid.SetRect(0, 0, w, h)
	grid.Set(
		ui.NewRow(Tenth,
			ui.NewCol(.25, totalQueries),
			ui.NewCol(.25, queriesBlocked),
			ui.NewCol(.25, percentBlocked),
			ui.NewCol(.25, domainsOnBlocklist),
		),
		ui.NewRow(.37,
			ui.NewCol(.5, topQueries),
			ui.NewCol(.5, topAds),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(.5, bottomText),
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

		// bottom text
		formattedTime := piCLIData.LastUpdated.Format("15:04:05")
		if summary.Status == "enabled" {
			bottomText.TextStyle.Fg = ui.ColorGreen
		} else {
			bottomText.TextStyle.Fg = ui.ColorRed
		}
		bottomText.Text = "Last updated: " + formattedTime + " | Status: " + strings.Title(summary.Status)

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
