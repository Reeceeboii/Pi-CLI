package database

import (
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
)

/*
	Extracts the the top queries of all time. This will include both blocked and non
	blocked queries. The only factor in ordering/appearance is the number of times that
	a query for that domain has hit the Pi-Hole.

	This query is parameterised on the limit, the user can choose how many top queries
	they want returned (i.e. top 10, top 20 etc...).

	This database dump includes:
		- The domain
		- The number of queries that have been sent for that domain
		- A total sum of all of the occurrences
*/
func TopQueries(db *sql.DB, limit int64) {
	rows, err := db.Query(`
		SELECT domain, COUNT(domain)
		FROM queries
		GROUP BY domain
		ORDER BY COUNT(domain) DESC
		LIMIT ?
	`, limit)

	if err != nil {
		log.Fatalf("Error in database top queries query: %s", err.Error())
	}

	var domain string
	var occurrence int

	var occurrenceSum uint64

	tabWriter := NewConfiguredTabWriter(1)
	localisedNumberWriter := message.NewPrinter(language.English)

	// insert column headers
	_, _ = fmt.Fprintln(tabWriter, "Domain\t", "Occurrences\t")
	// insert blank line separator
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t")

	for rows.Next() {
		_ = rows.Scan(&domain, &occurrence)

		occurrenceSum = occurrenceSum + uint64(occurrence)

		_, _ = fmt.Fprintln(
			tabWriter,
			fmt.Sprintf("%s\t", domain),
			fmt.Sprintf("%s\t", localisedNumberWriter.Sprintf("%d", occurrence)))
	}

	// insert blank line separator
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t")
	// insert column headers
	_, _ = fmt.Fprintln(tabWriter, "\t", "Total\t")

	// insert the total of the occurrences
	_, _ = fmt.Fprintln(
		tabWriter,
		"\t",
		fmt.Sprintf("%s\t", localisedNumberWriter.Sprintf("%d", occurrenceSum)))

	if err := tabWriter.Flush(); err != nil {
		return
	}
}
