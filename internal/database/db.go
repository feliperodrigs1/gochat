package database

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gochat/internal/models"
)

var DB *gorm.DB

func Connect() {
	env := os.Getenv("ENV")
	var dbName string

	if env == "test" {
		dbName = ":memory:"
	} else {
		dbName = "data/gochat.db"

		dir := filepath.Dir(dbName)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal("failed to create database directory:", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Chunk{},
	)

	DB = db
}