package dto

import "booking-service/internal/constants"

type BookingCreateRequest struct {
	SessionID     uint                    `json:"session_id" binding:"required"`
	UserID        uint                    `json:"user_id" binding:"required"`
	BookingStatus constants.BookingStatus `json:"booking_status"`
}

type BookingUpdateRequest struct {
	BookingStatus *constants.BookingStatus `json:"booking_status"`
}
