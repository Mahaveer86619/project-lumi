package models

import (
	"gorm.io/gorm"
)

type RegisteredChat struct {
	gorm.Model
	ChatID      string `gorm:"uniqueIndex;not null" json:"chat_id"` // e.g. 123@c.us
	Name        string `json:"name"`                                // Friendly name
	Type        string `json:"type"`                                // "chat" or "group"
	IsBotActive bool   `gorm:"default:false" json:"is_bot_active"`  // Is the NLP session active?
}

type ChatMessage struct {
	gorm.Model
	ChatID  string `gorm:"index;not null"`
	Role    string `json:"role"`    // "user" or "model"
	Content string `json:"content"` // Text content
}
