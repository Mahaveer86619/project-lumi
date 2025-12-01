package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/services/bot"
	connService "github.com/Mahaveer86619/lumi/pkg/services/connections"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type WahaHandler struct {
	wahaService connService.WahaClient
	chatService *services.ChatService
	botService  *bot.BotService
}

func NewWahaHandler(group *echo.Group, wahaService connService.WahaClient, chatService *services.ChatService, botService *bot.BotService) *WahaHandler {
	handler := &WahaHandler{
		wahaService: wahaService,
		chatService: chatService,
		botService:  botService,
	}

	group.GET("/connect", handler.ConnectWhatsApp)
	group.GET("/code", handler.RequestCode)
	group.POST("/start", handler.StartDefaultSession)

	group.GET("/me", handler.GetMe)

	group.POST("/send/text", handler.SendText)
	group.POST("/send/image", handler.SendImage)

	return handler
}

func (h *WahaHandler) HandleWebhook(c echo.Context) error {
	var webhook connections.WAHAWebhook
	if err := c.Bind(&webhook); err != nil {
		return c.JSON(http.StatusBadRequest, views.Failure{StatusCode: http.StatusBadRequest, Message: "Invalid payload"})
	}

	switch webhook.Event {
	case "session.status":
		var statusPayload connections.SessionStatusPayload
		if err := json.Unmarshal(webhook.Payload, &statusPayload); err != nil {
			log.Printf("Failed to unmarshal session status: %v", err)
			break
		}

		config.SetWhatsappConnectionStatus(statusPayload.Status)

		if statusPayload.Status == "WORKING" {
			go func() {
				profile, err := h.wahaService.GetMe()
				if err == nil && profile != nil {
					h.ensureSelfRegistered(profile)
				}
			}()
		}

	case "message":
		var msg connections.WAMessage
		if err := json.Unmarshal(webhook.Payload, &msg); err != nil {
			log.Printf("Failed to unmarshal message payload: %v", err)
			break
		}

		me, err := h.wahaService.GetMe()
		if err != nil {
			log.Printf("Error fetching me: %v", err)
		}

		isSelfMsg := msg.From == me.ID
		if isSelfMsg {
			log.Printf("Self message: %s", msg.Body)
		}

		chatID := msg.From
		if msg.FromMe {
			chatID = msg.To
		}

		isSelfChat := msg.From == msg.To

		if h.chatService.IsChatAllowed(chatID) || isSelfChat {
			log.Printf("[PROCESSING MESSAGE] Chat: %s, IsSelf: %v", chatID, isSelfChat)
			go h.botService.ProcessMessage(msg)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *WahaHandler) ConnectWhatsApp(c echo.Context) error {
	_, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, views.Failure{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User ID not found",
		})
	}

	err := h.wahaService.StartSession()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to start WhatsApp session: " + err.Error(),
		})
	}

	profile, err := h.wahaService.GetMe()
	if err == nil && profile != nil {
		h.ensureSelfRegistered(profile)

		return c.JSON(http.StatusOK, views.Success{
			StatusCode: http.StatusOK,
			Message:    "Already logged in",
			Data: map[string]interface{}{
				"status":  "connected",
				"profile": profile,
			},
		})
	}

	qrBytes, err := h.wahaService.GetQRCode()
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.Failure{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to retrieve QR code from Waha: " + err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "image/png", qrBytes)
}

func (h *WahaHandler) RequestCode(c echo.Context) error {
	phoneNumber := c.QueryParam("phoneNumber")
	method := c.QueryParam("method")

	if phoneNumber == "" {
		return c.JSON(http.StatusBadRequest, views.Failure{
			StatusCode: http.StatusBadRequest,
			Message:    "phoneNumber query parameter is required",
		})
	}

	if err := h.wahaService.StartSession(); err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to start session: " + err.Error(),
		})
	}

	resp, err := h.wahaService.RequestCode(phoneNumber, method)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to request code: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Code requested successfully",
		Data:       resp,
	})
}

func (h *WahaHandler) StartDefaultSession(c echo.Context) error {
	_, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, views.Failure{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User ID not found",
		})
	}

	err := h.wahaService.StartSession()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to start WhatsApp session: " + err.Error(),
		})
	}

	profile, err := h.wahaService.GetMe()
	if err == nil && profile != nil {
		return c.JSON(http.StatusOK, views.Success{
			StatusCode: http.StatusOK,
			Message:    "Already logged in",
			Data: map[string]interface{}{
				"status":  "connected",
				"profile": profile,
			},
		})
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Profile fetched successfully",
		Data:       profile,
	})
}

func (h *WahaHandler) GetMe(c echo.Context) error {
	_, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, views.Failure{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
		})
	}

	profile, err := h.wahaService.GetMe()
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.Failure{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to fetch profile. Ensure session is connected: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Profile fetched successfully",
		Data:       profile,
	})
}

func (h *WahaHandler) SendText(c echo.Context) error {
	var req views.SendTextChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, views.Failure{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	if !h.chatService.IsChatAllowed(req.ChatID) {
		return c.JSON(http.StatusForbidden, views.Failure{
			StatusCode: http.StatusForbidden,
			Message:    "Chat ID is not registered. Please register the chat/group first.",
		})
	}

	resp, err := h.wahaService.SendText(req.ChatID, req.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, views.Success{StatusCode: http.StatusOK, Message: "Message sent", Data: resp})
}

func (h *WahaHandler) SendImage(c echo.Context) error {
	var req connections.MessageImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, views.Failure{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	if !h.chatService.IsChatAllowed(req.ChatID) {
		return c.JSON(http.StatusForbidden, views.Failure{
			StatusCode: http.StatusForbidden,
			Message:    "Chat ID is not registered.",
		})
	}

	imagePayload := connections.ImagePayload{Caption: req.Caption, File: req.File}
	resp, err := h.wahaService.SendImage(req.ChatID, imagePayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{StatusCode: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, views.Success{StatusCode: http.StatusOK, Message: "Image sent", Data: resp})
}

func (h *WahaHandler) ensureSelfRegistered(profile *connections.MeInfo) {
	if profile == nil || profile.ID == "" {
		return
	}

	if !h.chatService.IsChatAllowed(profile.ID) {
		name := profile.PushName
		if name == "" {
			name = "Me"
		}
		_, err := h.chatService.RegisterChat(profile.ID, name+" (Self)", "self")
		if err != nil {
			log.Printf("Failed to auto-register self chat: %v", err)
		} else {
			log.Printf("Auto-registered self chat: %s (%s)", profile.ID, name)
		}
	}
}
