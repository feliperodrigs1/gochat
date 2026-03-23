package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	Title	   string `gorm:"not null"`
	Chunks	   []Chunk
}
