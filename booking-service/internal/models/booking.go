package models

import (
	"booking-service/internal/constants"
	"time"
)

type Booking struct {
	Base

	SessionID     uint                    `json:"session_id" gorm:"not null;index"`
	UserID        uint                    `json:"user_id" gorm:"not null;index"`
	BookingStatus constants.BookingStatus `json:"booking_status" gorm:"default:pending;index"`
	PaymentStatus constants.PaymentStatus `json:"payment_status" gorm:"default:pending;index"`
	ExpiresAt     time.Time               `json:"expires_at" gorm:"not null;index"`
	BookedSeats   []BookedSeat            `json:"booked_seats" gorm:"foreignKey:BookingID"`

	SessionStartTime time.Time `json:"session_start_time" gorm:"not null;index"`
	SessionEndTime   time.Time `json:"session_end_time" gorm:"not null;index"`
}

type BookedSeat struct {
	Base

	BookingID uint `json:"booking_id" gorm:"not null;index"`
	SeatID    uint `json:"seat_id" gorm:"not null;index"`
}
