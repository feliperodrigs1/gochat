package models

type Conversation struct {
	ID		   uint   `gorm:"primaryKey"`
	DocumentID uint   `gorm:"not null"`
	ExternalID string   `gorm:"not null"`
}