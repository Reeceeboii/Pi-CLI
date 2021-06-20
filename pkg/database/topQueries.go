package database

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
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
func TopQueries(db *sql.DB, limit int64, domainFilter string) {
	color.Yellow("\nLimit: %d | Filter: '%s' \n\n", limit, domainFilter)

	var rows *sql.Rows
	var err error

	// if filter has been provided, we want to plug it into the SQL query
	if domainFilter == "" {
		rows, err = db.Query(`
		SELECT domain, COUNT(domain)
		FROM queries
		GROUP BY domain
		ORDER BY COUNT(domain) DESC
		LIMIT ?
	`, limit)
	} else {
		sqlFilter := "%" + domainFilter + "%"

		rows, err = db.Query(`
		SELECT domain, COUNT(domain)
		FROM queries
		WHERE queries.domain LIKE ?
		GROUP BY domain
		ORDER BY COUNT(domain) DESC
		LIMIT ?
	`, sqlFilter, limit)
	}

	if err != nil {
		log.Fatalf("Error in database top queries query: %s", err.Error())
	}

	var domain string
	var occurrence int

	var occurrenceSum uint64

	tabWriter := NewConfiguredTabWriter(1)
	localisedNumberWriter := message.NewPrinter(language.English)

	// insert column headers
	_, _ = fmt.Fprintln(tabWriter, "#\t", "Domain\t", "Occurrences\t")
	// insert blank line separator
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t", "\t")

	// used to count the rows as they're outputted
	var row int64 = 1

	for rows.Next() {
		_ = rows.Scan(&domain, &occurrence)

		occurrenceSum = occurrenceSum + uint64(occurrence)

		_, _ = fmt.Fprintln(
			tabWriter,
			fmt.Sprintf("%d\t", row),
			fmt.Sprintf("%s\t", domain),
			localisedNumberWriter.Sprintf("%d\t", occurrence),
		)
		row++
	}

	// insert blank line separator
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t", "\t")
	// insert column headers
	_, _ = fmt.Fprintln(tabWriter, "\t", "\t", "Total\t")

	// insert the total of the occurrences
	_, _ = fmt.Fprintln(
		tabWriter,
		"\t",
		"\t",
		fmt.Sprintf("%s\t", localisedNumberWriter.Sprintf("%d", occurrenceSum)))

	// if the row counter has never been incremented, the database query returned zero results
	if row == 1 {
		color.Red("0 results in database")
	}

	if err := tabWriter.Flush(); err != nil {
		return
	}
}
