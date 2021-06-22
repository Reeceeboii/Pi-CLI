package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/ui"
	"github.com/urfave/cli/v2"
	"time"
)

/*
	This is the main CLI app, it contains all of the various commands and subcommands
	that Pi-CLI is capable of responding to, and manages all of their corresponding flags
*/
var App = cli.App{
	Name:        "Pi-CLI",
	Usage:       "Third party program to retrieve and display Pi-Hole data right from your terminal",
	Description: "Pi-Hole data right from your terminal. Live updating view, query history extraction and more!",
	Copyright:   fmt.Sprintf("Copyright (c) %d Reece Mercer", time.Now().Year()),
	Authors: []*cli.Author{
		{
			Name:  "Reece Mercer",
			Email: "reecemercer981@gmail.com",
		},
	},
	Commands: []*cli.Command{
		{
			Name:    "setup",
			Aliases: []string{"s"},
			Usage:   "Configure Pi-CLI",
			Action:  SetupCommand,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Interact with stored configuration settings",
			Subcommands: []*cli.Command{
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Delete stored config data (config file and API key)",
					Action:  ConfigDeleteCommand,
				},
				{
					Name:    "view",
					Aliases: []string{"v"},
					Usage:   "View config stored config data (config file and API key)",
					Action:  ConfigViewCommand,
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a one off command without booting the live view",
			Subcommands: []*cli.Command{
				{
					Name:    "summary",
					Aliases: []string{"s"},
					Usage:   "Extract a basic summary of data from the Pi-Hole",
					Action:  RunSummaryCommand,
				},
				{
					Name:    "top-forwarded",
					Aliases: []string{"tf"},
					Usage:   "Extract the current top 10 forwarded DNS queries",
					Action:  RunTopTenForwardedCommand,
				},
				{
					Name:    "top-blocked",
					Aliases: []string{"tb"},
					Usage:   "Extract the current top 10 blocked DNS queries",
					Action:  RunTopTenBlockedCommand,
				},
				{
					Name:    "latest-queries",
					Aliases: []string{"lq"},
					Usage:   "Extract the latest queries",
					Flags: []cli.Flag{
						&cli.Int64Flag{
							Name:        "queries",
							Aliases:     []string{"q"},
							Usage:       "The number of queries to extract",
							DefaultText: "10",
						},
					},
					Action: RunLatestQueriesCommand,
				},
				{
					Name:    "enable",
					Aliases: []string{"e"},
					Usage:   "Enable the Pi-Hole",
					Action:  RunEnablePiHoleCommand,
				},
				{
					Name:    "disable",
					Aliases: []string{"d"},
					Usage:   "Disable the Pi-Hole",
					Flags: []cli.Flag{
						&cli.Int64Flag{
							Name:        "timeout",
							Aliases:     []string{"t"},
							Usage:       "A timeout in seconds. Pi-Hole will re-enable after this time has elapsed.",
							DefaultText: "permanent",
						},
					},
					Action: RunDisablePiHoleCommand,
				},
			},
		},
		{
			Name:    "database",
			Aliases: []string{"d"},
			Usage:   "Analytics options to run on a Pi-Hole's FTL database",
			Subcommands: []*cli.Command{
				{
					Name:    "client-summary",
					Aliases: []string{"cs"},
					Usage:   "Summary of all Pi-Hole clients",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "path",
							Aliases:     []string{"p"},
							Usage:       "Path to a Pi-Hole FTL database file",
							DefaultText: "./pihole-FTL.db",
						},
					},
					Action: RunDatabaseClientSummaryCommand,
				},
				{
					Name:    "top-queries",
					Aliases: []string{"tq"},
					Usage:   "Returns the top (all time) queries",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "path",
							Aliases:     []string{"p"},
							Usage:       "Path to the Pi-Hole FTL database file",
							DefaultText: "./pihole-FTL.db",
						},
						&cli.Int64Flag{
							Name:        "limit",
							Aliases:     []string{"l"},
							Usage:       "The limit on the number of queries to extract",
							DefaultText: "10",
						},
						&cli.StringFlag{
							Name:        "filter",
							Aliases:     []string{"f"},
							Usage:       "Filter by domain or word. (e.g. 'google.com', 'spotify', 'facebook' etc...)",
							DefaultText: "No filter",
						},
					},
					Action: RunDatabaseClientSummaryCommand,
				},
			},
		},
	},

	Action: func(context *cli.Context) error {
		InitialisePICLI()
		ui.StartUI()
		return nil
	},
}
