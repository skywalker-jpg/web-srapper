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

func ProcessRootNodes(url string, db *sql.DB) {
	res, err := retryHttpGet(url, 4)
	handleError(err)

	defer res.Body.Close()

	handleHTTP(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	handleError(err)

	rootNodeTable := doc.Find("h2:contains('Root Tree Nodes')").Next()

	var child_nodes []string

	rootNodeTable.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
		node_id, _ := rowHtml.Find("td a").Attr("href")
		node_name := rowHtml.Find("td:nth-child(2)").Text()
		subChildren := rowHtml.Find("td:nth-child(3)").Text()
		subNodesTotal := rowHtml.Find("td:nth-child(4)").Text()
		description := rowHtml.Find("td:nth-child(5)").Text()
		information := rowHtml.Find("td:nth-child(6)").Text()

		work_DB.InsertIntoRoot_nodes(node_name, subChildren, subNodesTotal, description, information, db)

		if node_id != "" {
			node_id = MainUrl + node_id
			child_nodes = append(child_nodes, node_id)
		}
	})

	WorkChildNodes(child_nodes, db)
}

func ProcessChildNodes(url string, db *sql.DB) {
	var child_nodes []string

	res, err := retryHttpGet(url, 4)
	handleError(err)
	defer res.Body.Close()

	handleHTTP(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	handleError(err)

	childNodeTable := doc.Find("h3:contains('Children')").Next()

	childNodeTable.Find("tr").Each(func(index int, rowHtml *goquery.Selection) {
		nodeURL, _ := rowHtml.Find("td a").Attr("href")
		nodeName := rowHtml.Find("td a").Text()
		subChildren := rowHtml.Find("td:nth-child(3)").Text()
		subNodesTotal := rowHtml.Find("td:nth-child(4)").Text()

		work_DB.InsertIntoChild_nodes(nodeName, subChildren, subNodesTotal, db)

		TotalNodes += 1
		if TotalNodes%100 == 0 {
			fmt.Printf("Processed nodes: %d\n", TotalNodes)
		}

		if nodeURL != "" {
			nodeURL = MainUrl + nodeURL
			child_nodes = append(child_nodes, nodeURL)
		}
	})
	if len(child_nodes) != 0 {
		WorkChildNodes(child_nodes, db)
	}
}

func WorkChildNodes(child_nodes []string, db *sql.DB) {
	for _, childURL := range child_nodes {
		if !work_DB.IsNodeProcessed(childURL, db) {
			ProcessChildNodes(childURL, db)
			work_DB.MarkNodeAsProcessed(childURL, db)
		}
	}
}

func handleHTTP(StatusCode int) {
	if StatusCode != 200 {
		log.Fatalf("Request failed with status code %d", StatusCode)
	}
}

func retryHttpGet(url string, maxRetries int) (*http.Response, error) {
	for i := 0; i < maxRetries; i++ {
		res, err := http.Get(url)
		if err == nil && res.StatusCode == 200 {
			return res, nil
		}
		fmt.Printf("Attempt %d: Error fetching %s: %v\n", i+1, url, err)
		time.Sleep(time.Second * 1)
	}

	return nil, fmt.Errorf("failed to fetch %s after %d attempts", url, maxRetries)
}
