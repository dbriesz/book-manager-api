package main_test

import (
	"os"
	"testing"

	"."
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
