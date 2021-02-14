package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"os"
	"time"
)

// stores the data needed by Pi-CLI during runtime
type PiCLIData struct {
	Settings            *Settings
	FormattedAPIAddress string
	APIKey              string
	LastUpdated         time.Time
}

var piCLIData = PiCLIData{}

var basicData = BasicData{
	QueriesToday:        0,
	BlockedToday:        0,
	PercentBlockedToday: 0.0,
	DomainsOnBlocklist:  0,
}

func startUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	totalQueries := widgets.NewParagraph()
	totalQueries.Title = "Total queries"
	totalQueries.SetRect(0, 0, 25, 3)

	queriesBlocked := widgets.NewParagraph()
	queriesBlocked.Title = "Queries Blocked"
	queriesBlocked.SetRect(25, 0, 50, 3)

	percentBlocked := widgets.NewParagraph()
	percentBlocked.Title = "Percent Blocked"
	percentBlocked.SetRect(50, 0, 75, 3)

	domainsOnBlocklist := widgets.NewParagraph()
	domainsOnBlocklist.Title = "Domains on Blocklist"
	domainsOnBlocklist.SetRect(75, 0, 100, 3)

	lastUpdated := widgets.NewParagraph()
	lastUpdated.SetRect(0, 3, 25, 6)
	lastUpdated.Border = false

	draw := func() {
		basicData.update()
		totalQueries.Text = fmt.Sprintf("%d", basicData.QueriesToday)
		queriesBlocked.Text = fmt.Sprintf("%d", basicData.BlockedToday)
		percentBlocked.Text = fmt.Sprintf("%f", basicData.PercentBlockedToday)
		domainsOnBlocklist.Text = fmt.Sprintf("%d", basicData.DomainsOnBlocklist)
		lastUpdated.Text = fmt.Sprintf("Last updated: %s", piCLIData.LastUpdated.Format("15:04:05"))
		ui.Render(totalQueries, queriesBlocked, percentBlocked, domainsOnBlocklist, lastUpdated)
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
			draw()
		}
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
