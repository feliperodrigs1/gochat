package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gochat/internal/models"
)

var DB *gorm.DB

func Connect() {
	env := os.Getenv("ENV")
	dbName := "docchat.db"

	if env == "test" {
		dbName = ":memory:"
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