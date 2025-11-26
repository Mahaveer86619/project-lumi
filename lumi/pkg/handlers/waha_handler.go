package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/labstack/echo/v4"
)

type WahaHandler struct {
	wahaService *services.WahaService
}

func NewWahaHandler(group *echo.Group, wahaService *services.WahaService) *WahaHandler {
	handler := &WahaHandler{
		wahaService: wahaService,
	}

	group.GET("/whatsapp/connect", handler.ConnectWhatsApp)

	return handler
}

func (h *WahaHandler) ConnectWhatsApp(c echo.Context) error {
	sessionName := "default"

	err := h.wahaService.StartSession(sessionName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to start WhatsApp session: " + err.Error(),
		})
	}

	qrBytes, err := h.wahaService.GetQRCode(sessionName)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{
			"error": "Failed to retrieve QR code from Waha: " + err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "image/png", qrBytes)
}
