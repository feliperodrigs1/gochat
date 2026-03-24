package models

type Message struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"not null"`
	Role		   string `gorm:"type:varchar(20);not null"`
	Content        string `gorm:"type:text;not null"`
}