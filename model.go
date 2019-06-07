package main

import (
	"database/sql"
)

type book struct {
	ID        string `json:"ID"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Date      string `json:"date"`
	Rating    int    `json:"rating"`
	Status    string `json:"status"`
}

func (b *book) getBook(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM books WHERE id=$1",
		b.ID).Scan(&b.Title, &b.Author, &b.Publisher, &b.Date, &b.Rating, &b.Status)
}

func (b *book) updateBook(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE books SET title=$1, author=$2, publisher=$3, date=$4, rating=$5, status=$6 WHERE id=$7",
			b.Title, b.Author, b.Publisher, b.Date, b.Rating, b.Status, b.ID)

	return err
}

func (b *book) deleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", b.ID)

	return err
}

func (b *book) createBook(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO books(title, author, publisher, date, rating, status) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		b.Title, b.Author, b.Publisher, b.Date, b.Rating, b.Status).Scan(&b.ID)

	if err != nil {
		return err
	}

	return nil
}

func getBooks(db *sql.DB, start, count int) ([]book, error) {
	rows, err := db.Query("SELECT * FROM books")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []book{}

	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Publisher, &b.Date, &b.Rating, &b.Status); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}
