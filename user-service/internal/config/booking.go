package config

import "os"

func BookingServiceURL() string {
	url := os.Getenv("BOOKING_SERVICE_URL")
	if url == "" {
		return "http://localhost:8082"
	}
	return url
}
