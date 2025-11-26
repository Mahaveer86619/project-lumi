package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
)

type WahaService struct {
	httpClient *http.Client
}

func NewWahaService() *WahaService {
	return &WahaService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WahaService) StartSession(sessionName string) error {
	url := fmt.Sprintf("%s/api/sessions", config.GConfig.WahaServiceURL)

	payload := map[string]string{"name": sessionName}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode != 409 && resp.StatusCode != 400 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to start session: %d %s", resp.StatusCode, string(body))
		}
	}
	return nil
}

func (s *WahaService) GetQRCode(sessionName string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s/auth/qr?format=image", config.GConfig.WahaServiceURL, sessionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get QR code, status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (s *WahaService) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", config.GConfig.WahaAPIKey)
}
