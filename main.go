package main

import (
	"database/sql"
	"eshelon/work_DB"
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	MainUrl    = "https://oidref.com"
	DbPath     = "default.db"
	TotalNodes = 0
)

func main() {

	flag.StringVar(&DbPath, "output", "default.db", "Path to the database file")
	flag.Parse()

	database, err := sql.Open("sqlite3", DbPath)
	handleError(err)

	defer database.Close()

	err = work_DB.CreateDBrootNode(database)
	handleError(err)

	err = work_DB.CreateDBChildNodes(database)
	handleError(err)

	err = work_DB.CreateProgressTable(database)
	handleError(err)

	resumeParsing(database)

	ProcessRootNodes(MainUrl, database)
}

func resumeParsing(db *sql.DB) {
	err := db.QueryRow("SELECT COUNT(node_id) FROM parsing_progress WHERE processed = ?", true).Scan(&TotalNodes)

	handleError(err)

	rows, err := db.Query("SELECT node_id FROM parsing_progress WHERE processed = ?", false)
	handleError(err)
	defer rows.Close()

	for rows.Next() {
		var nodeURL string
		if err := rows.Scan(&nodeURL); err != nil {
			log.Fatal(err)
		}
		ProcessChildNodes(nodeURL, db)
		work_DB.MarkNodeAsProcessed(nodeURL, db)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
