package work_DB

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInsertIntoRoot_nodes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	CreateDBrootNode(db)
	InsertIntoRoot_nodes("Node1", "SubChildren1", "SubNodesTotal1", "Description1", "Information1", db)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM root_nodes WHERE node_name = ?", "Node1").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row in root_nodes, but got %d", count)
	}
}

func TestInsertIntoChild_nodes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	CreateDBChildNodes(db)
	InsertIntoChild_nodes("ChildNode1", "SubChildren1", "SubNodesTotal1", db)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM child_nodes WHERE node_name = ?", "ChildNode1").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row in child_nodes, but got %d", count)
	}
}

func TestIsNodeProcessed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	CreateProgressTable(db)

	_, err := db.Exec("INSERT INTO parsing_progress (node_id, processed) VALUES (?, ?)", "NodeURL1", true)
	if err != nil {
		t.Fatal(err)
	}

	processed := IsNodeProcessed("NodeURL1", db)
	if !processed {
		t.Errorf("Expected NodeURL1 to be processed, but it wasn't")
	}

	notProcessed := IsNodeProcessed("NodeURL2", db)
	if notProcessed {
		t.Errorf("Expected NodeURL2 to be not processed, but it was")
	}
}

func TestMarkNodeAsProcessed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	CreateProgressTable(db)
	MarkNodeAsProcessed("NodeURL2", db)

	var processed bool
	err := db.QueryRow("SELECT processed FROM parsing_progress WHERE node_id = ?", "NodeURL2").Scan(&processed)
	if err != nil {
		t.Fatal(err)
	}
	if !processed {
		t.Errorf("Expected NodeURL2 to be marked as processed, but it wasn't")
	}
}
