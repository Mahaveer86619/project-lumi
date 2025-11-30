package handlers

import (
	"net/http"
	"strings"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(group *echo.Group, chatService *services.ChatService) *ChatHandler {
	handler := &ChatHandler{chatService: chatService}

	// Remote (from WAHA)
	group.GET("/remote/chats", handler.GetRemoteChats)
	group.GET("/remote/groups", handler.GetRemoteGroups)

	// Local (Registered/Allowed)
	group.GET("/registered", handler.GetRegisteredChats)
	group.POST("/register", handler.RegisterChat)
	group.DELETE("/register/:chatId", handler.UnregisterChat)

	return handler
}

func (h *ChatHandler) GetRemoteChats(c echo.Context) error {
	rawChats, err := h.chatService.GetRemoteChats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: 500, Message: err.Error()})
	}

	var response []views.RemoteChatListResponse

	for _, chat := range rawChats {
		chatType := "chat"
		if strings.HasSuffix(chat.ID, "@g.us") {
			chatType = "group"
		} else if strings.HasSuffix(chat.ID, "@newsletter") {
			chatType = "channel"
		}

		lastMsg := ""
		var timestamp int64 = 0

		if chat.LastMessage != nil {
			timestamp = chat.LastMessage.Timestamp

			if chat.LastMessage.Body != "" {
				lastMsg = chat.LastMessage.Body
			} else {
				if val, ok := chat.LastMessage.Data["caption"]; ok && val != nil {
					if caption, ok := val.(string); ok {
						lastMsg = caption
					}
				}
			}

			if lastMsg == "" {
				msgType := ""
				if val, ok := chat.LastMessage.Data["type"]; ok && val != nil {
					if t, ok := val.(string); ok {
						msgType = t
					}
				}

				if msgType != "" {
					lastMsg = "[" + msgType + "]"
				} else {
					lastMsg = "[unknown]"
				}
			}
		}

		response = append(response, views.RemoteChatListResponse{
			ID:          chat.ID,
			Name:        chat.Name,
			Picture:     chat.Picture,
			LastMessage: lastMsg,
			Timestamp:   timestamp,
			Type:        chatType,
		})
	}

	return c.JSON(http.StatusOK, views.Success{StatusCode: 200, Data: response})
}

func (h *ChatHandler) GetRemoteGroups(c echo.Context) error {
	groups, err := h.chatService.GetRemoteGroups()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: 500, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, views.Success{StatusCode: 200, Data: groups})
}

func (h *ChatHandler) GetRegisteredChats(c echo.Context) error {
	chats, err := h.chatService.GetRegisteredChats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: 500, Message: err.Error()})
	}

	resp := views.NewRegisteredChatResponse(chats)
	return c.JSON(http.StatusOK, views.Success{StatusCode: 200, Data: resp, Message: "All Registered chats fetched"})
}

func (h *ChatHandler) RegisterChat(c echo.Context) error {
	var req views.RegisterChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, views.Failure{StatusCode: 400, Message: "Invalid payload"})
	}

	chat, err := h.chatService.RegisterChat(req.ChatID, req.Name, req.Type)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: 500, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, views.Success{StatusCode: 200, Message: "Chat registered", Data: chat})
}

func (h *ChatHandler) UnregisterChat(c echo.Context) error {
	id := c.Param("chatId")
	if err := h.chatService.UnregisterChat(id); err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: 500, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, views.Success{StatusCode: 200, Message: "Chat unregistered"})
}
