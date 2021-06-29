package main

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/cli"
	"github.com/Reeceeboii/Pi-CLI/pkg/logger"
	"github.com/Reeceeboii/Pi-CLI/pkg/update"
	"log"
	"os"
)

// these are placeholder strings, and are overridden via the makefile
var PICLIVersion = "Undefined"
var GitHash = "Undefined"

func main() {
	// pass the set variables through to the update package
	update.SetVersion(PICLIVersion)
	update.SetGitHash(GitHash)

	// start the CLI app
	if err := cli.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	logger.LivePiCLILogger.LogStatus("Pi-CLI exiting")

	// attempt to close the log file's handle
	if logger.LivePiCLILogger.Enabled {
		if err := logger.LivePiCLILogger.LogFileHandle.Close(); err != nil {
			log.Println("Unable to close logger file handle")
			log.Println(err)
		}
	}
}
