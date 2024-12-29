package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func New(storagePath string) (*DB, error) {
	const op = "db.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url (
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		create index if not exists idx_alias on url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &DB{db: db}, nil

}

func (db *DB) SaveURl(urlToSave string, alias string) (int64, error) {
	const op = "db.sqlite.SaveURl"

	stmt, err := db.db.Prepare(
		`INSERT INTO url (alias, url) VALUES (?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (db *DB) GetUrl(alias string) (string, error) {
	const op = "db.sqlite.GetUrl"

	stmt, err := db.db.Prepare(
		`select url from url wheare alias = ?`,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var url string

	err = stmt.QueryRow(alias).Scan(&url)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil

}

func (db *DB) DeleteUrl(alias string) error {
	const op = "db.sqlite.GetUrl"

	stmt, err := db.db.Prepare(
		`delete from url where alias = ?`,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}
