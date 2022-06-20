package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const uris_table string = `
  CREATE TABLE IF NOT EXISTS uris (
  uri TEXT PRIMARY KEY NOT NULL
  );`

const endpoints_table string = `
  CREATE TABLE IF NOT EXISTS endpoints (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  path TEXT NOT NULL,
  body TEXT,
  headers TEXT,
  method TEXT NOT NULL,
  uri INTEGER NOT NULL,
  FOREIGN KEY (uri)
       REFERENCES uris (uri)
  );`

func Conn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		return nil, err
	}

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
