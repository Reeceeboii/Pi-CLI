package cli

import (
	"fmt"
	"github.com/Reeceeboii/Pi-CLI/pkg/auth"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/logger"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

/*
	Searches for and deletes:
		- the API key from the system keyring (if exists)
		- the config file from the user's home directory (if exists)
*/
func ConfigDeleteCommand(c *cli.Context) error {
	logger.LivePiCLILogger.LogCommand("ConfigDeleteCommand")

	if auth.DeleteAPIKeyFromKeyring() {
		color.Green("System keyring API entry has been deleted!")
		logger.LivePiCLILogger.LogInformation("System keyring API entry has been deleted!")
	} else {
		color.Yellow("Pi-CLI did not find a keyring entry to delete")
		logger.LivePiCLILogger.LogInformation("Pi-CLI did not find a keyring entry to delete")
	}

	if data.DeleteConfigFile(data.GetConfigFileLocation()) {
		color.Green("Stored config file has been deleted!")
		logger.LivePiCLILogger.LogInformation("Stored config file has been deleted!")
	} else {
		color.Yellow("Pi-CLI did not find a config file to delete")
		logger.LivePiCLILogger.LogInformation("Pi-CLI did not find a config file to delete")
	}

	return nil
}

/*
	Displays any saved configuration data to the user.
	If a config file is present, that can be loaded and displayed,
	otherwise, the user can be prompted to create one.
*/
func ConfigViewCommand(c *cli.Context) error {
	logger.LivePiCLILogger.LogCommand("ConfigViewCommand")

	/*
		- Pi-Hole IP address
		- Pi-Hole port
		- Data refresh rate
	*/
	if data.ConfigFileExists(data.GetConfigFileLocation()) {
		// Display the location of the config file in the filesystem
		color.Green("Config location: %s\n", data.GetConfigFileLocation())

		// Open the config file so we can extract data from it
		data.PICLISettings.LoadFromFile(data.GetConfigFileLocation())
		fmt.Printf("Pi-Hole address: %s\n", data.PICLISettings.PiHoleAddress)
		fmt.Printf("Pi-Hole port: %d\n", data.PICLISettings.PiHolePort)
		fmt.Printf("Refresh rate: %ds\n", data.PICLISettings.RefreshS)

		// display the auto update check setting
		if data.PICLISettings.AutoCheckForUpdates {
			fmt.Println("Automatically check for updates: true")
		} else {
			fmt.Println("Automatically check for updates: false")
		}
	} else {
		color.Yellow("No config file is present - run the setup command to create one")
	}

	// and the same with the API key
	if auth.APIKeyIsInKeyring() {
		fmt.Printf("API key (keyring): %s\n", auth.RetrieveAPIKeyFromKeyring())
	} else if data.PICLISettings.APIKeyIsInFile() {
		fmt.Printf("API key (config file): %s\n", data.PICLISettings.APIKey)
	} else {
		color.Yellow("No API key has been provided - run the setup command to enter it")
	}

	return nil
}
