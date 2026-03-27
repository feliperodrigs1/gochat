package database

import (
	"log"
	"os"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gochat/internal/models"
)

var DB *gorm.DB

func Connect() {
	env := os.Getenv("ENV")
	
	var db *gorm.DB
	var err error

	if env == "test" {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			getEnv("DB_HOST", "db"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "gochat"),
			getEnv("DB_PORT", "5432"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Chunk{},
		&models.Question{},
		&models.Conversation{},
		&models.Message{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}