package main

import (
	"database/sql"
	"errors"
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
	return errors.New("Not implemented")
}

func (b *book) updateBook(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (b *book) deleteBook(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (b *book) createBook(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getBooks(db *sql.DB, start, count int) ([]book, error) {
	return nil, errors.New("Not implemented")
}
