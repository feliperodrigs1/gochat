package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	ID 	       uint   `gorm:"primaryKey"`
	PublicId string   `gorm:"uniqueIndex;not null"`
	UserID     uint   `gorm:"not null"`
	Title	   string `gorm:"not null"`
	Chunks	   []Chunk
}
