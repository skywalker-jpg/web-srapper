package work_DB

import (
	"database/sql"
	"log"
)

func InsertIntoRoot_nodes(node_name, subChildren, subNodesTotal, description, information string, db *sql.DB) {
	insertRootNodeStmt, err := db.Prepare(`
        INSERT INTO root_nodes (node_name, subChildren, subNodesTotal, description, information)
        VALUES (?, ?, ?, ?, ?)
    `)

	handleError(err)

	insertRootNodeStmt.Exec(node_name, subChildren, subNodesTotal, description, information)
	insertRootNodeStmt.Close()
}

func InsertIntoChild_nodes(node_name, sub_children, sub_nodes_total string, db *sql.DB) {
	insertChildNodeStmt, err := db.Prepare(`
		INSERT INTO child_nodes (node_name, sub_children, sub_nodes_total)
		VALUES (?, ?, ?)
	`)

	handleError(err)

	insertChildNodeStmt.Exec(node_name, sub_children, sub_nodes_total)
	insertChildNodeStmt.Close()
}

func IsNodeProcessed(nodeURL string, db *sql.DB) bool {
	var processed bool
	err := db.QueryRow("SELECT processed FROM parsing_progress WHERE node_id = ?", nodeURL).Scan(&processed)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return processed
}

func MarkNodeAsProcessed(nodeURL string, db *sql.DB) {
	_, err := db.Exec("INSERT INTO parsing_progress (node_id, processed) VALUES (?, ?)", nodeURL, true)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
