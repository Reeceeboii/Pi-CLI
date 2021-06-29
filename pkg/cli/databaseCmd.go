package cli

import (
	"github.com/Reeceeboii/Pi-CLI/pkg/database"
	"github.com/Reeceeboii/Pi-CLI/pkg/logger"
	"github.com/urfave/cli/v2"
)

/*
	FOR ALL DATABASE COMMANDS:
		If no path is provided by the user, Pi-CLI will assume that the database file's
		name hasn't been changed from it's default name, and that is has been placed in the
		same working directory that it is being executed from. This saves some command typing.
*/

/*
	Extracts a summary of data regarding the Pi-Hole's clients
*/
func RunDatabaseClientSummaryCommand(c *cli.Context) error {
	logger.LivePiCLILogger.LogCommand("RunDatabaseClientSummaryCommand")

	path := c.String("path")
	if path == "" {
		path = database.DefaultDatabaseFileLocation
		logger.LivePiCLILogger.LogCommand("DB path not provided, reverting to default")
	}

	logger.LivePiCLILogger.LogCommand("DB path: " + path)

	conn := database.Connect(path)
	database.ClientSummary(conn)

	return nil
}

/*
	Extracts all time top query data from the database file.
*/
func RunDatabaseTopQueriesCommand(c *cli.Context) error {
	logger.LivePiCLILogger.LogCommand("RunDatabaseTopQueriesCommand")

	path := c.String("path")
	if path == "" {
		path = database.DefaultDatabaseFileLocation
		logger.LivePiCLILogger.LogCommand("DB path not provided, reverting to default")
	}

	logger.LivePiCLILogger.LogCommand("DB path: " + path)

	conn := database.Connect(path)
	database.TopQueries(conn, c.Int64("limit"), c.String("filter"))

	return nil
}
