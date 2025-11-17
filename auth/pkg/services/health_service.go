package services

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (h *HealthService) GetHealth() (string, error) {
	return "pong", nil
}
