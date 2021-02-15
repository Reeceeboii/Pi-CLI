package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net"
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

				// read in the data refresh rate
				fmt.Print("Please enter your preferred data refresh rate in seconds (default 1s): ")
				refreshS, _ := reader.ReadString('\n')
				trimmed = strings.TrimSpace(refreshS)
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

				fmt.Print("Enter API key (stored securely, not in a file): ")
				apiKey, _ := reader.ReadString('\n')
				apiKey = strings.TrimSpace(apiKey)
				if len(apiKey) < 1 {
					log.Fatal("Please provide your API key for authentication")
				}
				storeAPIKey(&apiKey)
				// write config to disk
				settings.saveToFile()
				return nil
			},
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "View currently saved configuration data for verification",
			Action: func(c *cli.Context) error {
				// if the config file is present, that can be loaded and displayed
				if configFileExists() {
					configFileData := Settings{}
					configFileData.loadFromFile()
					fmt.Printf("%s%s\n", "Pi-Hole address: ", configFileData.PiHoleAddress)
					fmt.Printf("%s%d\n", "Pi-Hole port: ", configFileData.PiHolePort)
					fmt.Printf("%s%d%s\n", "Refresh rate: ", configFileData.RefreshS, "s")
				} else {
					fmt.Println("No config file is present - run the setup command to create one")
				}

				// and the same with the API key
				if APIKeyExists() {
					fmt.Printf("%s%s\n", "API key: ", retrieveAPIKey())
				} else {
					fmt.Println("No API key has been provided - run the setup command to enter it")
				}
				return nil
			},
		},
	},

	Action: func(c *cli.Context) error {
		if !configFileExists() || !APIKeyExists() {
			log.Fatal("Please configure Pi-CLI via the 'setup' command")
		}

		settings.loadFromFile()
		piCLIData.Settings = &settings
		piCLIData.APIKey = retrieveAPIKey()
		piCLIData.FormattedAPIAddress = generateAPIAddress()
		startUI()
		return nil
	},
}
