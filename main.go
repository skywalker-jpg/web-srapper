package main

import (
	"database/sql"
	"eshelon/work_DB"
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	MainUrl    = "https://oidref.com" // Основной URL, с которого начинается парсинг.
	DbPath     = "default.db"         // Путь к файлу базы данных по умолчанию.
	TotalNodes = 0                    // Общее количество узлов, которые было обработано.
)

func main() {
	// Парсинг флагов командной строки для определения пути к базе данных.
	flag.StringVar(&DbPath, "output", "default.db", "Path to the database file")
	flag.Parse()

	// Открытие базы данных.
	database, err := sql.Open("sqlite3", DbPath)
	handleError(err)

	defer database.Close()

	// Создание таблиц в базе данных, если они не существуют.
	err = work_DB.CreateDBrootNode(database)
	handleError(err)

	err = work_DB.CreateDBChildNodes(database)
	handleError(err)

	err = work_DB.CreateProgressTable(database)
	handleError(err)

	// Возобновление парсинга узлов, которые еще не были обработаны.
	resumeParsing(database)

	// Начало обработки корневых узлов.
	ProcessRootNodes(MainUrl, database)
}

// resumeParsing продолжает обработку узлов, которые еще не были обработаны.
func resumeParsing(db *sql.DB) {
	// Получение общего количества узлов, которые уже были обработаны.
	err := db.QueryRow("SELECT COUNT(node_id) FROM parsing_progress WHERE processed = ?", true).Scan(&TotalNodes)

	handleError(err)

	// Получение списка узлов, которые еще не были обработаны.
	rows, err := db.Query("SELECT node_id FROM parsing_progress WHERE processed = ?", false)
	handleError(err)
	defer rows.Close()

	// Итерация по узлам и их обработка.
	for rows.Next() {
		var nodeURL string
		if err := rows.Scan(&nodeURL); err != nil {
			log.Fatal(err)
		}
		ProcessChildNodes(nodeURL, db)           // Обработка дочерних узлов.
		work_DB.MarkNodeAsProcessed(nodeURL, db) // Пометка узла как обработанного.
	}
}

// handleError обрабатывает ошибку, выводя ее и завершая выполнение программы при наличии ошибки.
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
