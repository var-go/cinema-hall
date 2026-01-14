package models

import "booking-service/internal/constants"

type Booking struct {
	Base

	CinemaID      uint                    `json:"cinema_id" gorm:"not null;index"`
	UserID        uint                    `json:"user_id" gorm:"not null;index"`
	BookingStatus constants.BookingStatus `json:"booking_status" gorm:"default:pending;index"`
}
