package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
)

type HealthService struct {
	httpClient *http.Client
}

func NewHealthService() *HealthService {
	return &HealthService{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HealthService) GetHealth() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["service"] = "lumi"
	result["status"] = "healthy"
	result["timestamp"] = time.Now()

	authHealth, err := h.checkAuthService()
	if err != nil {
		result["auth_service"] = map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	} else {
		result["auth_service"] = authHealth
	}

	return result, nil
}

func (h *HealthService) checkAuthService() (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/ping", config.GConfig.AuthServiceURL)

	resp, err := h.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}
