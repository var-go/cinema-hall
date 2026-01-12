package dto

import "booking-service/internal/constants"

type BookingCreateRequest struct {
	BookingStatus constants.BookingStatus `json:"booking_status" gorm:"not null;index"`
}

type BookingUpdateRequest struct {
	BookingStatus constants.BookingStatus `json:"booking_status" gorm:"not null;index"`
}
