package ui

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

/*
	Update data so it can be displayed.
	This function makes calls to the Pi-Hole's API
*/
func updateData() {
	go api.LiveSummary.Update(&wg)
	go api.LiveTopItems.Update(&wg)
	go api.LiveAllQueries.Update(&wg)
	wg.Wait()
	data.LivePiCLIData.LastUpdated = time.Now()
}

/*
	Given a value representing the current privacy level, return the level name.
	https://docs.pi-hole.net/ftldns/privacylevels/
*/
func getPrivacyLevel(level *string) string {
	return api.LiveSummary.PrivacyLevelNumberMapping[*level]
}

/*
	Is the UI free to draw to? Currently this only takes into account the fact
	that the keybinds view may be showing. Adding more conditions for halting live
	UI redraws is as simple as ANDing them here
*/
func uiCanDraw() bool {
	return !data.LivePiCLIData.ShowKeybindsScreen
}

// Create the UI and start rendering
func StartUI() {
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
	topQueries.Rows = api.LiveTopItems.PrettyTopQueries

	topAds := widgets.NewList()
	topAds.Title = "Top 10 Blocked Domains"
	topAds.Rows = api.LiveTopItems.PrettyTopAds

	queryLog := widgets.NewList()
	queryLog.Title = fmt.Sprintf("Latest %d queries", api.LiveAllQueries.AmountOfQueriesInLog)
	queryLog.Rows = api.LiveAllQueries.Table

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
	keybindsList.Rows = data.LivePiCLIData.Keybinds

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
		if uiCanDraw() {
			// 4 top summary boxes
			totalQueries.Text = api.LiveSummary.QueriesToday
			queriesBlocked.Text = api.LiveSummary.BlockedToday
			percentBlocked.Text = api.LiveSummary.PercentBlockedToday + "%"
			domainsOnBlocklist.Text = api.LiveSummary.DomainsOnBlocklist

			// domain lists
			topQueries.Rows = api.LiveTopItems.PrettyTopQueries
			topAds.Rows = api.LiveTopItems.PrettyTopAds

			// query log
			queryLog.Rows = api.LiveAllQueries.Table
			queryLog.Title = fmt.Sprintf("Latest %d queries", api.LiveAllQueries.AmountOfQueriesInLog)

			// status text
			formattedTime := data.LivePiCLIData.LastUpdated.Format("15:04:05")

			piHoleInfo.Rows = []string{
				fmt.Sprintf("Pi-Hole Status: %s", strings.Title(api.LiveSummary.Status)),
				fmt.Sprintf(
					"Data last updated: %s (update every %ds)",
					formattedTime,
					data.LivePiCLIData.Settings.RefreshS),
				fmt.Sprintf("Privacy Level: %s", getPrivacyLevel(&api.LiveSummary.PrivacyLevel)),
				fmt.Sprintf("Total Clients Seen: %s", api.LiveSummary.TotalClientsSeen),
			}

			// render the grid
			ui.Render(grid)
		} else {
			ui.Render(keybindsGrid)
		}
	}

	uiEvents := ui.PollEvents()

	// channel used to capture ticker events to time data update events
	tickerDuration := time.Duration(data.LivePiCLIData.Settings.RefreshS)
	dataUpdateTicker := time.NewTicker(time.Second * tickerDuration).C

	// channel used to capture ticker events to time redraws
	drawTicker := time.NewTicker(time.Second / time.Duration(data.LivePiCLIData.Settings.UIFramesPerSecond)).C

	updateData()
	draw()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {

			// quit
			case "q", "<C-c>":
				return

			// respond to terminal resize events
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				if uiCanDraw() {
					grid.SetRect(0, 0, payload.Width, payload.Height)
					ui.Render(grid)
					break
				}
				keybindsGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(keybindsGrid)
				break

			// increase (by 1) the number of queries in the query log
			case "e":
				if uiCanDraw() {
					api.LiveAllQueries.AmountOfQueriesInLog++
					api.LiveAllQueries.Queries = append(api.LiveAllQueries.Queries, api.Query{})
				}
				break

			// increase (by 10) the number of queries in the query log
			case "r":
				if uiCanDraw() {
					api.LiveAllQueries.AmountOfQueriesInLog += 10
					api.LiveAllQueries.Queries = append(api.LiveAllQueries.Queries, make([]api.Query, 10)...)
				}
				break

			// decrease (by 1) the number of queries in the query log
			case "d":
				if uiCanDraw() && api.LiveAllQueries.AmountOfQueriesInLog > 1 {
					api.LiveAllQueries.AmountOfQueriesInLog--
					api.LiveAllQueries.Queries = api.LiveAllQueries.Queries[:len(api.LiveAllQueries.Queries)-1]
				}
				break

			// decrease (by 10) the number of queries in the query log
			case "f":
				if uiCanDraw() {
					if api.LiveAllQueries.AmountOfQueriesInLog-10 <= 0 {
						api.LiveAllQueries.AmountOfQueriesInLog = 1
						api.LiveAllQueries.Queries =
							api.LiveAllQueries.Queries[:len(api.LiveAllQueries.Queries)-(len(api.LiveAllQueries.Queries)-1)]
					} else {
						api.LiveAllQueries.AmountOfQueriesInLog -= 10
						api.LiveAllQueries.Queries = api.LiveAllQueries.Queries[:len(api.LiveAllQueries.Queries)-10]
					}
				}
				break

			// scroll down (by 1) in the query log list
			case "<Down>":
				if uiCanDraw() {
					queryLog.ScrollDown()
				}
				break

			// scroll down (by 10) in the query log list
			case "<PageDown>":
				if uiCanDraw() {
					queryLog.ScrollAmount(10)
				}
				break

			// scroll up (by 1) in the query log list
			case "<Up>":
				if uiCanDraw() {
					queryLog.ScrollUp()
				}
				break

			// scroll up (by 10) in the query log list
			case "<PageUp>":
				if uiCanDraw() {
					queryLog.ScrollAmount(-10)
				}
				break

			// enable or disable the Pi-Hole
			case "p":
				if uiCanDraw() {
					if api.LiveSummary.Status == "enabled" {
						api.DisablePiHole(false, 0)
					} else {
						api.EnablePiHole()
					}
				}
				break

			// switch grids between the keybinds view and the main screen
			case "<F1>":
				ui.Clear()
				data.LivePiCLIData.ShowKeybindsScreen = !data.LivePiCLIData.ShowKeybindsScreen
				break
			}

		/*
			Capturing 2 separate ticker channels like this allows the update of the data and the update of the
			UI to occur independently. Key presses will still be visually responded to and the program itself will
			*feel* quick and responsive even if the user has set a much longer data refresh rate.
		*/

		// refresh event used to time API polls for up to date data
		case <-dataUpdateTicker:
			// there's only a need to make API calls when the keybinds screen isn't being shown
			if uiCanDraw() {
				updateData()
			}
			break

		// draw event used to time UI redraws
		case <-drawTicker:
			draw()
			break
		}
	}
}
