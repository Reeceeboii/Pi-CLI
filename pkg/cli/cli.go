package cli

import (
	"bufio"
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/api"
	"github.com/Reeceeboii/Pi-CLI/pkg/auth"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/database"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/Reeceeboii/Pi-CLI/pkg/settings"
	"github.com/Reeceeboii/Pi-CLI/pkg/ui"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// set up the cli app and all of its flags
var App = cli.App{
	EnableBashCompletion: true,
	Name:                 "Pi-CLI",
	Description:          `Pi-Hole data right from your terminal. Live updating view, query history extraction and more!`,
	Compiled:             time.Now(),
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
			// read in information from the user and create a config file with it
			Action: func(context *cli.Context) error {
				reader := bufio.NewReader(os.Stdin)

				addressDetailsValid := false
				for !addressDetailsValid {
					// read in the IP address and check that it is valid
					fmt.Print("Please enter the IP address of your Pi-Hole: ")
					piHoleAddress, _ := reader.ReadString('\n')
					ip := net.ParseIP(strings.TrimSpace(piHoleAddress))
					if ip == nil {
						log.Fatal("Please enter a valid IP address")
					}
					settings.PICLISettings.PiHoleAddress = ip.String()

					// read in the port
					fmt.Print("Please enter the port that exposes the web interface (default 80): ")
					piHolePort, _ := reader.ReadString('\n')
					trimmed := strings.TrimSpace(piHolePort)
					// if the user entered nothing, keep the default. Else, check and apply theirs
					if len(trimmed) > 0 {
						intPiHolePort, _ := strconv.Atoi(trimmed)
						if intPiHolePort < 1 || intPiHolePort > 65535 {
							log.Fatal("Please enter a valid port number")
						}
						settings.PICLISettings.PiHolePort = intPiHolePort
					}

					// send a request to the PiHole to validate that the IP and port actually point to it
					tempURL := fmt.Sprintf(
						"http://%s:%d/admin/api.php",
						settings.PICLISettings.PiHoleAddress,
						settings.PICLISettings.PiHolePort)
					req, err := http.NewRequest("GET", tempURL, nil)
					if err != nil {
						log.Fatal(err)
					}
					res, err := network.HttpClient.Do(req)

					// if the details are valid and the request didn't time out...
					// lazy evaluation saves us from deref errors here and saves a check
					if err == nil && network.ValidatePiHoleDetails(res) {
						addressDetailsValid = true
					} else {
						color.Yellow("Pi-Hole doesn't seem to be alive, check your details and try again!")
						fmt.Println()
					}
				}

				// read in the data refresh rate
				fmt.Print("Please enter your preferred data refresh rate in seconds (default 1s): ")
				refreshS, _ := reader.ReadString('\n')
				trimmed := strings.TrimSpace(refreshS)
				if len(trimmed) > 0 {
					intRefreshS, err := strconv.Atoi(trimmed)
					if err != nil {
						log.Fatal("Please enter a number")
					}
					if intRefreshS < 1 {
						log.Fatal("Refresh time cannot be less than 1 second")
					}
					settings.PICLISettings.RefreshS = intRefreshS
				}

				for {
					// read in the API key and work out where the user wants to store it (keyring or config file)
					fmt.Print("Please enter your Pi-Hole API key: ")
					apiKey, _ := reader.ReadString('\n')
					apiKey = strings.TrimSpace(apiKey)
					if len(apiKey) < 1 {
						fmt.Println("Please provide your API key for authentication")
						continue
					}

					settings.PICLISettings.APIKey = apiKey

					// before we store the API token (keyring or config file), we should check that it's valid
					// the address + port have been validated by this point so we're safe to shoot requests at it
					data.LivePiCLIData.Settings = settings.PICLISettings
					data.LivePiCLIData.FormattedAPIAddress = network.GenerateAPIAddress(
						settings.PICLISettings.PiHoleAddress,
						settings.PICLISettings.PiHolePort)

					if !auth.ValidateAPIKey(settings.PICLISettings.APIKey) {
						color.Yellow("That API token doesn't seem to be correct, check it and try again!")
					} else {
						break
					}
				}

				fmt.Print("Do you wish to store the API key in your system keyring? (y/n - default y): ")
				storageChoice, _ := reader.ReadString('\n')
				storageChoice = strings.ToLower(strings.TrimSpace(storageChoice))

				// if they wish to use their system's keyring...
				if storageChoice == "y" || len(storageChoice) == 0 {
					err := auth.StoreAPIKeyInKeyring(settings.PICLISettings.APIKey)

					if err == nil {
						color.Green("Your API token has been securely stored in your system keyring")
						/*
							After the API key has been saved to the keyring, there is no longer a need to save it
							to the config file, so the stored copy of it can be removed from the in-memory settings
							instance before it gets serialised to disk
						*/
						settings.PICLISettings.APIKey = ""
					} else {
						color.Yellow("System keyring call failed, falling back to config file")
					}
				}

				// write config file to disk
				// all fields in the settings struct would have been set by this point
				settings.PICLISettings.SaveToFile()
				color.Green("Configuration successful")
				return nil
			},
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
					Action: func(context *cli.Context) error {
						if auth.DeleteAPIKeyFromKeyring() {
							color.Green("System keyring API entry has been deleted!")
						} else {
							color.Yellow("Pi-CLI did not find a keyring entry to delete")
						}
						if settings.DeleteConfigFile() {
							color.Green("Stored config file has been deleted!")
						} else {
							color.Yellow("Pi-CLI did not find a config file to delete")
						}
						return nil
					},
				},
				{
					Name:    "view",
					Aliases: []string{"v"},
					Usage:   "View config stored config data (config file and API key)",
					Action: func(context *cli.Context) error {
						// if the config file is present, that can be loaded and displayed,
						// otherwise, prompt the user to create one
						if settings.ConfigFileExists() {
							settings.PICLISettings.LoadFromFile()
							fmt.Printf(
								"%s%s\n",
								"Pi-Hole address: ",
								settings.PICLISettings.PiHoleAddress)
							fmt.Printf(
								"%s%d\n",
								"Pi-Hole port: ",
								settings.PICLISettings.PiHolePort)
							fmt.Printf("%s%d%s\n",
								"Refresh rate: ",
								settings.PICLISettings.RefreshS,
								"s")
						} else {
							color.Yellow("No config file is present - run the setup command to create one")
						}

						// and the same with the API key
						if auth.APIKeyIsInKeyring() {
							fmt.Printf("%s%s\n", "API key (keyring): ", auth.RetrieveAPIKeyFromKeyring())
						} else if settings.PICLISettings.APIKeyIsInFile() {
							fmt.Printf("%s%s\n", "API key (config file): ", settings.PICLISettings.APIKey)
						} else {
							color.Yellow("No API key has been provided - run the setup command to enter it")
						}

						return nil
					},
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
					Action: func(context *cli.Context) error {
						initialisePICLI()
						api.LiveSummary.Update(nil)
						fmt.Printf("Summary @ %s\n", time.Now().Format(time.Stamp))
						fmt.Println()
						fmt.Printf("Pi-Hole status: %s\n", strings.Title(api.LiveSummary.Status))
						fmt.Println()
						fmt.Printf("Queries /24hr: %s\n", api.LiveSummary.QueriesToday)
						fmt.Printf("Blocked /24hr: %s\n", api.LiveSummary.BlockedToday)
						fmt.Printf("Percent blocked: %s%s\n", api.LiveSummary.PercentBlockedToday, "%")
						fmt.Printf("Domains on blocklist: %s\n", api.LiveSummary.DomainsOnBlocklist)
						fmt.Printf("Privacy level: %s - %s\n",
							api.LiveSummary.PrivacyLevel,
							api.LiveSummary.PrivacyLevelNumberMapping[api.LiveSummary.PrivacyLevel],
						)
						fmt.Printf("Total clients seen: %s\n", api.LiveSummary.TotalClientsSeen)
						fmt.Println()
						return nil
					},
				},
				{
					Name:    "top-queries",
					Aliases: []string{"tq"},
					Usage:   "Extract the current top 10 permitted DNS queries",
					Action: func(context *cli.Context) error {
						initialisePICLI()
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
						initialisePICLI()
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
						initialisePICLI()
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
						initialisePICLI()
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
						initialisePICLI()
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
							Usage:       "Filter by domain or word. (e.g. 'google.com', 'spotify', 'adservice' etc...)",
							DefaultText: "No filter",
						},
					},
					Action: func(context *cli.Context) error {
						path := context.String("path")
						if path == "" {
							path = "./pihole-FTL.db"
						}

						conn := database.Connect(path)

						limit := context.Int64("limit")
						if limit == 0 {
							limit = 10
						}

						database.TopQueries(conn, limit, context.String("filter"))
						return nil
					},
				},
			},
		},
	},

	Action: func(context *cli.Context) error {
		initialisePICLI()
		ui.StartUI()
		return nil
	},
}

/*
	Validate that the config file and API key are in place.
	Load the required data and settings into memory
*/
func initialisePICLI() {
	// firstly, has a config file been created?
	if !settings.ConfigFileExists() {
		log.Fatal("Please configure Pi-CLI via the 'setup' command")
	}

	settings.PICLISettings.LoadFromFile()

	// retrieve the API key depending upon its storage location
	if !settings.PICLISettings.APIKeyIsInFile() && !auth.APIKeyIsInKeyring() {
		log.Fatal("Please configure Pi-CLI via the 'setup' command")
	} else {
		if settings.PICLISettings.APIKeyIsInFile() {
			data.LivePiCLIData.APIKey = settings.PICLISettings.APIKey
		} else {
			data.LivePiCLIData.APIKey = auth.RetrieveAPIKeyFromKeyring()
		}
	}

	data.LivePiCLIData.Settings = settings.PICLISettings
	data.LivePiCLIData.FormattedAPIAddress = network.GenerateAPIAddress(
		settings.PICLISettings.PiHoleAddress,
		settings.PICLISettings.PiHolePort)
}
