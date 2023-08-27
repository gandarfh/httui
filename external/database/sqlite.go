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

var Client *gorm.DB

func SqliteConnection() error {
	home, _ := os.UserHomeDir()
	db, err := gorm.Open(sqlite.Open(filepath.Join(home, "httui.db")), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return err
	}

	sqldb, _ := db.DB()

	if err := sqldb.Ping(); err != nil {
		defer sqldb.Close()
		return err
	}

	Client = db
	return nil
}
