package repositories

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var newLogger = logger.New(
	log.Default(),
	logger.Config{
		// LogLevel: logger.Error,
	},
)

var Database *gorm.DB

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

	db.AutoMigrate(&Request{}, &Response{})
	db.AutoMigrate(&Workspace{})
	db.AutoMigrate(&Env{})
	db.AutoMigrate(&Default{})

	Database = db
	return nil
}
