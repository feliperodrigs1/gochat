package models

import "gorm.io/gorm"

type Chunk struct {
	gorm.Model
	DocumentID uint   `gorm:"not null"`
	Content    string `gorm:"type:text;not null"`
	Embedding  string
}