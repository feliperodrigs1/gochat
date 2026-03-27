package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	PublicID string  `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID   uint    `json:"-" gorm:"not null"`
	Title    string  `json:"title" gorm:"not null"`
	Chunks   []Chunk
}
