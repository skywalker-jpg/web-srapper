package work_DB

import (
	"database/sql"
)

// CreateDBrootNode создает таблицу root_nodes в базе данных, если она не существует.
// Эта таблица предназначена для хранения информации о корневых узлах.
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

// CreateDBChildNodes создает таблицу child_nodes в базе данных, если она не существует.
// Эта таблица предназначена для хранения информации о дочерних узлах.
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

// CreateProgressTable создает таблицу parsing_progress в базе данных, если она не существует.
// Эта таблица предназначена для отслеживания прогресса обработки узлов.
func CreateProgressTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS parsing_progress (
            node_id TEXT PRIMARY KEY,
            processed BOOLEAN
        )
    `)
	return err
}
