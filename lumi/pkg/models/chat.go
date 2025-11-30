package models

import (
	"gorm.io/gorm"
)

type RegisteredChat struct {
	gorm.Model
	ChatID string `gorm:"uniqueIndex;not null" json:"chat_id"` // e.g. 123@c.us
	Name   string `json:"name"`                                // Friendly name
	Type   string `json:"type"`                                // "chat" or "group"
}
