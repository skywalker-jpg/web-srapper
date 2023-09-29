package work_DB

import (
	"database/sql"
	"log"
)

// InsertIntoRootNodes выполняет вставку данных о корневых узлах в таблицу root_nodes базы данных.
func InsertIntoRoot_nodes(node_name, subChildren, subNodesTotal, description, information string, db *sql.DB) {
	// Подготовка SQL-запроса для вставки данных.
	insertRootNodeStmt, err := db.Prepare(`
        INSERT INTO root_nodes (node_name, subChildren, subNodesTotal, description, information)
        VALUES (?, ?, ?, ?, ?)
    `)

	handleError(err)

	// Выполнение вставки данных в таблицу root_nodes.
	_, err = insertRootNodeStmt.Exec(node_name, subChildren, subNodesTotal, description, information)
	handleError(err)

	// Закрытие подготовленного выражения после использования.
	insertRootNodeStmt.Close()
}

// InsertIntoChildNodes выполняет вставку данных о дочерних узлах в таблицу child_nodes базы данных.
func InsertIntoChild_nodes(node_name, sub_children, sub_nodes_total string, db *sql.DB) {
	// Подготовка SQL-запроса для вставки данных.
	insertChildNodeStmt, err := db.Prepare(`
		INSERT INTO child_nodes (node_name, sub_children, sub_nodes_total)
		VALUES (?, ?, ?)
	`)

	handleError(err)

	// Выполнение вставки данных в таблицу child_nodes.
	_, err = insertChildNodeStmt.Exec(node_name, sub_children, sub_nodes_total)
	handleError(err)

	// Закрытие подготовленного выражения после использования.
	insertChildNodeStmt.Close()
}

// IsNodeProcessed проверяет, был ли узел обработан, основываясь на его URL.
// Возвращает true, если узел был обработан, иначе возвращает false.
func IsNodeProcessed(nodeURL string, db *sql.DB) bool {
	var processed bool
	err := db.QueryRow("SELECT processed FROM parsing_progress WHERE node_id = ?", nodeURL).Scan(&processed)
	if err == sql.ErrNoRows {
		// Если нет строк с данным URL, считаем, что узел не был обработан.
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return processed
}

// MarkNodeAsProcessed устанавливает статус обработки узла с указанным URL в true.
func MarkNodeAsProcessed(nodeURL string, db *sql.DB) {
	_, err := db.Exec("INSERT INTO parsing_progress (node_id, processed) VALUES (?, ?)", nodeURL, true)
	handleError(err)
}

// handleError обрабатывает ошибку, выводя ее и завершая выполнение программы при наличии ошибки.
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
