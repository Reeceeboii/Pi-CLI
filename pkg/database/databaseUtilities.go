package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

// Database constants
const (
	// The name of the database driver to use
	DBDriverName = "sqlite3"
	// The default limit on the number of queries returned from some database queries
	DefaultQueryTableLimit = 10
	/*
		The default location and name used by database commands to look for the Pi-Hole's
		FTL database file
	*/
	DefaultDatabaseFileLocation = "./pihole-FTL.db"
)

/*
	Attempts to connect to a database, returns an sql.DB handle
	if this connection succeeds
*/
func Connect(pathToPotentialDB string) *sql.DB {
	conn := &sql.DB{}
	if validateDatabase(pathToPotentialDB) {
		conn, _ = sql.Open(DBDriverName, pathToPotentialDB)
	}
	return conn
}

/*
	Returns a RFC822 formatted version of a given Unix time integer retrieved
	from a database row. For example, given the Unix time of 1612548060, the function
	will return the string "05 Feb 21 18:01 GMT"
*/
func FormattedDBUnixTimestamp(stamp int) string {
	return time.Unix(int64(stamp), 0).Format(time.RFC822)
}

/*
	Returns a newly configured tabwriter.Writer, with a parameterised padding,
	allowing optional changes to the padding between an element and the edge of
	its cell
*/
func NewConfiguredTabWriter(padding int) *tabwriter.Writer {
	return tabwriter.NewWriter(
		os.Stdout,
		0,
		0,
		padding,
		' ',
		tabwriter.Debug)
}

// Checks if the filepath to a database is valid and that a connection can be opened
func validateDatabase(pathToPotentialDB string) bool {
	if err := doesDatabaseFileExist(pathToPotentialDB); err != nil {
		log.Fatal(err.Error())
	}

	if err := canOpenConnectionToDB(pathToPotentialDB); err != nil {
		log.Fatal(err.Error())
	}

	return true
}

// Does the database file exist?
func doesDatabaseFileExist(pathToPotentialDB string) error {
	if _, err := os.Stat(pathToPotentialDB); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("'%s' does not exist", pathToPotentialDB))
	}
	return nil
}

// Can a connection be opened with the DB file?
func canOpenConnectionToDB(pathToPotentialDB string) error {
	// attempt to open a connection to the database
	conn, err := sql.Open(DBDriverName, pathToPotentialDB)

	/*
		 	If the connection failed, return an error. If we get to this point, the file is valid
			and is present in the local filesystem. However, either the file is not an SQLite database,
			or it is somehow unreadable.
	*/
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"Failed to connect. Check that the path is correct & points to a valid file: %s", err.Error()))
	}

	if err := conn.Close(); err != nil {
		return errors.New(fmt.Sprintf("Failed to close connection: %s", err.Error()))
	}

	return nil
}
