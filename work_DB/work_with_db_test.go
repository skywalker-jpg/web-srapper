package work_DB

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB создает временную базу данных для тестов и возвращает на нее ссылку.
// Если произойдет ошибка, функция завершит выполнение теста.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// TestInsertIntoRoot_nodes проверяет функцию InsertIntoRoot_nodes на вставку данных в таблицу root_nodes.
func TestInsertIntoRoot_nodes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создание таблицы root_nodes в тестовой базе данных.
	CreateDBrootNode(db)

	// Вызов функции InsertIntoRoot_nodes для вставки данных.
	InsertIntoRoot_nodes("Node1", "SubChildren1", "SubNodesTotal1", "Description1", "Information1", db)

	// Проверка, что данные были успешно вставлены в таблицу root_nodes.
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM root_nodes WHERE node_name = ?", "Node1").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row in root_nodes, but got %d", count)
	}
}

// TestInsertIntoChild_nodes проверяет функцию InsertIntoChild_nodes на вставку данных в таблицу child_nodes.
func TestInsertIntoChild_nodes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создание таблицы child_nodes в тестовой базе данных.
	CreateDBChildNodes(db)

	// Вызов функции InsertIntoChild_nodes для вставки данных.
	InsertIntoChild_nodes("ChildNode1", "SubChildren1", "SubNodesTotal1", db)

	// Проверка, что данные были успешно вставлены в таблицу child_nodes.
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM child_nodes WHERE node_name = ?", "ChildNode1").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row in child_nodes, but got %d", count)
	}
}

// TestIsNodeProcessed проверяет функцию IsNodeProcessed на проверку статуса обработки узла.
func TestIsNodeProcessed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создание таблицы parsing_progress в тестовой базе данных.
	CreateProgressTable(db)

	// Вставка записи о обработанном узле.
	_, err := db.Exec("INSERT INTO parsing_progress (node_id, processed) VALUES (?, ?)", "NodeURL1", true)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка, что узел NodeURL1 был обработан.
	processed := IsNodeProcessed("NodeURL1", db)
	if !processed {
		t.Errorf("Expected NodeURL1 to be processed, but it wasn't")
	}

	// Проверка, что узел NodeURL2 не был обработан.
	notProcessed := IsNodeProcessed("NodeURL2", db)
	if notProcessed {
		t.Errorf("Expected NodeURL2 to be not processed, but it was")
	}
}

// TestMarkNodeAsProcessed проверяет функцию MarkNodeAsProcessed на установку статуса обработки узла.
func TestMarkNodeAsProcessed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Создание таблицы parsing_progress в тестовой базе данных.
	CreateProgressTable(db)

	// Установка статуса обработки для узла NodeURL2.
	MarkNodeAsProcessed("NodeURL2", db)

	// Проверка, что узел NodeURL2 был успешно помечен как обработанный.
	var processed bool
	err := db.QueryRow("SELECT processed FROM parsing_progress WHERE node_id = ?", "NodeURL2").Scan(&processed)
	if err != nil {
		t.Fatal(err)
	}
	if !processed {
		t.Errorf("Expected NodeURL2 to be marked as processed, but it wasn't")
	}
}
