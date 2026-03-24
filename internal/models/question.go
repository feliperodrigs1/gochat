package models

type Question struct {
	ID 	   	   uint   `gorm:"primaryKey"`
	DocumentID uint   `gorm:"not null"`
	Question   string `gorm:"not null"`
	Answer     string `gorm:"not null"`
	Embedding  string
}
