package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/Reeceeboii/Pi-CLI/pkg/database"
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
					Name:    "top-queries",
					Aliases: []string{"tq"},
					Usage:   "Extract the current top 10 permitted DNS queries",
					Action: func(context *cli.Context) error {
						InitialisePICLI()
						api.LiveTopItems.Update(nil)
						fmt.Printf("Top queries as of @ %s\n\n", time.Now().Format(time.Stamp))
						for _, q := range api.LiveTopItems.PrettyTopQueries {
							fmt.Println(q)
						}
						return nil
					},
				},
				{
					Name:    "top-ads",
					Aliases: []string{"ta"},
					Usage:   "Extract the current top 10 blocked domains",
					Action: func(c *cli.Context) error {
						InitialisePICLI()
						api.LiveTopItems.Update(nil)
						fmt.Printf("Top ads as of @ %s\n\n", time.Now().Format(time.Stamp))
						for _, q := range api.LiveTopItems.PrettyTopAds {
							fmt.Println(q)
						}
						return nil
					},
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
					Action: func(c *cli.Context) error {
						queryAmount := c.Int("queries")
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
					},
				},
				{
					Name:    "enable",
					Aliases: []string{"e"},
					Usage:   "Enable the Pi-Hole",
					Action: func(context *cli.Context) error {
						InitialisePICLI()
						api.LiveSummary.Update(nil)
						if api.LiveSummary.Status == "enabled" {
							fmt.Println("Pi-Hole is already enabled!")

						} else {
							api.EnablePiHole()
							fmt.Println("Pi-Hole enabled!")
						}

						return nil
					},
				},
				{
					Name:    "disable",
					Aliases: []string{"d"},
					Usage:   "Disable the Pi-Hole",
					Flags: []cli.Flag{
						&cli.Int64Flag{
							Name:        "timeout",
							Aliases:     []string{"t"},
							Usage:       "A timeout in seconds. Pi-Hole will re-enable when this time has elapsed.",
							DefaultText: "permanent",
						},
					},
					Action: func(context *cli.Context) error {
						InitialisePICLI()
						api.LiveSummary.Update(nil)
						if api.LiveSummary.Status == "disabled" {
							fmt.Println("Pi-Hole is already disabled!")
						} else {
							timeout := context.Int64("timeout")
							if timeout == 0 {
								api.DisablePiHole(false, 0)
								fmt.Println("Pi-Hole disabled until explicitly re-enabled")
							} else {
								api.DisablePiHole(true, timeout)
								fmt.Printf("Pi-Hole disabled. Will re-enable in %d seconds\n", timeout)
							}
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "database",
			Aliases: []string{"d"},
			Usage:   "Analytics options to run on a Pi-Hole's FTL database",

			/*
				FOR ALL DATABASE COMMANDS:
					If no path is provided by the user, Pi-CLI will assume that the database file's
					name hasn't been changed from it's default name, and that is has been placed in the
					same working directory that it is being executed from. This saves some command typing.
			*/
			Subcommands: []*cli.Command{
				{
					Name:    "client-summary",
					Aliases: []string{"cs"},
					Usage:   "Summary of all Pi-Hole clients",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "path",
							Aliases:     []string{"p"},
							Usage:       "Path to the Pi-Hole FTL database file",
							DefaultText: "./pihole-FTL.db",
						},
					},
					Action: func(context *cli.Context) error {
						path := context.String("path")
						if path == "" {
							path = "./pihole-FTL.db"
						}

						conn := database.Connect(path)
						database.ClientSummary(conn)
						return nil
					},
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
					Action: func(context *cli.Context) error {
						path := context.String("path")
						if path == "" {
							path = "./pihole-FTL.db"
						}

						conn := database.Connect(path)

						database.TopQueries(conn, context.Int64("limit"), context.String("filter"))

						return nil
					},
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
