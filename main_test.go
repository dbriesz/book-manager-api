package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/dbriesz/book-manager-api"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM books")
	a.DB.Exec("ALTER SEQUENCE books_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS books
(
	id SERIAL,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	publisher TEXT NOT NULL,
	date TEXT NOT NULL,
	rating NUMERIC(10) NOT NULL DEFAULT 0,
	status TEXT NOT NULL,
	CONSTRAINT books_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/books", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentBook(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/book/42", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", m["error"])
	}
}

func TestCreateBook(t *testing.T) {
	clearTable()

	payload := []byte(
		`{"title":"test title",
		"author":"test author",
		"publisher":"test publisher",
		"date":"1/1/2019",
		"rating":3,
		"status":"CheckedOut"}`)

	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "test title" {
		t.Errorf("Expected book title to be 'test title'. Got '%v'", m["title"])
	}

	if m["author"] != "test author" {
		t.Errorf("Expected book author to be 'test author'. Got '%v'", m["author"])
	}

	if m["publisher"] != "test publisher" {
		t.Errorf("Expected book publisher to be 'test publisher'. Got '%v'", m["publisher"])
	}

	if m["date"] != "1/1/2019" {
		t.Errorf("Expected publish date to be '1/1/2019'. Got '%v'", m["date"])
	}

	if m["rating"] != 3 {
		t.Errorf("Expected book rating to be '3'. Got '%v'", m["rating"])
	}

	if m["status"] != "CheckedOut" {
		t.Errorf("Expected book status to be 'CheckedOut'. Got '%v'", m["status"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected book ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetBook(t *testing.T) {
	clearTable()
	addBooks(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addBooks(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO books(title, author, publisher, date, rating, status) VALUES($1, $2, $3, $4, $5, $6)", "Book "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestUpdateBook(t *testing.T) {
	clearTable()
	addBooks(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	payload := []byte(`{"title":"test book - updated title","price":11.22}`)

	req, _ = http.NewRequest("PUT", "/book/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalBook["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalBook["id"], m["id"])
	}

	if m["title"] == originalBook["title"] {
		t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalBook["title"], m["title"], m["title"])
	}

	if m["author"] == originalBook["author"] {
		t.Errorf("Expected the author to change from '%v' to '%v'. Got '%v'", originalBook["author"], m["author"], m["author"])
	}

	if m["publisher"] == originalBook["publisher"] {
		t.Errorf("Expected the publisher to change from '%v' to '%v'. Got '%v'", originalBook["publisher"], m["publisher"], m["publisher"])
	}

	if m["date"] == originalBook["date"] {
		t.Errorf("Expected the date to change from '%v' to '%v'. Got '%v'", originalBook["date"], m["date"], m["date"])
	}

	if m["rating"] == originalBook["rating"] {
		t.Errorf("Expected the rating to change from '%v' to '%v'. Got '%v'", originalBook["rating"], m["rating"], m["rating"])
	}

	if m["status"] == originalBook["status"] {
		t.Errorf("Expected the status to change from '%v' to '%v'. Got '%v'", originalBook["status"], m["status"], m["status"])
	}
}

func TestDeleteBook(t *testing.T) {
	clearTable()
	addBooks(1)

	req, _ := http.NewRequest("GET", "/book/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/book/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/book/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
