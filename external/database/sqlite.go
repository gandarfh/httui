package database

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{},
)

func SqliteConnection() (*gorm.DB, error) {
	home, _ := os.UserHomeDir()
	db, err := gorm.Open(sqlite.Open(filepath.Join(home, "httui.db")), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	sqldb, _ := db.DB()

	if err := sqldb.Ping(); err != nil {
		defer sqldb.Close()
		return nil, err
	}

	return db, nil
}
