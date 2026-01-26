package clients

import (
	"booking-service/internal/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func getCinemaServiceURL() string {
	url := os.Getenv("CINEMA_SERVICE_URL")
	if url == "" {
		return "http://localhost:8081"
	}
	return url
}

func GetSession(sessionID uint) (*dto.SessionResponse, error) {
	cinemaServiceUrl := getCinemaServiceURL()
	url := fmt.Sprintf("%s/sessions/%d", cinemaServiceUrl, sessionID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cinema service returned status %d for session %d", resp.StatusCode, sessionID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var session dto.SessionResponse

	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}

	return &session, nil
}
