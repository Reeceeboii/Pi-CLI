package cli

import (
	"bufio"
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/auth"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
Reads in data from the user and uses it to construct a config file that Pi-CLI can use
in the future.

The setup commands takes:
  - The IP address of the target Pi-Hole instance
  - The port exposing the Pi-Hole's web interface
  - A data refresh rate in seconds
  - User's Pi-Hole API key (used for authentication)

It will then ask them if they wish to store the API key in their system keyring or the config
file itself.
*/
func SetupCommand(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		// read in the IP address and check that it is valid
		fmt.Print(" > Please enter the IP address of your Pi-Hole: ")
		piHoleAddress, _ := reader.ReadString('\n')
		ip := net.ParseIP(strings.TrimSpace(piHoleAddress))
		if ip == nil {
			color.Yellow("Please enter a valid IP address")
			continue
		}
		data.PICLISettings.PiHoleAddress = ip.String()
		break
	}

	for {
		// read in the port
		fmt.Print(" > Please enter the port that exposes the web interface (default 80): ")
		piHolePort, _ := reader.ReadString('\n')
		trimmed := strings.TrimSpace(piHolePort)
		// if the user entered something, validate and apply. Else, revert to the default
		if len(trimmed) > 0 {
			intPiHolePort, err := strconv.Atoi(trimmed)
			if err != nil {
				color.Yellow("Please enter a number")
				continue
			}
			if intPiHolePort < 1 || intPiHolePort > 65535 {
				color.Yellow("Port must be between 1 and 65535")
				continue
			}
			testAddressWithPort := network.GenerateAPIAddress(data.PICLISettings.PiHoleAddress, intPiHolePort)
			if network.IsAlive(testAddressWithPort) {
				data.PICLISettings.PiHolePort = intPiHolePort
			} else {
				continue
			}
		} else {
			testAddressWithPort := network.GenerateAPIAddress(data.PICLISettings.PiHoleAddress, data.DefaultPort)
			if network.IsAlive(testAddressWithPort) {
				data.PICLISettings.PiHolePort = data.DefaultPort
			} else {
				continue
			}
		}

		// send a request to the PiHole to validate that the IP and port actually point to it
		tempURL := fmt.Sprintf(
			"http://%s:%d/admin/api.php",
			data.PICLISettings.PiHoleAddress,
			data.PICLISettings.PiHolePort)
		req, err := http.NewRequest("GET", tempURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		res, err := network.HttpClient.Do(req)

		// if the details are valid and the request didn't time out...
		// lazy evaluation saves us from deref errors here and saves a check
		if err == nil && network.ValidatePiHoleDetails(res) {
			break
		} else {
			color.Yellow("Pi-Hole doesn't seem to be alive, check your details and try again!")
			fmt.Println()
		}
	}

	color.Green(
		"Pi-Hole reachable at %s:%d!\n",
		data.PICLISettings.PiHoleAddress,
		data.PICLISettings.PiHolePort)

	// read in the data refresh rate
	for {
		fmt.Print(" > Please enter your preferred data refresh rate in seconds (default 1s): ")
		refreshS, _ := reader.ReadString('\n')
		trimmed := strings.TrimSpace(refreshS)
		if len(trimmed) > 0 {
			intRefreshS, err := strconv.Atoi(trimmed)
			if err != nil {
				color.Yellow("Please enter a number")
				continue
			}
			if intRefreshS < 1 {
				color.Yellow("Refresh time cannot be less than 1 second")
				continue
			}
			data.PICLISettings.RefreshS = intRefreshS
			break
		} else {
			break
		}
	}

	// read in the API key and work out where the user wants to store it (keyring or config file)
	for {
		fmt.Print(" > Please enter your Pi-Hole API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)
		if len(apiKey) < 1 {
			color.Yellow("Please provide your API key for authentication")
			continue
		}

		data.PICLISettings.APIKey = apiKey

		// before we store the API token (keyring or config file), we should check that it's valid
		// the address + port have been validated by this point so we're safe to shoot requests at it
		data.LivePiCLIData.Settings = data.PICLISettings
		data.LivePiCLIData.FormattedAPIAddress = network.GenerateAPIAddress(
			data.PICLISettings.PiHoleAddress,
			data.PICLISettings.PiHolePort)

		if !auth.ValidateAPIKey(data.LivePiCLIData.FormattedAPIAddress, data.PICLISettings.APIKey) {
			color.Yellow("That API token doesn't seem to be correct, check it and try again!")
		} else {
			break
		}
	}

	color.Green("Authenticated with API key!\n")

	fmt.Print(" > Do you wish to store the API key in your system keyring? (y/n - default y): ")
	storageChoice, _ := reader.ReadString('\n')
	storageChoice = strings.ToLower(strings.TrimSpace(storageChoice))

	// if they wish to use their system's keyring...
	if storageChoice == "y" || len(storageChoice) == 0 {
		err := auth.StoreAPIKeyInKeyring(data.PICLISettings.APIKey)

		if err == nil {
			color.Green("Your API token has been securely stored in your system keyring")
			/*
				After the API key has been saved to the keyring, there is no longer a need to save it
				to the config file, so the stored copy of it can be removed from the in-memory settings
				instance before it gets serialised to disk
			*/
			data.PICLISettings.APIKey = ""
		} else {
			color.Yellow("System keyring call failed, falling back to config file")
		}
	}

	// write config file to disk
	// all fields in the settings struct would have been set by this point
	if err := data.PICLISettings.SaveToFile(); err != nil {
		color.Red("Failed to save settings")
		log.Fatal(err.Error())
	}

	color.Green("\nConfiguration successfully saved to %s", data.GetConfigFileLocation())
	return nil
}
