package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteConnection() (*gorm.DB, error) {
	fmt.Println("Connecting to database...")
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqldb, _ := db.DB()

	if err := sqldb.Ping(); err != nil {
		defer sqldb.Close()
		return nil, err
	}

	fmt.Println("[Success] | Connected")
	return db, nil
}
