package dto

import (
	"booking-service/internal/constants"
	"time"
)

type BookingCreateRequest struct {
	SessionID uint   `json:"session_id" binding:"required"`
	UserID    uint   `json:"user_id" binding:"required"`
	SeatsID   []uint `json:"seats_id" binding:"required,min=1"`
}

type BookingUpdateRequest struct {
	BookingStatus *constants.BookingStatus `json:"booking_status"`
}

type SessionResponse struct {
	MovieID   uint      `json:"movie_id"`
	HallID    uint      `json:"hall_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}
