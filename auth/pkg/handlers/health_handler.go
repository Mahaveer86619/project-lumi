package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/ms/auth/pkg/services"
	"github.com/Mahaveer86619/ms/auth/pkg/views"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(group *echo.Group, healthService *services.HealthService) *HealthHandler {
	handler := &HealthHandler{
		healthService: healthService,
	}

	group.GET("/ping", handler.GetHealth)

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

	jResp := &views.Success{}
	jResp.SetStatusCode(http.StatusOK)
	jResp.SetMessage("Health check successful")
	jResp.SetData(resp)
	
	return jResp.JSON(c)
}
