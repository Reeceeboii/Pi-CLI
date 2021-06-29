package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/database"
	"github.com/Reeceeboii/Pi-CLI/pkg/logger"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/Reeceeboii/Pi-CLI/pkg/ui"
	"github.com/Reeceeboii/Pi-CLI/pkg/update"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"log"
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
	Copyright:   fmt.Sprintf("Copyright (C) %d Reece Mercer", time.Now().Year()),
	Before:      initialiseWithGlobals,
	Authors: []*cli.Author{
		{
			Name:  "Reece Mercer",
			Email: "reecemercer981@gmail.com",
		},
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "environment",
			Aliases:     []string{"env"},
			Usage:       "Load .env file",
			DefaultText: "false",
		},
		&cli.BoolFlag{
			Name:        "log",
			Aliases:     []string{"l"},
			Usage:       "Enable debug logging. Saves to user's home directory",
			DefaultText: "false",
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
							Name:        "limit",
							Aliases:     []string{"l"},
							Usage:       "The limit on the number of queries to extract",
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
							DefaultText: database.DefaultDatabaseFileLocation,
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
							DefaultText: database.DefaultDatabaseFileLocation,
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
					Action: RunDatabaseTopQueriesCommand,
				},
			},
		},
	},

	Action: func(c *cli.Context) error {
		InitialisePICLI()
		ui.StartUI()
		return nil
	},
}

// Carries out some preliminary app initialisation with the global arguments
func initialiseWithGlobals(c *cli.Context) error {
	/*
		If logging flag is given, we can enable the logger. If is it not enabled (it is
		initialised disabled by default), all logging statements using logger.PiCLIFileLogger
		will simply have no effect.
	*/
	if c.Bool("log") {
		logger.LivePiCLILogger.Enabled = true
		logger.LivePiCLILogger.LogStatus("Logging has been enabled")
	}

	logger.LivePiCLILogger.LogStatus("Pi-CLI started")
	logger.LivePiCLILogger.LogStartupInformation()

	// If the "environment" flag has been enabled, we want to read in environment variables
	if c.Bool("environment") {
		logger.LivePiCLILogger.LogInformation("Environment variables flag provided, loading from environment")
		if err := godotenv.Load(); err != nil {
			color.Red("Failed to load .env")
			logger.LivePiCLILogger.LogError(".env file not found, environment variables could not be loaded")
		}
	}

	// load in the config file if it exists
	if data.ConfigFileExists(data.GetConfigFileLocation()) {
		data.PICLISettings.LoadFromFile(data.GetConfigFileLocation())
	} else {
		return nil
	}

	if data.PICLISettings.AutoCheckForUpdates {
		logger.LivePiCLILogger.LogInformation("Starting update check")
		latestReleaseCheckedTime, err := time.Parse(time.RFC3339, data.PICLISettings.LatestRemoteRelease.TimeChecked)
		if err != nil {
			logger.LivePiCLILogger.LogError("Error parsing latest release time: " + err.Error())
			return err
		} else {
			// if previous update check was carried out > 20 minutes ago
			if latestReleaseCheckedTime.Before(time.Now().Add(time.Minute * -20)) {
				logger.LivePiCLILogger.LogInformation("Previous check > 20 minutes prior - contacting GitHub")
				data.PICLISettings.LatestRemoteRelease = update.GetLatestGitHubRelease(network.HttpClient)

				log.Println(data.PICLISettings.LatestRemoteRelease)
				//if err := data.PICLISettings.SaveToFile(); err != nil {
				//	color.Red("Unable to access log file!")
				//	log.Fatal(err)
				//}
			}
		}
	}

	return nil
}
