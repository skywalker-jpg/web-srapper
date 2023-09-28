package work_DB

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDatabase(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening test database: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func TestCreateDBrootNode(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	err := CreateDBrootNode(db)
	if err != nil {
		t.Fatalf("Error creating root_nodes table: %v", err)
	}

	_, err = db.Exec("SELECT * FROM root_nodes")
	if err != nil {
		t.Fatalf("Error querying root_nodes table: %v", err)
	}
}

func TestCreateDBChildNodes(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	err := CreateDBChildNodes(db)
	if err != nil {
		t.Fatalf("Error creating child_nodes table: %v", err)
	}

	_, err = db.Exec("SELECT * FROM child_nodes")
	if err != nil {
		t.Fatalf("Error querying child_nodes table: %v", err)
	}
}

func TestCreateProgressTable(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	err := CreateProgressTable(db)
	if err != nil {
		t.Fatalf("Error creating parsing_progress table: %v", err)
	}

	_, err = db.Exec("SELECT * FROM parsing_progress")
	if err != nil {
		t.Fatalf("Error querying parsing_progress table: %v", err)
	}
}
