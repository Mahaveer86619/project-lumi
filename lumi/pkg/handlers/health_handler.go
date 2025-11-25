package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(group *echo.Group, healthService *services.HealthService) *HealthHandler {
	handler := &HealthHandler{
		healthService: healthService,
	}

	group.GET("/health", handler.GetHealth)

	return handler
}

func (h *HealthHandler) GetHealth(c echo.Context) error {
	resp, err := h.healthService.GetHealth()
	if err != nil {
		res := &views.Failure{}
		res.SetStatusCode(http.StatusBadRequest)
		res.SetMessage("Health check failed")
		return res.JSON(c)
	}

	success := &views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("Health check successful")
	success.SetData(resp)

	return success.JSON(c)
}
