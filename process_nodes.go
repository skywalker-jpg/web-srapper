package main

import (
	"database/sql"
	"eshelon/work_DB"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ProcessRootNodes обрабатывает корневые узлы, извлекая информацию о них и их дочерних узлах.
func ProcessRootNodes(url string, db *sql.DB) {
	// Выполнение HTTP GET-запроса с возможностью повторных попыток.
	res, err := retryHttpGet(url, 4)
	handleError(err)
	defer res.Body.Close()

	// Обработка статуса HTTP-ответа.
	handleHTTP(res.StatusCode)

	// Создание объекта goquery для разбора HTML-страницы.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	handleError(err)

	// Поиск таблицы с информацией о корневых узлах.
	rootNodeTable := doc.Find("h2:contains('Root Tree Nodes')").Next()

	var childNodesURLs []string

	// Итерация по строкам таблицы с корневыми узлами.
	rootNodeTable.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
		// Извлечение данных о корневых узлах.
		nodeURL, _ := rowHtml.Find("td a").Attr("href")
		nodeName := rowHtml.Find("td:nth-child(2)").Text()
		subChildren := rowHtml.Find("td:nth-child(3)").Text()
		subNodesTotal := rowHtml.Find("td:nth-child(4)").Text()
		description := rowHtml.Find("td:nth-child(5)").Text()
		information := rowHtml.Find("td:nth-child(6)").Text()

		// Вставка данных о корневом узле в базу данных.
		work_DB.InsertIntoRoot_nodes(nodeName, subChildren, subNodesTotal, description, information, db)

		// Если у узла есть URL, добавляем его в список дочерних узлов для обработки.
		if nodeURL != "" {
			nodeURL = MainUrl + nodeURL
			childNodesURLs = append(childNodesURLs, nodeURL)
		}
	})

	// Обработка дочерних узлов корневых узлов.
	WorkChildNodes(childNodesURLs, db)
}

// ProcessChildNodes обрабатывает дочерние узлы, извлекая информацию о них.
func ProcessChildNodes(url string, db *sql.DB) {
	var childNodesURLs []string

	// Выполнение HTTP GET-запроса с возможностью повторных попыток.
	res, err := retryHttpGet(url, 4)
	handleError(err)
	defer res.Body.Close()

	// Обработка статуса HTTP-ответа.
	handleHTTP(res.StatusCode)

	// Создание объекта goquery для разбора HTML-страницы.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	handleError(err)

	// Поиск таблицы с информацией о дочерних узлах.
	childNodeTable := doc.Find("h3:contains('Children')").Next()

	// Итерация по строкам таблицы с дочерними узлами.
	childNodeTable.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
		// Извлечение данных о дочерних узлах.
		nodeURL, _ := rowHtml.Find("td a").Attr("href")
		nodeName := rowHtml.Find("td a").Text()
		subChildren := rowHtml.Find("td:nth-child(3)").Text()
		subNodesTotal := rowHtml.Find("td:nth-child(4)").Text()

		// Вставка данных о дочернем узле в базу данных.
		work_DB.InsertIntoChild_nodes(nodeName, subChildren, subNodesTotal, db)

		// Увеличение счетчика обработанных узлов и вывод статистики.
		TotalNodes++
		if TotalNodes%100 == 0 {
			fmt.Printf("Processed nodes: %d\n", TotalNodes)
		}

		// Если у дочернего узла есть URL, добавляем его в список для обработки.
		if nodeURL != "" {
			nodeURL = MainUrl + nodeURL
			childNodesURLs = append(childNodesURLs, nodeURL)
		}
	})

	// Если есть дочерние узлы, продолжаем обработку.
	if len(childNodesURLs) != 0 {
		WorkChildNodes(childNodesURLs, db)
	}
}

// WorkChildNodes обрабатывает дочерние узлы, начиная с указанных URL.
func WorkChildNodes(childNodes []string, db *sql.DB) {
	for _, childURL := range childNodes {
		// Проверяем, не был ли узел уже обработан.
		if !work_DB.IsNodeProcessed(childURL, db) {
			// Если узел не обработан, выполняем его обработку и помечаем как обработанный.
			ProcessChildNodes(childURL, db)
			work_DB.MarkNodeAsProcessed(childURL, db)
		}
	}
}

// handleHTTP проверяет статус HTTP-ответа и завершает программу с ошибкой, если статус не равен 200.
func handleHTTP(StatusCode int) {
	if StatusCode != 200 {
		log.Fatalf("Request failed with status code %d", StatusCode)
	}
}

// retryHttpGet выполняет HTTP GET-запрос с возможностью повторных попыток в случае ошибки.
func retryHttpGet(url string, maxRetries int) (*http.Response, error) {
	for i := 0; i < maxRetries; i++ {
		res, err := http.Get(url)
		if err == nil && res.StatusCode == 200 {
			return res, nil
		}
		fmt.Printf("Attempt %d: Error fetching %s: %v\n", i+1, url, err)
		time.Sleep(time.Second * 60)
	}

	return nil, fmt.Errorf("failed to fetch %s after %d attempts", url, maxRetries)
}
