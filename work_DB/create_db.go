package work_DB

import (
	"database/sql"
)

func CreateDBrootNode(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS root_nodes (
            node_name TEXT PRIMARY KEY,
            subChildren TEXT,
            subNodesTotal TEXT,
            description TEXT,
            information TEXT
        )
    `)
	return err
}

func CreateDBChildNodes(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS child_nodes (
            node_name TEXT PRIMARY KEY,
            sub_children TEXT,
            sub_nodes_total TEXT
        )
    `)
	return err
}

func CreateProgressTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS parsing_progress (
            node_id TEXT PRIMARY KEY,
            processed BOOLEAN
        )
    `)
	return err
}
