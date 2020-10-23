package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todos", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todo/43", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Todo Not Found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Todo Not Found'. Got '%s'.", m["error"])
	}
}

func TestAddTodo(t *testing.T) {
	clearTable()

	body := []byte(`{"todo": "Add Test Todo"}`)
	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "applcation/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["todo"] != "Add Test Todo" {
		t.Errorf("Expected todo to be 'Add Test Todo'. Got '%v'.", m["todo"])
	}

	if m["completed"] != false {
		t.Errorf("Expected completed to default to false. Got %v.", m["completed"])
	}
}

func TestGetTodo(t *testing.T) {
	clearTable()
	// Execute an add to insert item into database to get it.
	a.DB.Exec("INSERT INTO todos(todo) VALUES($1)", "Test Get Todo Add")

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["todo"] != "Test Get Todo Add" {
		t.Errorf("Expected todo be to 'Test Get Todo Add'. Got '%s'.", m["todo"])
	}
}

func TestUpdateTodo(t *testing.T) {
	clearTable()
	// Execute an addTodo to update its value
	a.DB.Exec("INSERT INTO todos(todo) VALUES($1)", "Test Update Todo Add")

	body := []byte(`{"todo": "Updated todo", "completed": true}`)
	req, _ := http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != 1.0 {
		t.Errorf("Expected id to stay the same at 1. Got %d.", m["id"])
	}

	if m["todo"] != "Updated todo" {
		t.Errorf("Expected todo to be 'Updated todo'. Got '%s'.", m["todo"])
	}

	if m["completed"] != true {
		t.Errorf("Expected completed to be true. Got %v.", m["completed"])
	}
}

func TestDeleteTodo(t *testing.T) {
	clearTable()
	a.DB.Exec("INSERT INTO todos(todo) VALUES($1)", "Delete this")

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func checkResponseCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("Expected response code %d.  Got %d.\n", want, got)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todos
(
	id SERIAL,
	todo TEXT NOT NULL,
	completed BOOL DEFAULT false,
	CONSTRAINT todo_pkey PRIMARY KEY (id)
)`

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM todos")
	a.DB.Exec("ALTER SEQUENCE todos_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}
