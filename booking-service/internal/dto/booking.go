package dto

import "booking-service/internal/constants"

type BookingCreateRequest struct {
	CinemaID      uint                    `json:"cinema_id" gorm:"not null"`
	UserID        uint                    `json:"user_id" gorm:"not null"`
	BookingStatus constants.BookingStatus `json:"booking_status" gorm:"default:pending"`
}

type BookingUpdateRequest struct {
	BookingStatus *constants.BookingStatus `json:"booking_status" gorm:"not null;index"`
}
