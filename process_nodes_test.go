package main

import (
	"database/sql"
	"eshelon/work_DB"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRetryHttpGet_Successful проверяет успешный HTTP GET запрос и обработку успешного ответа.
func TestRetryHttpGet_Successful(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Successful response")
	}))
	defer server.Close()

	maxRetries := 3
	res, err := retryHttpGet(server.URL, maxRetries)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}

// TestRetryHttpGet_MaxRetriesExceeded проверяет обработку ситуации, когда превышено максимальное количество попыток.
func TestRetryHttpGet_MaxRetriesExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer server.Close()

	maxRetries := 2
	res, err := retryHttpGet(server.URL, maxRetries)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if res != nil {
		t.Errorf("Expected response to be nil, got %+v", res)
	}
}

// TestRetryHttpGet_NetworkError проверяет обработку сетевой ошибки при HTTP GET запросе.
func TestRetryHttpGet_NetworkError(t *testing.T) {
	nonExistentURL := "http://localhost:9999"
	maxRetries := 3
	res, err := retryHttpGet(nonExistentURL, maxRetries)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if res != nil {
		t.Errorf("Expected response to be nil, got %+v", res)
	}
}

// TestRetryHttpGet_SuccessAfterRetries проверяет успешный HTTP GET запрос после нескольких попыток.
func TestRetryHttpGet_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Internal Server Error")
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Successful response")
		}
	}))
	defer server.Close()

	maxRetries := 2
	res, err := retryHttpGet(server.URL, maxRetries)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}

// TestProcessChildNodes проверяет обработку данных дочерних узлов.
func TestProcessChildNodes(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при создании фиктивной базы данных: %v", err)
	}
	defer db.Close()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<h3>Children (505)</h3>
		<table>
		
			<tr><th>OID</th><th>Name</th><th>Sub children</th><th>Sub Nodes Total</th></tr>
			
				<tr>
					<td><a href="" title="Node 1">Node 1</a></td>
					<td>1111</td>
					<td>0</td>
					<td>0</td>
				</tr>
			</table>
		</html>
        `))
	}))
	defer mockServer.Close()

	work_DB.CreateDBChildNodes(db)
	work_DB.CreateProgressTable(db)
	ProcessChildNodes(mockServer.URL, db)

	var subChildren1, subNodesTotal1 string
	err = db.QueryRow("SELECT sub_children, sub_nodes_total FROM child_nodes WHERE node_name = 'Node 1'").Scan(&subChildren1, &subNodesTotal1)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса к базе данных: %v", err)
	}

	if subChildren1 != "0" {
		t.Errorf("Ожидалось, что subChildren1 равно '0', но получено: %s", subChildren1)
	}

	if subNodesTotal1 != "0" {
		t.Errorf("Ожидалось, что subNodesTotal1 равно '0', но получено: %s", subNodesTotal1)
	}
}

// TestProcessRootNodes проверяет обработку данных корневых узлов.
func TestProcessRootNodes(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Ошибка при создании фиктивной базы данных: %v", err)
	}
	defer db.Close()

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<h2>Root Tree Nodes</h2>
		<table>
			<tr><th>Node</th><th>Name</th><th>Sub children</th><th>Sub Nodes Total</th><th>Description</th><th>Information</th></tr>
			<tr><td><a href="">0</a></td><td>Name1</td><td>0</td><td>0</td><td>Description1</td><td>Information1</td></tr>
		</table>
		</html>		
        `))
	}))
	defer mockServer.Close()

	work_DB.CreateDBrootNode(db)
	ProcessRootNodes(mockServer.URL, db)

	var name, sub_children, sub_nodes_total, description, information string
	err = db.QueryRow("SELECT node_name, subChildren, subNodesTotal, description, information FROM root_nodes WHERE node_name = ?", "Name1").Scan(&name, &sub_children, &sub_nodes_total, &description, &information)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса к базе данных: %v", err)
	}
	if sub_children != "0" {
		t.Errorf("Ожидалось, что sub_children будет '0', но получено: %s", sub_children)
	}
	if sub_nodes_total != "0" {
		t.Errorf("Ожидалось, что sub_nodes_total будет '0', но получено: %s", sub_nodes_total)
	}
	if description != "Description1" {
		t.Errorf("Ожидалось, что description будет 'Description1', но получено: %s", description)
	}
	if information != "Information1" {
		t.Errorf("Ожидалось, что information будет 'Information1', но получено: %s", information)
	}
}
