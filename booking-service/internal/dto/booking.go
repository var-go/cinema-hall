package dto

import "booking-service/internal/constants"

type BookingCreateRequest struct {
	SessionID uint   `json:"session_id" binding:"required"`
	UserID    uint   `json:"user_id"`
	SeatsID   []uint `json:"seats_id" binding:"required,min=1"`
}

type BookingUpdateRequest struct {
	BookingStatus *constants.BookingStatus `json:"booking_status"`
}
