package database

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strings"
)

/*
	Extracts a summary of all of the clients that have been served by the Pi-Hole instance.

	NOTE:
	This will include duplicate DNS entries if the same client has been seen multiple times
	(i.e. a phone being seen as a LAN client both locally and via a VPN - the client itself
	is the same but has separate query counts and different local addresses and as such will
	be listed as many times as it has appeared in different contexts).

	This database dump includes:
		- The client's address: (IP or mac addr)
		- The date that the client was first seen
		- The date that the last query from the client was received
		- The total number of queries received from the client
		- The client's DNS name
*/
func ClientSummary(db *sql.DB) {
	rows, err := db.Query(`
		SELECT DISTINCT n.hwaddr, n.firstSeen, n.lastQuery, n.numQueries, na.name
		FROM network n
		INNER JOIN network_addresses na on n.id = na.network_id
		WHERE n.numQueries != 0
		ORDER BY numQueries DESC
	`)

	if err != nil {
		log.Fatalf("Error in database client summary query: %s", err.Error())
	}

	var address string
	var firstSeen int
	var lastQuery int
	var numQueries int
	var name string

	tabWriter := NewConfiguredTabWriter(1)
	localisedNumberWriter := message.NewPrinter(language.English)

	// insert column headers
	_, _ = fmt.Fprintln(
		tabWriter,
		"#\t",
		"Address\t",
		"First seen\t",
		"Last query\t",
		"No. queries\t",
		"DNS\t")

	// insert blank line separator
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t", "\t", "\t", "\t", "\t")

	row := 1

	// print out each row from the query results
	for rows.Next() {
		_ = rows.Scan(&address, &firstSeen, &lastQuery, &numQueries, &name)

		// if the string is denoting an IP, we can chop off the IP identifier from the row entry
		if strings.Contains(address, "ip-") {
			address = strings.Split(address, "ip-")[1]
		}

		_, _ = fmt.Fprintln(
			tabWriter,
			fmt.Sprintf("%d\t", row),
			fmt.Sprintf("%s\t", address),
			fmt.Sprintf("%s\t", FormattedDBUnixTimestamp(firstSeen)),
			fmt.Sprintf("%s\t", FormattedDBUnixTimestamp(lastQuery)),
			fmt.Sprintf("%s\t", localisedNumberWriter.Sprintf("%d", numQueries)),
			fmt.Sprintf("%s\t", name))
		row++
	}

	// if the row counter has never been incremented, the database query returned zero results
	if row == 1 {
		color.Red("0 results in database")
	}

	if err := tabWriter.Flush(); err != nil {
		return
	}
}
