package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var settings = Settings{
	PiHolePort: defaultPort,
	RefreshS:   defaultRefreshS,
}

// set up the cli app and all of its flags
var app = cli.App{
	EnableBashCompletion: true,
	Name:                 "Pi-CLI",
	Usage:                "Third party program to retrieve and display Pi-Hole data right from your terminal.",
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
			Action: func(c *cli.Context) error {
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
					settings.PiHoleAddress = ip.String()

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
						settings.PiHolePort = intPiHolePort
					}

					// send a request to the PiHole to validate that the IP and port actually point to it
					tempURL := fmt.Sprintf("http://%s:%d/admin/api.php", settings.PiHoleAddress, settings.PiHolePort)
					req, err := http.NewRequest("GET", tempURL, nil)
					if err != nil {
						log.Fatal(err)
					}
					res, err := client.Do(req)

					// if the details are valid and the request didn't time out...
					// lazy evaluation saves us from deref errors here and saves a check
					if err == nil && validatePiHoleDetails(res) {
						addressDetailsValid = true
					} else {
						fmt.Println("Pi-Hole doesn't seem to be alive, check your details and try again!")
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
					settings.RefreshS = intRefreshS
				}

				// read in the API key and work out where the user wants to store it (keyring or config file)
				fmt.Print("Please enter your Pi-Hole API key: ")
				apiKey, _ := reader.ReadString('\n')
				apiKey = strings.TrimSpace(apiKey)
				if len(apiKey) < 1 {
					log.Fatal("Please provide your API key for authentication")
				}

				fmt.Print("Do you wish to store the API key in your system keyring? (y/n - default y): ")
				storageChoice, _ := reader.ReadString('\n')
				storageChoice = strings.ToLower(strings.TrimSpace(storageChoice))

				// if they wish to use their system's keyring...
				if storageChoice == "y" || len(storageChoice) == 0 {
					storeAPIKeyInKeyring(&apiKey)
					fmt.Println("Your API token has been securely stored in your system keyring")
				} else {
					settings.APIKey = apiKey
				}

				// write config file to disk
				settings.saveToFile()
				fmt.Println("Configuration successful!")
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
						if deleteAPIKeyFromKeyring() {
							fmt.Println("System keyring API entry has been deleted!")
						} else {
							fmt.Println("Pi-CLI did not find a keyring entry to delete")
						}
						if deleteConfigFile() {
							fmt.Println("Stored config file has been deleted!")
						} else {
							fmt.Println("Pi-CLI did not find a config file to delete")
						}
						return nil
					},
				},
				{
					Name:    "view",
					Aliases: []string{"v"},
					Usage:   "View config stored config data (config file and API key)",
					Action: func(context *cli.Context) error {
						settings.loadFromFile()
						// if the config file is present, that can be loaded and displayed
						if configFileExists() {
							settings.loadFromFile()
							fmt.Printf("%s%s\n", "Pi-Hole address: ", settings.PiHoleAddress)
							fmt.Printf("%s%d\n", "Pi-Hole port: ", settings.PiHolePort)
							fmt.Printf("%s%d%s\n", "Refresh rate: ", settings.RefreshS, "s")
						} else {
							fmt.Println("No config file is present - run the setup command to create one")
						}

						// and the same with the API key
						if APIKeyIsInKeyring() {
							fmt.Printf("%s%s\n", "API key (keyring): ", retrieveAPIKeyFromKeyring())
						} else if settings.APIKeyIsInFile() {
							fmt.Printf("%s%s\n", "API key (config file): ", settings.APIKey)
						} else {
							fmt.Println("No API key has been provided - run the setup command to enter it")
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
					Action: func(c *cli.Context) error {
						initialisePICLI()
						summary.update()
						fmt.Printf("Summary @ %s\n", time.Now().Format(time.Stamp))
						fmt.Println()
						fmt.Printf("Pi-Hole status: %s\n", strings.Title(summary.Status))
						fmt.Println()
						fmt.Printf("Queries /24hr: %s\n", summary.QueriesToday)
						fmt.Printf("Blocked /24hr: %s\n", summary.BlockedToday)
						fmt.Printf("Percent blocked: %s%s\n", summary.PercentBlockedToday, "%")
						fmt.Printf("Domains on blocklist: %s\n", summary.DomainsOnBlocklist)
						fmt.Printf("Privacy level: %s - %s\n",
							summary.PrivacyLevel,
							summary.PrivacyLevelNumberMapping[summary.PrivacyLevel],
						)
						fmt.Printf("Total clients seen: %s\n", summary.TotalClientsSeen)
						fmt.Println()
						return nil
					},
				},
				{
					Name:    "top-queries",
					Aliases: []string{"tq"},
					Usage:   "Extract the current top 10 permitted DNS queries",
					Action: func(c *cli.Context) error {
						initialisePICLI()
						topItems.update()
						fmt.Printf("Top queries as of @ %s\n\n", time.Now().Format(time.Stamp))
						for _, q := range topItems.PrettyTopQueries {
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
						topItems.update()
						fmt.Printf("Top ads as of @ %s\n\n", time.Now().Format(time.Stamp))
						for _, q := range topItems.PrettyTopAds {
							fmt.Println(q)
						}
						return nil
					},
				},
				{
					Name:    "latest-queries",
					Aliases: []string{"lq"},
					Usage:   "Extract the latest x queries. Takes a flag for -q, the number of queries to extract",
					Flags: []cli.Flag{
						&cli.Int64Flag{
							Name:    "queries",
							Aliases: []string{"q"},
							Usage:   "The number of queries to extract",
						},
					},
					Action: func(c *cli.Context) error {
						queryAmount := c.Int("queries")
						if queryAmount < 1 {
							fmt.Println("Please enter a number of queries >= 1")
							return nil
						}
						initialisePICLI()
						allQueries.AmountOfQueriesInLog = queryAmount
						allQueries.Queries = make([]Query, allQueries.AmountOfQueriesInLog)
						allQueries.update()
						for _, query := range allQueries.Table {
							fmt.Println(query)
						}
						return nil
					},
				},
			},
		},
	},

	Action: func(c *cli.Context) error {
		initialisePICLI()
		startUI()
		return nil
	},
}

// validate that the config file and API key are in place
// load the required data and settings into memory
func initialisePICLI() {
	// firstly, has a config file been created?
	if !configFileExists() {
		log.Fatal("Please configure Pi-CLI via the 'setup' command")
	}

	settings.loadFromFile()

	// retrieve the API key depending upon its storage location
	if !settings.APIKeyIsInFile() && !APIKeyIsInKeyring() {
		log.Fatal("Please configure Pi-CLI via the 'setup' command")
	} else {
		if settings.APIKeyIsInFile() {
			piCLIData.APIKey = settings.APIKey
		} else {
			piCLIData.APIKey = retrieveAPIKeyFromKeyring()
		}
	}

	piCLIData.Settings = &settings
	piCLIData.FormattedAPIAddress = generateAPIAddress()
}
