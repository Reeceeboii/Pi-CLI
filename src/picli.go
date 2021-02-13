package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"os"
	"time"
)

// stores the data needed by Pi-CLI during runtime
type PiCLIData struct {
	settings *Settings
	APIKey   string
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

	draw := func() {
		ui.Render(totalQueries, queriesBlocked, percentBlocked, domainsOnBlocklist)
	}

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second / 5).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw()
		}
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
