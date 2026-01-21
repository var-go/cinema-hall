package clients

import (
	"booking-service/internal/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const cinemaServiceUrl = "http://localhost:8082"

func GetSession(sessionID uint) (*dto.SessionResponse, error) {
	url := fmt.Sprintf("%s/sessions/%d", cinemaServiceUrl, sessionID)

	resp, err := http.Get(url)
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
