package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	PublicId string  `json:"public_id" gorm:"uniqueIndex;not null"`
	UserID   uint    `json:"-" gorm:"not null"`
	Title    string  `json:"title" gorm:"not null"`
	Chunks   []Chunk
}
