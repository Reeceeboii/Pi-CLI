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
	go summary.update()
	go topItems.update()
	go allQueries.update()
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

	queryLog := widgets.NewList()
	queryLog.Title = fmt.Sprintf("Latest %d queries", allQueries.AmountOfQueriesInLog)
	queryLog.Rows = allQueries.Table

	keybindsPrompt := widgets.NewParagraph()
	keybindsPrompt.Text = "Press F1 at any time to view keybinds..."
	keybindsPrompt.Border = false

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
		ui.NewRow(.35,
			ui.NewCol(.5, topQueries),
			ui.NewCol(.5, topAds),
		),
		ui.NewRow(.35,
			ui.NewCol(1, queryLog),
		),
		ui.NewRow(.1,
			ui.NewCol(1, keybindsPrompt),
		),
	)

	keybindsList := widgets.NewList()
	keybindsList.Title = "Pi-CLI keybinds"
	keybindsList.Rows = piCLIData.Keybinds

	returnHomePrompt := widgets.NewParagraph()
	returnHomePrompt.Text = "Press F1 at any time to return home..."
	returnHomePrompt.Border = false

	keybindsGrid := ui.NewGrid()
	w, h = ui.TerminalDimensions()
	keybindsGrid.SetRect(0, 0, w, h)
	keybindsGrid.Set(
		ui.NewRow(.9,
			ui.NewCol(1, keybindsList),
		),
		ui.NewRow(.1,
			ui.NewCol(1, returnHomePrompt),
		),
	)

	draw := func() {
		if !piCLIData.ShowKeybindsScreen {
			// 4 top summary boxes
			totalQueries.Text = summary.QueriesToday
			queriesBlocked.Text = summary.BlockedToday
			percentBlocked.Text = summary.PercentBlockedToday + "%"
			domainsOnBlocklist.Text = summary.DomainsOnBlocklist

			// domain lists
			topQueries.Rows = topItems.PrettyTopQueries
			topAds.Rows = topItems.PrettyTopAds

			// query log
			queryLog.Rows = allQueries.Table
			queryLog.Title = fmt.Sprintf("Latest %d queries", allQueries.AmountOfQueriesInLog)

			// status text
			formattedTime := piCLIData.LastUpdated.Format("15:04:05")

			piHoleInfo.Rows = []string{
				fmt.Sprintf("Pi-Hole Status: %s", strings.Title(summary.Status)),
				fmt.Sprintf("Data last updated: %s (update every %ds)", formattedTime, piCLIData.Settings.RefreshS),
				fmt.Sprintf("Privacy Level: %s", getPrivacyLevel(&summary.PrivacyLevel)),
				fmt.Sprintf("Total Clients Seen: %s", summary.TotalClientsSeen),
			}

			// update the grid given the current terminal dimensions
			w, h := ui.TerminalDimensions()
			grid.SetRect(0, 0, w, h)
			ui.Render(grid)
		} else {
			ui.Render(keybindsGrid)
		}
	}

	uiEvents := ui.PollEvents()
	tickerDuration := time.Duration(piCLIData.Settings.RefreshS)
	ticker := time.NewTicker(time.Second * tickerDuration).C

	updateData()
	for {
		draw()
		select {
		case e := <-uiEvents:
			switch e.ID {

			// quit
			case "q", "<C-c>":
				return

			// respond to terminal resize events
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				if !piCLIData.ShowKeybindsScreen {
					grid.SetRect(0, 0, payload.Width, payload.Height)
					ui.Clear()
					ui.Render(grid)
					break
				}
				keybindsGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(keybindsGrid)
				break

			// increase number of queries in query log by 1
			case "e":
				if !piCLIData.ShowKeybindsScreen {
					allQueries.AmountOfQueriesInLog++
					allQueries.Queries = append(allQueries.Queries, Query{})
				}
				break

			// decrease number of queries in log by 1
			case "d":
				if !piCLIData.ShowKeybindsScreen && allQueries.AmountOfQueriesInLog > 1 {
					allQueries.AmountOfQueriesInLog--
					allQueries.Queries = allQueries.Queries[:len(allQueries.Queries)-1]
				}
				break

			// scroll down in the query log list
			case "<Down>":
				if !piCLIData.ShowKeybindsScreen {
					queryLog.ScrollDown()
				}
				break

			// scroll up in the query log list
			case "<Up>":
				if !piCLIData.ShowKeybindsScreen {
					queryLog.ScrollUp()
				}
				break

			// switch grids between the keybinds view and the main screen
			case "<F1>":
				ui.Clear()
				piCLIData.ShowKeybindsScreen = !piCLIData.ShowKeybindsScreen
				break
			}

		// refresh event
		case <-ticker:
			// there's only a need to make API calls when the keybinds screen isn't being shown
			if !piCLIData.ShowKeybindsScreen {
				updateData()
			}
		}
	}
}
