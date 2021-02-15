package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"time"
)

func startUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	totalQueries := widgets.NewParagraph()
	totalQueries.Title = "Queries /24hr"
	totalQueries.TitleStyle.Fg = ui.ColorGreen
	totalQueries.BorderStyle.Fg = ui.ColorGreen
	totalQueries.SetRect(0, 0, 25, 3)

	queriesBlocked := widgets.NewParagraph()
	queriesBlocked.Title = "Blocked /24hr"
	queriesBlocked.TitleStyle.Fg = ui.ColorBlue
	queriesBlocked.BorderStyle.Fg = ui.ColorBlue
	queriesBlocked.SetRect(25, 0, 50, 3)

	percentBlocked := widgets.NewParagraph()
	percentBlocked.Title = "Percent Blocked"
	percentBlocked.TitleStyle.Fg = ui.ColorYellow
	percentBlocked.BorderStyle.Fg = ui.ColorYellow
	percentBlocked.SetRect(50, 0, 75, 3)

	domainsOnBlocklist := widgets.NewParagraph()
	domainsOnBlocklist.Title = "Domains on Blocklist"
	domainsOnBlocklist.TitleStyle.Fg = ui.ColorRed
	domainsOnBlocklist.BorderStyle.Fg = ui.ColorRed
	domainsOnBlocklist.SetRect(75, 0, 100, 3)

	topQueries := widgets.NewList()
	topQueries.Title = "Top Permitted Domains"
	topQueries.Rows = topItems.PrettyTopQueries
	topQueries.SetRect(0, 3, 50, 15)

	topAds := widgets.NewList()
	topAds.Title = "Top Blocked Domains"
	topAds.Rows = topItems.PrettyTopAds
	topAds.SetRect(50, 3, 100, 15)

	lastUpdated := widgets.NewParagraph()
	lastUpdated.SetRect(0, 20, 25, 23)
	lastUpdated.Border = false

	draw := func() {
		totalQueries.Text = summary.QueriesToday
		queriesBlocked.Text = summary.BlockedToday
		percentBlocked.Text = summary.PercentBlockedToday + "%"
		domainsOnBlocklist.Text = summary.DomainsOnBlocklist
		topQueries.Rows = topItems.PrettyTopQueries
		topAds.Rows = topItems.PrettyTopAds
		lastUpdated.Text = fmt.Sprintf("Last updated: %s", piCLIData.LastUpdated.Format("15:04:05"))

		ui.Render(
			totalQueries,
			queriesBlocked,
			percentBlocked,
			domainsOnBlocklist,
			topQueries,
			topAds,
			lastUpdated,
		)
	}

	uiEvents := ui.PollEvents()
	tickerDuration := time.Duration(piCLIData.Settings.RefreshS)
	ticker := time.NewTicker(time.Second * tickerDuration).C
	draw()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			piCLIData.LastUpdated = time.Now()
			summary.update()
			topItems.update()
			topItems.prettyConvert()
			draw()
		}
	}
}
