package repository

import (
	"booking-service/internal/config"
	"booking-service/internal/models"

	"gorm.io/gorm"
)

type BookingSeatRepository interface {
	Create(tx *gorm.DB, bookingID uint, seatList []uint) error
	DeleteByBookingID(tx *gorm.DB, bookingID uint) error
}

type gormBookingSeat struct {
	db *gorm.DB
}

func NewBookingSeatRepository(db *gorm.DB) BookingSeatRepository {
	return &gormBookingSeat{
		db: db,
	}
}

func (r *gormBookingSeat) Create(tx *gorm.DB, bookingID uint, seatList []uint) error {
	var bookedSeats = make([]models.BookedSeat, 0, len(seatList))

	for _, seat := range seatList {
		bookedSeats = append(bookedSeats, models.BookedSeat{
			BookingID: bookingID,
			SeatID:    seat,
		})
	}

	if err := tx.Create(&bookedSeats).Error; err != nil {
		config.GetLogger().Error("Failed to create booked seats", "error", err, "booking_id", bookingID, "seats", seatList)
		return err
	}

	return nil
}

func (r *gormBookingSeat) DeleteByBookingID(tx *gorm.DB, bookingID uint) error {
	if err := tx.Where("booking_id = ?", bookingID).Delete(&models.BookedSeat{}).Error; err != nil {
		config.GetLogger().Error("Failed to delete booked seats by booking_id", "error", err, "booking_id", bookingID)
		return err
	}

	return nil
}
