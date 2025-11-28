package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/views"
)

type HealthService struct {
	wahaService *WahaService

	httpClient *http.Client
}

func NewHealthService() *HealthService {
	return &HealthService{
		wahaService: NewWahaService(),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HealthService) GetHealth() (*views.HealthResponse, error) {
	var servicesList []views.Health

	// 1. Check Waha Service
	servicesList = append(servicesList, h.checkWahaService())

	// 2. Check Database
	servicesList = append(servicesList, h.checkDBService())

	// 3. Check Lumi Service
	servicesList = append(servicesList, views.Health{
		Name:    "lumi-service",
		IsUp:    true,
		Message: "Service is running",
	})

	return &views.HealthResponse{
		Services: servicesList,
	}, nil
}

func (h *HealthService) checkWahaService() views.Health {
	err := h.wahaService.PingWaha()
	if err != nil {
		return views.Health{
			Name:    "waha-service",
			IsUp:    false,
			Message: fmt.Sprintf("Ping failed: %v", err),
		}
	}

	return views.Health{
		Name:    "waha-service",
		IsUp:    true,
		Message: "Healthy",
	}
}

func (h *HealthService) checkDBService() views.Health {
	if db.DB == nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: "Database connection not initialized",
		}
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: fmt.Sprintf("Failed to retrieve DB instance: %v", err),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: fmt.Sprintf("Ping failed: %v", err),
		}
	}

	return views.Health{
		Name:    "database",
		IsUp:    true,
		Message: "Healthy",
	}
}
