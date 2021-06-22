package cli

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/auth"
	"github.com/Reeceeboii/Pi-CLI/pkg/data"
	"github.com/Reeceeboii/Pi-CLI/pkg/network"
	"github.com/fatih/color"
	"os"
)

/*
	Validate that the config file and API key are in place.
	Load the required settings into memory
*/
func InitialisePICLI() {
	// firstly, has a config file been created?
	if !data.ConfigFileExists() {
		color.Red("Please configure Pi-CLI via the 'setup' command")
		os.Exit(1)
	}

	data.PICLISettings.LoadFromFile()

	// retrieve the API key depending upon its storage location
	if !data.PICLISettings.APIKeyIsInFile() && !auth.APIKeyIsInKeyring() {
		color.Red("Please configure Pi-CLI via the 'setup' command")
		os.Exit(1)
	} else {
		if data.PICLISettings.APIKeyIsInFile() {
			data.LivePiCLIData.APIKey = data.PICLISettings.APIKey
		} else {
			data.LivePiCLIData.APIKey = auth.RetrieveAPIKeyFromKeyring()
		}
	}

	data.LivePiCLIData.Settings = data.PICLISettings
	data.LivePiCLIData.FormattedAPIAddress = network.GenerateAPIAddress(
		data.PICLISettings.PiHoleAddress,
		data.PICLISettings.PiHolePort)
}
