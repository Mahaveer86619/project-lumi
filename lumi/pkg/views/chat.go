package views

import (
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/utils"
)

type RemoteChatListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	LastMessage string `json:"last_message"`
	Timestamp   int64  `json:"timestamp"`
	Type        string `json:"type"`
}

type RegisterChatRequest struct {
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type SendTextChatRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type RegisteredChat struct {
	ID     utils.MaskedId `json:"id"`
	ChatID string         `gorm:"uniqueIndex;not null" json:"chat_id"` // e.g. 123@c.us
	Name   string         `json:"name"`                                // Friendly name
	Type   string         `json:"type"`                                // "chat" or "group"
}

func NewRegisteredChatResponse(chat []models.RegisteredChat) *[]RegisteredChat {
	var resp []RegisteredChat
	for _, c := range chat {
		resp = append(resp, RegisteredChat{
			ID:     utils.Mask(c.ID),
			ChatID: c.ChatID,
			Name:   c.Name,
			Type:   c.Type,
		})
	}
	return &resp
}
