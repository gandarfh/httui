package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const uris_table string = `
  CREATE TABLE IF NOT EXISTS uris (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name TEXT NOT NULL UNIQUE
  );`

const endpoints_table string = `
  CREATE TABLE IF NOT EXISTS endpoints (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  path TEXT NOT NULL,
  body TEXT,
  headers TEXT,
  method TEXT NOT NULL,
  uri TEXT NULL,
  FOREIGN KEY(uri) REFERENCES uris(name)
  );`

func Conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db?_foreign_keys=on")

	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("PRAGMA foreign_keys = ON;")
	defer stmt.Close()
	stmt.Exec()

	if err := createTable(db, uris_table); err != nil {
		return nil, err
	}

	if err := createTable(db, endpoints_table); err != nil {
		return nil, err
	}

	return db, nil

}

func createTable(db *sql.DB, table string) error {
	query, err := db.Prepare(table)

	if err != nil {
		return err
	}

	query.Exec()

	return nil
}
